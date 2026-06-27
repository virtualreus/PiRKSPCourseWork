#!/usr/bin/env python3
"""API fuzz: Schemathesis (openapi.yaml) + business scenarios. Run via `make fuzz`."""

from __future__ import annotations

import json
import os
import shutil
import subprocess
import sys
import time
import uuid
from dataclasses import dataclass, field
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Callable

from schemathesis import Case
from schemathesis import checks as st_checks
from schemathesis.internal.result import Ok
from hypothesis import settings as hypo_settings
import schemathesis

FUZZ_DIR = Path(__file__).resolve().parent
ROOT = FUZZ_DIR.parent
API_BASE = os.environ.get("FUZZ_API_BASE", "http://127.0.0.1:8080/api/v1").rstrip("/")
REPORT_DIR = Path(os.environ.get("FUZZ_REPORT_DIR", FUZZ_DIR / "reports"))
OPENAPI = ROOT / "openapi.yaml"
MAX_EXAMPLES = int(os.environ.get("SCHEMATHESIS_MAX_EXAMPLES", "12"))
API_WAIT_SEC = int(os.environ.get("FUZZ_API_WAIT_SEC", "180"))

ST_CHECKS = (
    st_checks.not_a_server_error,
    st_checks.status_code_conformance,
    st_checks.content_type_conformance,
)

SCHEMA_PHASES: list[tuple[str, list[str], str | None]] = [
    ("public", ["Hackathons"], None),
    ("participant", ["Users", "Participation", "Teams", "Submissions"], "participant"),
    ("organizer", ["Organizer", "Tracks", "Cases"], "organizer"),
]

DEMO_ACCOUNTS = {
    "participant": ("user@user.ru", "user"),
    "organizer": ("admin@admin.ru", "admin"),
}

BUSINESS_META: dict[str, dict[str, str]] = {
    "health_ok": {
        "title": "Доступность API",
        "purpose": "Проверка healthcheck и подключения к PostgreSQL",
        "expected": "GET /health → 200, status=ok, database=ok",
    },
    "register_weak_password": {
        "title": "Валидация пароля при регистрации",
        "purpose": "Отклонение слишком короткого пароля",
        "expected": "POST /auth/register с password=short → 400",
    },
    "login_invalid_credentials": {
        "title": "Неверные учётные данные",
        "purpose": "Отклонение входа с неверным паролем",
        "expected": "POST /auth/login → 401",
    },
    "register_duplicate_email": {
        "title": "Повторная регистрация email",
        "purpose": "Конфликт при повторном email",
        "expected": "Второй POST /auth/register с тем же email → 409",
    },
    "users_me_requires_auth": {
        "title": "Защита профиля",
        "purpose": "Доступ к профилю только с JWT",
        "expected": "GET /users/me без токена → 401",
    },
    "organizer_forbidden_for_participant": {
        "title": "Разграничение ролей (участник)",
        "purpose": "Участник не может открыть кабинет организатора",
        "expected": "GET /organizer/hackathons с токеном participant → 403",
    },
    "organizer_list_ok": {
        "title": "Доступ организатора",
        "purpose": "Организатор видит свои хакатоны",
        "expected": "GET /organizer/hackathons с токеном organizer → 200",
    },
    "hackathon_not_found": {
        "title": "Несуществующий хакатон",
        "purpose": "Корректная обработка неизвестного UUID",
        "expected": "GET /hackathons/{id} → 404",
    },
    "organizer_cannot_register": {
        "title": "Организатор не регистрируется как участник",
        "purpose": "Роль organizer не может записаться на хакатон",
        "expected": "POST /hackathons/{id}/register → 403",
    },
    "register_and_double_register": {
        "title": "Повторная регистрация на хакатон",
        "purpose": "Запрет второй регистрации на тот же хакатон",
        "expected": "Первый POST → 201, второй POST → 400",
    },
    "no_case_blocks_submission": {
        "title": "Сабмит без выбранного кейса",
        "purpose": "Блокировка сдачи работы при submit_block_reason=no_case",
        "expected": "PUT /teams/{id}/submission без case_id → 400",
    },
    "unregister_blocked_in_team": {
        "title": "Отмена регистрации в команде",
        "purpose": "Нельзя отписаться от хакатона, пока состоишь в команде",
        "expected": "DELETE /hackathons/{id}/register → 400",
    },
    "already_in_team": {
        "title": "Вторая команда на хакатоне",
        "purpose": "Участник не может создать две команды на одном хакатоне",
        "expected": "Второй POST /hackathons/{id}/teams → 400",
    },
}

SCHEMA_PHASE_META: dict[str, str] = {
    "public": "Публичные GET-эндпоинты без авторизации",
    "participant": "GET-эндпоинты с JWT участника (user@user.ru)",
    "organizer": "GET-эндпоинты с JWT организатора (admin@admin.ru)",
}

FUZZ_SERVICES = os.environ.get("FUZZ_COMPOSE_SERVICES", "db,api").split(",")

ST_CHECK_NAMES = "not_a_server_error, status_code_conformance, content_type_conformance"


# ---------------------------------------------------------------------------
# HTTP client
# ---------------------------------------------------------------------------


@dataclass
class CaseResult:
    name: str
    passed: bool
    detail: str = ""
    requests: list[str] = field(default_factory=list)


@dataclass
class SchemaOpResult:
    phase: str
    tag: str
    method: str
    path: str
    passed: bool
    detail: str = ""


@dataclass
class ApiClient:
    base_url: str
    request_log: list[str] = field(default_factory=list)

    def request(
        self,
        method: str,
        path: str,
        *,
        body: dict[str, Any] | None = None,
        token: str | None = None,
        expect: int | tuple[int, ...] | None = None,
    ) -> tuple[int, Any]:
        import urllib.error
        import urllib.request

        url = f"{self.base_url}{path}"
        headers = {"Accept": "application/json"}
        payload = None
        if body is not None:
            headers["Content-Type"] = "application/json"
            payload = json.dumps(body).encode("utf-8")
        if token:
            headers["Authorization"] = f"Bearer {token}"

        req = urllib.request.Request(url, data=payload, headers=headers, method=method)
        try:
            with urllib.request.urlopen(req, timeout=30) as resp:
                status = resp.status
                raw = resp.read().decode("utf-8")
        except urllib.error.HTTPError as exc:
            status = exc.code
            raw = exc.read().decode("utf-8", errors="replace")
        except urllib.error.URLError as exc:
            raise AssertionError(f"{method} {path}: {exc}") from exc

        parsed: Any = None
        if raw:
            try:
                parsed = json.loads(raw)
            except json.JSONDecodeError:
                parsed = raw

        self.request_log.append(f"{method} {path} -> {status}")
        if expect is not None:
            allowed = (expect,) if isinstance(expect, int) else expect
            if status not in allowed:
                raise AssertionError(
                    f"{method} {path}: expected {allowed}, got {status}: {parsed}"
                )
        return status, parsed

    def login(self, email: str, password: str) -> str:
        _, data = self.request(
            "POST",
            "/auth/login",
            body={"email": email, "password": password},
            expect=200,
        )
        token = data.get("access_token")
        if not token:
            raise AssertionError(f"login {email}: missing access_token")
        return token

    def register_participant(self, email: str | None = None) -> tuple[str, str]:
        email = email or f"fuzz-{uuid.uuid4().hex[:12]}@example.com"
        _, data = self.request(
            "POST",
            "/auth/register",
            body={
                "email": email,
                "password": "password123",
                "full_name": "Fuzz User",
                "platform_role": "participant",
            },
            expect=201,
        )
        return email, data["access_token"]


def parse_dt(value: str) -> datetime:
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def hackathon_with_open_registration(client: ApiClient) -> dict[str, Any]:
    _, data = client.request("GET", "/hackathons", expect=200)
    now = datetime.now(timezone.utc)
    for item in data.get("items") or []:
        if item.get("status") != "registration":
            continue
        _, detail = client.request("GET", f"/hackathons/{item['id']}", expect=200)
        timeline = detail.get("timeline") or {}
        opens = parse_dt(timeline["registration_opens_at"])
        closes = parse_dt(timeline["registration_closes_at"])
        if opens <= now < closes:
            return detail
    raise AssertionError("no hackathon with open registration window")


# ---------------------------------------------------------------------------
# Business scenarios
# ---------------------------------------------------------------------------

BusinessTest = Callable[[ApiClient], None]
BUSINESS_TESTS: list[tuple[str, BusinessTest]] = []


def business_case(name: str) -> Callable[[BusinessTest], BusinessTest]:
    def register(fn: BusinessTest) -> BusinessTest:
        BUSINESS_TESTS.append((name, fn))
        return fn

    return register


@business_case("health_ok")
def test_health_ok(client: ApiClient) -> None:
    _, data = client.request("GET", "/health", expect=200)
    if data.get("status") != "ok" or data.get("database") != "ok":
        raise AssertionError(data)


@business_case("register_weak_password")
def test_register_weak_password(client: ApiClient) -> None:
    client.request(
        "POST",
        "/auth/register",
        body={
            "email": f"weak-{uuid.uuid4().hex[:8]}@example.com",
            "password": "short",
            "full_name": "Weak",
            "platform_role": "participant",
        },
        expect=400,
    )


@business_case("login_invalid_credentials")
def test_login_invalid_credentials(client: ApiClient) -> None:
    client.request(
        "POST",
        "/auth/login",
        body={"email": "user@user.ru", "password": "wrong"},
        expect=401,
    )


@business_case("register_duplicate_email")
def test_register_duplicate_email(client: ApiClient) -> None:
    email, _ = client.register_participant()
    client.request(
        "POST",
        "/auth/register",
        body={
            "email": email,
            "password": "password123",
            "full_name": "Duplicate",
            "platform_role": "participant",
        },
        expect=409,
    )


@business_case("users_me_requires_auth")
def test_users_me_requires_auth(client: ApiClient) -> None:
    client.request("GET", "/users/me", expect=401)


@business_case("organizer_forbidden_for_participant")
def test_organizer_forbidden_for_participant(client: ApiClient) -> None:
    email, password = DEMO_ACCOUNTS["participant"]
    token = client.login(email, password)
    client.request("GET", "/organizer/hackathons", token=token, expect=403)


@business_case("organizer_list_ok")
def test_organizer_list_ok(client: ApiClient) -> None:
    email, password = DEMO_ACCOUNTS["organizer"]
    token = client.login(email, password)
    client.request("GET", "/organizer/hackathons", token=token, expect=200)


@business_case("hackathon_not_found")
def test_hackathon_not_found(client: ApiClient) -> None:
    client.request(
        "GET",
        "/hackathons/00000000-0000-4000-8000-000000000000",
        expect=404,
    )


@business_case("organizer_cannot_register")
def test_organizer_cannot_register(client: ApiClient) -> None:
    hackathon = hackathon_with_open_registration(client)
    email, password = DEMO_ACCOUNTS["organizer"]
    token = client.login(email, password)
    client.request(
        "POST",
        f"/hackathons/{hackathon['id']}/register",
        token=token,
        expect=403,
    )


@business_case("register_and_double_register")
def test_register_and_double_register(client: ApiClient) -> None:
    hackathon = hackathon_with_open_registration(client)
    _, token = client.register_participant()
    hid = hackathon["id"]
    client.request("POST", f"/hackathons/{hid}/register", token=token, expect=201)
    client.request("POST", f"/hackathons/{hid}/register", token=token, expect=400)


@business_case("no_case_blocks_submission")
def test_no_case_blocks_submission(client: ApiClient) -> None:
    hackathon = hackathon_with_open_registration(client)
    _, token = client.register_participant()
    hid = hackathon["id"]
    client.request("POST", f"/hackathons/{hid}/register", token=token, expect=201)
    _, team = client.request(
        "POST",
        f"/hackathons/{hid}/teams",
        token=token,
        body={"name": f"Fuzz-{uuid.uuid4().hex[:6]}"},
        expect=201,
    )
    _, participation = client.request(
        "GET",
        f"/hackathons/{hid}/participation",
        token=token,
        expect=200,
    )
    if participation.get("submit_block_reason") != "no_case":
        raise AssertionError(participation.get("submit_block_reason"))
    client.request(
        "PUT",
        f"/teams/{team['id']}/submission",
        token=token,
        body={
            "title": "Draft",
            "summary": "No case selected",
            "repo_url": "https://github.com/example/repo",
        },
        expect=400,
    )


@business_case("unregister_blocked_in_team")
def test_unregister_blocked_in_team(client: ApiClient) -> None:
    hackathon = hackathon_with_open_registration(client)
    _, token = client.register_participant()
    hid = hackathon["id"]
    client.request("POST", f"/hackathons/{hid}/register", token=token, expect=201)
    client.request(
        "POST",
        f"/hackathons/{hid}/teams",
        token=token,
        body={"name": f"Team-{uuid.uuid4().hex[:6]}"},
        expect=201,
    )
    client.request("DELETE", f"/hackathons/{hid}/register", token=token, expect=400)


@business_case("already_in_team")
def test_already_in_team(client: ApiClient) -> None:
    hackathon = hackathon_with_open_registration(client)
    _, token = client.register_participant()
    hid = hackathon["id"]
    client.request("POST", f"/hackathons/{hid}/register", token=token, expect=201)
    client.request(
        "POST",
        f"/hackathons/{hid}/teams",
        token=token,
        body={"name": f"First-{uuid.uuid4().hex[:6]}"},
        expect=201,
    )
    client.request(
        "POST",
        f"/hackathons/{hid}/teams",
        token=token,
        body={"name": f"Second-{uuid.uuid4().hex[:6]}"},
        expect=400,
    )


def run_business_tests() -> tuple[int, list[CaseResult], ApiClient]:
    client = ApiClient(API_BASE)
    results: list[CaseResult] = []

    for name, test_fn in BUSINESS_TESTS:
        start = len(client.request_log)
        try:
            test_fn(client)
            results.append(
                CaseResult(
                    name,
                    True,
                    requests=client.request_log[start:],
                )
            )
        except AssertionError as exc:
            results.append(
                CaseResult(
                    name,
                    False,
                    str(exc),
                    requests=client.request_log[start:],
                )
            )

    passed = sum(1 for item in results if item.passed)
    print(f"  business: {passed}/{len(results)} passed")
    return (0 if passed == len(results) else 1), results, client


# ---------------------------------------------------------------------------
# Schemathesis (OpenAPI contract fuzz)
# ---------------------------------------------------------------------------


def build_schema_test(auth_token: str | None) -> Callable[[Case], None]:
    def test_operation(case: Case) -> None:
        headers = {"Authorization": f"Bearer {auth_token}"} if auth_token else None
        response = case.call(headers=headers)
        case.validate_response(response, checks=ST_CHECKS)

    return test_operation


def run_schemathesis(tokens: dict[str, str]) -> tuple[int, list[SchemaOpResult]]:
    print(f"  tool: schemathesis {schemathesis.__version__}")
    settings = hypo_settings(max_examples=MAX_EXAMPLES, deadline=None)
    REPORT_DIR.mkdir(parents=True, exist_ok=True)
    failed_phases = 0
    all_ops: list[SchemaOpResult] = []
    phase_chunks: list[str] = []

    for label, tags, role in SCHEMA_PHASES:
        token = tokens.get(role) if role else None
        lines = [
            "Schemathesis protocol",
            f"version: {schemathesis.__version__}",
            f"spec: {OPENAPI}",
            f"base_url: {API_BASE}",
            f"phase: {label}",
            f"description: {SCHEMA_PHASE_META.get(label, '')}",
            f"tags: {', '.join(tags)}",
            f"max_examples: {MAX_EXAMPLES}",
            f"checks: {ST_CHECK_NAMES}",
            "",
        ]
        phase_failures = 0
        op_index = 0

        for tag in tags:
            schema = schemathesis.from_path(
                str(OPENAPI),
                base_url=API_BASE,
                method="GET",
                tag=tag,
            )
            test_fn = build_schema_test(token)

            for result in schema.get_all_tests(test_fn, settings=settings):
                if not isinstance(result, Ok):
                    phase_failures += 1
                    lines.append(f"FAIL schema ({tag}): {result}")
                    continue

                operation, hypothesis_test = result.ok()
                op_index += 1
                op_label = f"{operation.method.upper()} {operation.path}"
                try:
                    hypothesis_test()
                    op_result = SchemaOpResult(
                        label, tag, operation.method.upper(), operation.path, True
                    )
                    lines.extend(
                        [
                            f"{op_index}. {op_label} — PASS",
                            f"   Тег OpenAPI: {tag}",
                            f"   Примеры: до {MAX_EXAMPLES} сгенерированных запросов (Hypothesis)",
                            f"   Проверки: {ST_CHECK_NAMES}",
                            "",
                        ]
                    )
                except Exception as exc:
                    phase_failures += 1
                    detail = str(exc).strip() or type(exc).__name__
                    op_result = SchemaOpResult(
                        label, tag, operation.method.upper(), operation.path, False, detail
                    )
                    lines.extend(
                        [
                            f"{op_index}. {op_label} — FAIL",
                            f"   Тег OpenAPI: {tag}",
                            f"   Ошибка: {detail}",
                            "",
                        ]
                    )
                all_ops.append(op_result)

        phase_chunks.append("\n".join(lines))
        status = "FAIL" if phase_failures else "OK"
        print(f"  schemathesis/{label}: {status} ({phase_failures} failures)")
        if phase_failures:
            failed_phases += 1

    (REPORT_DIR / "schemathesis_latest.txt").write_text(
        f"schemathesis {schemathesis.__version__}\n\n" + "\n\n".join(phase_chunks) + "\n",
        encoding="utf-8",
    )
    return (1 if failed_phases else 0), all_ops


# ---------------------------------------------------------------------------
# Docker + reports
# ---------------------------------------------------------------------------


def api_is_healthy() -> bool:
    import urllib.error
    import urllib.request

    try:
        with urllib.request.urlopen(f"{API_BASE}/health", timeout=5) as resp:
            return resp.status == 200
    except (urllib.error.URLError, urllib.error.HTTPError, TimeoutError):
        return False


def wait_for_api() -> None:
    deadline = time.time() + API_WAIT_SEC
    while time.time() < deadline:
        if api_is_healthy():
            return
        time.sleep(2)
    raise RuntimeError(f"API not ready: {API_BASE}/health")


def ensure_docker_stack() -> bool:
    if api_is_healthy():
        print(f"API reachable: {API_BASE}")
        return False

    env_file = ROOT / ".env"
    if not env_file.exists():
        shutil.copy(ROOT / ".env.example", env_file)

    services = FUZZ_SERVICES
    print(f"Starting Docker ({', '.join(services)})...")

    up = subprocess.run(
        ["docker", "compose", "up", "-d", "--wait", *services],
        cwd=ROOT,
        capture_output=True,
        text=True,
    )
    if up.returncode != 0:
        print("  images missing or outdated, rebuilding api...")
        up = subprocess.run(
            ["docker", "compose", "up", "--build", "-d", "--wait", *services],
            cwd=ROOT,
            capture_output=True,
            text=True,
        )

    if up.returncode != 0:
        if up.stderr:
            print(up.stderr, file=sys.stderr)
        raise RuntimeError(
            "Не удалось поднять db + api. "
            "Проверьте Docker и сеть (Docker Hub). "
            "Или заранее: docker compose up -d --wait db api"
        )

    wait_for_api()
    return True


def stop_docker_stack() -> None:
    services = FUZZ_SERVICES
    print(f"Stopping Docker ({', '.join(services)})...")
    subprocess.run(
        ["docker", "compose", "stop", *services],
        cwd=ROOT,
        check=False,
    )
    subprocess.run(
        ["docker", "compose", "rm", "-f", *services],
        cwd=ROOT,
        check=False,
    )


def format_business_section(results: list[CaseResult]) -> list[str]:
    lines = [
        "ЧАСТЬ 1. БИЗНЕС-СЦЕНАРИИ",
        f"Всего тестов: {len(results)}",
        "",
    ]
    for index, item in enumerate(results, start=1):
        meta = BUSINESS_META.get(item.name, {})
        status = "PASS" if item.passed else "FAIL"
        lines.append(f"{index}. {item.name} — {status}")
        if meta.get("title"):
            lines.append(f"   Название: {meta['title']}")
        if meta.get("purpose"):
            lines.append(f"   Цель: {meta['purpose']}")
        if meta.get("expected"):
            lines.append(f"   Ожидание: {meta['expected']}")
        if item.detail:
            lines.append(f"   Ошибка: {item.detail}")
        if item.requests:
            lines.append("   Запросы:")
            for req in item.requests:
                lines.append(f"     • {req}")
        else:
            lines.append("   Запросы: (нет)")
        lines.append("")
    return lines


def format_schema_section(ops: list[SchemaOpResult]) -> list[str]:
    passed = sum(1 for op in ops if op.passed)
    lines = [
        "ЧАСТЬ 2. SCHEMATHESIS (фаззинг по openapi.yaml)",
        f"Инструмент: Schemathesis {schemathesis.__version__}",
        f"Метод: GET-эндпоинты, до {MAX_EXAMPLES} примеров на операцию (Hypothesis)",
        f"Проверки: {ST_CHECK_NAMES}",
        f"Операций: {len(ops)}, успешно: {passed}",
        "",
    ]
    current_phase = ""
    index = 0
    for op in ops:
        if op.phase != current_phase:
            current_phase = op.phase
            index = 0
            lines.append(f"--- Фаза: {op.phase} ---")
            lines.append(SCHEMA_PHASE_META.get(op.phase, ""))
            lines.append("")
        index += 1
        status = "PASS" if op.passed else "FAIL"
        lines.append(f"  {index}. {op.method} {op.path} — {status}")
        lines.append(f"     Тег OpenAPI: {op.tag}")
        lines.append(f"     Примеры: до {MAX_EXAMPLES} сгенерированных запросов")
        if op.detail:
            lines.append(f"     Ошибка: {op.detail}")
        lines.append("")
    return lines


def write_reports(
    business_results: list[CaseResult],
    schema_ops: list[SchemaOpResult],
    business_rc: int,
    schema_rc: int,
) -> Path:
    REPORT_DIR.mkdir(parents=True, exist_ok=True)
    finished_at = datetime.now(timezone.utc)
    overall_ok = business_rc == 0 and schema_rc == 0
    business_passed = sum(1 for item in business_results if item.passed)
    schema_passed = sum(1 for op in schema_ops if op.passed)

    divider = "=" * 78
    section = "-" * 78

    report_lines = [
        divider,
        "ПРОТОКОЛ ТЕСТИРОВАНИЯ API HackathonHub",
        divider,
        f"Дата (UTC): {finished_at.isoformat()}",
        f"Целевой URL: {API_BASE}",
        f"Спецификация: {OPENAPI}",
        f"Скрипт: fuzzing/fuzz.py",
        f"Запуск: make fuzz",
        "",
        section,
        *format_business_section(business_results),
        section,
        *format_schema_section(schema_ops),
        section,
        "ИТОГ",
        section,
        f"Бизнес-сценарии: {business_passed}/{len(business_results)} PASS",
        f"Schemathesis: {schema_passed}/{len(schema_ops)} операций PASS",
        f"Общий результат: {'PASS' if overall_ok else 'FAIL'}",
        divider,
    ]

    report_path = REPORT_DIR / "report_latest.txt"
    report_path.write_text("\n".join(report_lines) + "\n", encoding="utf-8")

    (REPORT_DIR / "summary_latest.txt").write_text(
        "\n".join(
            [
                "API fuzz summary",
                f"finished_at: {finished_at.isoformat()}",
                f"target: {API_BASE}",
                "",
                f"business: {business_passed}/{len(business_results)} PASS",
                f"schemathesis: {schema_passed}/{len(schema_ops)} PASS",
                f"overall: {'PASS' if overall_ok else 'FAIL'}",
                "",
                "Подробный протокол: report_latest.txt",
            ]
        )
        + "\n",
        encoding="utf-8",
    )
    return report_path


def login_demo_accounts(client: ApiClient) -> dict[str, str]:
    tokens: dict[str, str] = {}
    for role, (email, password) in DEMO_ACCOUNTS.items():
        tokens[role] = client.login(email, password)
    return tokens


def main() -> int:
    started_stack = False
    print("HackathonHub API fuzz")
    print(f"target: {API_BASE}")

    try:
        started_stack = ensure_docker_stack()

        print("\n[1/2] Business scenarios")
        business_rc, business_results, client = run_business_tests()

        print("\n[2/2] Schemathesis (openapi.yaml)")
        tokens = login_demo_accounts(client)
        schema_rc, schema_ops = run_schemathesis(tokens)

        report_path = write_reports(
            business_results, schema_ops, business_rc, schema_rc
        )

        overall = business_rc == 0 and schema_rc == 0
        print(f"\nResult: {'PASS' if overall else 'FAIL'}")
        print(f"Reports: {REPORT_DIR}")
        print(f"Protocol: {report_path}")
        return 0 if overall else 1
    except RuntimeError as exc:
        print(f"\nError: {exc}", file=sys.stderr)
        return 1
    finally:
        if started_stack:
            stop_docker_stack()


if __name__ == "__main__":
    sys.exit(main())
