import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import * as hackathonsApi from "../api/hackathons";
import type { HackathonDetail } from "../api/hackathonTypes";
import * as participationApi from "../api/participation";
import type { ParticipationStatus } from "../api/participationTypes";
import { ApiError } from "../api/client";
import { ParticipationAlert } from "../components/ParticipationAlert";
import { Reveal } from "../components/Reveal";
import { useAuth } from "../context/AuthContext";
import { formatDate, statusLabel } from "../utils/hackathon";
import {
  formatDeadlineRemaining,
  getStep3Hint,
  getSubmitBlockInfo,
  participationProgress,
} from "../utils/participation";

function stepChip(done: boolean, active: boolean, locked: boolean): string {
  if (done) {
    return "step-chip step-chip-done";
  }
  if (locked) {
    return "step-chip step-chip-locked";
  }
  if (active) {
    return "step-chip step-chip-active";
  }
  return "step-chip";
}

export function ParticipatePage() {
  const { id } = useParams();
  const { user } = useAuth();
  const [hackathon, setHackathon] = useState<HackathonDetail | null>(null);
  const [status, setStatus] = useState<ParticipationStatus | null>(null);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);

  async function reload() {
    if (!id) {
      return;
    }
    const [h, p] = await Promise.all([
      hackathonsApi.getHackathon(id),
      participationApi.getParticipation(id),
    ]);
    setHackathon(h);
    setStatus(p);
  }

  useEffect(() => {
    if (!id) {
      return;
    }
    (async () => {
      try {
        await reload();
      } catch (err) {
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError("Не удалось загрузить участие");
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  async function handleRegister() {
    if (!id) {
      return;
    }
    setBusy(true);
    setError("");
    setMessage("");
    try {
      await participationApi.registerForHackathon(id);
      setMessage("Вы зарегистрированы на хакатон");
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleUnregister() {
    if (!id || !confirm("Отменить регистрацию?")) {
      return;
    }
    setBusy(true);
    setError("");
    setMessage("");
    try {
      await participationApi.unregisterFromHackathon(id);
      setMessage("Регистрация отменена");
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    } finally {
      setBusy(false);
    }
  }

  if (loading) {
    return (
      <div className="page-loading" aria-busy>
        <div className="spinner" />
      </div>
    );
  }

  if (error && !status) {
    return (
      <div className="card glass form-card">
        <p className="form-error">{error}</p>
        <Link to={id ? `/hackathons/${id}` : "/"}>Назад к хакатону</Link>
      </div>
    );
  }

  if (!hackathon || !status || !id) {
    return null;
  }

  const step1Done = status.is_registered;
  const step2Done = Boolean(status.team);
  const step3Done = status.has_submission;
  const progress = participationProgress(status);
  const isCaptain = user?.id === status.team?.captain_id;
  const submitBlock =
    status.team && !status.can_submit
      ? getSubmitBlockInfo(status.submit_block_reason, id, Boolean(isCaptain))
      : null;
  const deadlineLeft = status.submission_deadline_at
    ? formatDeadlineRemaining(status.submission_deadline_at)
    : null;

  return (
    <div className="participate-page participate-page-wide">
      <Reveal>
        <nav className="breadcrumb">
          <Link to={`/hackathons/${id}`}>{hackathon.title}</Link>
          <span>Участие</span>
        </nav>
        <header className="participate-hero">
          <div>
            <h1 className="participate-title">Участие в хакатоне</h1>
            <p className="participate-lead">
              {hackathon.short_description ?? hackathon.title} ·{" "}
              {statusLabel(status.hackathon_status)}
            </p>
          </div>
          {status.submission_deadline_at && (
            <div className="deadline-card">
              <span className="deadline-label">Дедлайн сдачи</span>
              <strong>{formatDate(status.submission_deadline_at)}</strong>
              {deadlineLeft && (
                <span className="deadline-remaining">{deadlineLeft}</span>
              )}
            </div>
          )}
        </header>
      </Reveal>

      <Reveal delay={40}>
        <div className="participate-progress card glass">
          <div className="progress-header">
            <span>Прогресс участия</span>
            <strong>
              {progress.done} из {progress.total} шагов
            </strong>
          </div>
          <div
            className="progress-bar"
            role="progressbar"
            aria-valuenow={progress.percent}
          >
            <div
              className="progress-fill"
              style={{ width: `${progress.percent}%` }}
            />
          </div>
        </div>
      </Reveal>

      {error && <p className="form-error participate-banner">{error}</p>}
      {message && <p className="form-success participate-banner">{message}</p>}

      {submitBlock && step2Done && !step3Done && (
        <ParticipationAlert
          variant={submitBlock.variant}
          title={submitBlock.title}
          actionLabel={submitBlock.actionLabel}
          actionTo={submitBlock.actionTo}
        >
          {submitBlock.body}
        </ParticipationAlert>
      )}

      <div className="participate-steps">
        <Reveal delay={80}>
          <section
            className={`card glass participate-step ${step1Done ? "step-done" : ""} ${!status.can_register && !step1Done ? "step-disabled" : ""}`}
          >
            <div className="step-indicator">{step1Done ? "✓" : "1"}</div>
            <div className="step-body">
              <div className="step-title-row">
                <h2>Регистрация на хакатон</h2>
                <span className={stepChip(step1Done, !step1Done, false)}>
                  {step1Done
                    ? "Готово"
                    : status.can_register
                      ? "Нужно действие"
                      : "Закрыто"}
                </span>
              </div>
              <p>
                {step1Done
                  ? `Вы в списке участников с ${formatDate(status.registration?.registered_at ?? "")}.`
                  : "Подтвердите участие - это откроет создание команды и сдачу решения."}
              </p>
              {step1Done ? (
                <button
                  type="button"
                  className="btn-ghost"
                  disabled={busy || Boolean(status.team)}
                  onClick={handleUnregister}
                  title={status.team ? "Сначала выйдите из команды" : undefined}
                >
                  Отменить регистрацию
                </button>
              ) : (
                <button
                  type="button"
                  className="btn-primary"
                  disabled={busy || !status.can_register}
                  onClick={handleRegister}
                >
                  Зарегистрироваться
                </button>
              )}
            </div>
          </section>
        </Reveal>

        <Reveal delay={160}>
          <section
            className={`card glass participate-step ${step2Done ? "step-done" : ""} ${!step1Done ? "step-disabled" : ""}`}
          >
            <div className="step-indicator">{step2Done ? "✓" : "2"}</div>
            <div className="step-body">
              <div className="step-title-row">
                <h2>Команда</h2>
                <span
                  className={stepChip(
                    step2Done,
                    step1Done && !step2Done,
                    !step1Done,
                  )}
                >
                  {step2Done
                    ? "Готово"
                    : step1Done
                      ? "Нужно действие"
                      : "Сначала регистрация"}
                </span>
              </div>
              {step2Done && status.team ? (
                <>
                  <p>
                    Команда «<strong>{status.team.name}</strong>» ·{" "}
                    {status.team.members.length} из {hackathon.max_team_size}{" "}
                    мест
                  </p>
                  {!status.team.case_id && (
                    <p className="step-warning">
                      Кейс не выбран - капитан должен указать трек и кейс перед
                      сдачей.
                    </p>
                  )}
                </>
              ) : (
                <p>
                  Создайте свою команду или вступите в открытую - до{" "}
                  {hackathon.max_team_size} человек.
                </p>
              )}
              <Link
                to={`/hackathons/${id}/team`}
                className={`btn-secondary ${!step1Done ? "btn-disabled" : ""}`}
                aria-disabled={!step1Done}
                onClick={(e) => {
                  if (!step1Done) {
                    e.preventDefault();
                  }
                }}
              >
                {step2Done ? "Управление командой" : "Собрать команду"}
              </Link>
            </div>
          </section>
        </Reveal>

        <Reveal delay={240}>
          <section
            className={`card glass participate-step ${step3Done ? "step-done" : ""} ${!step2Done ? "step-disabled" : ""}`}
          >
            <div className="step-indicator">{step3Done ? "✓" : "3"}</div>
            <div className="step-body">
              <div className="step-title-row">
                <h2>Сдача решения</h2>
                <span
                  className={stepChip(
                    step3Done,
                    step2Done && !step3Done,
                    !step2Done,
                  )}
                >
                  {step3Done
                    ? "Отправлено"
                    : step2Done
                      ? status.can_submit
                        ? "Можно сдавать"
                        : "Требуется настройка"
                      : "Сначала команда"}
                </span>
              </div>
              <p>{getStep3Hint(status)}</p>
              <Link
                to={`/hackathons/${id}/submission`}
                className={`btn-primary ${!step2Done ? "btn-disabled" : ""}`}
                aria-disabled={!step2Done}
                onClick={(e) => {
                  if (!step2Done) {
                    e.preventDefault();
                  }
                }}
              >
                {step3Done ? "Редактировать сдачу" : "Перейти к сдаче"}
              </Link>
            </div>
          </section>
        </Reveal>
      </div>

      <Reveal delay={320}>
        <div className="participate-footer">
          <Link to={`/hackathons/${id}`} className="btn-ghost">
            К карточке хакатона
          </Link>
        </div>
      </Reveal>
    </div>
  );
}
