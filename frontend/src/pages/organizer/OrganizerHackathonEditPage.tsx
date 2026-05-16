import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import * as hackathonsApi from "../../api/hackathons";
import type { HackathonDetail } from "../../api/hackathonTypes";
import * as participationApi from "../../api/participation";
import type {
  HackathonRegistrationWithUser,
  SubmissionWithTeam,
} from "../../api/participationTypes";
import { ApiError } from "../../api/client";
import { TrackCaseManager } from "../../components/organizer/TrackCaseManager";
import { Reveal } from "../../components/Reveal";
import { formatDate, statusLabel } from "../../utils/hackathon";

export function OrganizerHackathonEditPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [item, setItem] = useState<HackathonDetail | null>(null);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [registrations, setRegistrations] = useState<HackathonRegistrationWithUser[]>([]);
  const [submissions, setSubmissions] = useState<SubmissionWithTeam[]>([]);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);

  async function reload() {
    if (!id) {
      return;
    }
    const data = await hackathonsApi.getOrganizerHackathon(id);
    setItem(data);
    if (data.status !== "draft") {
      const [regs, subs] = await Promise.all([
        participationApi.listOrganizerRegistrations(id),
        participationApi.listOrganizerSubmissions(id),
      ]);
      setRegistrations(regs.items);
      setSubmissions(subs.items);
    } else {
      setRegistrations([]);
      setSubmissions([]);
    }
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
          setError("Не удалось загрузить хакатон");
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  async function handlePublish() {
    if (!id) {
      return;
    }
    setError("");
    setMessage("");
    setBusy(true);
    try {
      await hackathonsApi.publishHackathon(id);
      setMessage("Хакатон опубликован");
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("Не удалось опубликовать");
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleDelete() {
    if (!id || !confirm("Удалить черновик?")) {
      return;
    }
    setBusy(true);
    try {
      await hackathonsApi.deleteHackathon(id);
      navigate("/organizer/hackathons");
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
      <div className="page-loading organizer-page" aria-busy>
        <div className="spinner" />
      </div>
    );
  }

  if (!item) {
    return (
      <div className="card glass form-card organizer-page">
        <p className="form-error">{error || "Не найдено"}</p>
        <Link to="/organizer/hackathons">К списку</Link>
      </div>
    );
  }

  const isDraft = item.status === "draft";
  const canPublish =
    isDraft && item.tracks.length > 0 && item.tracks.some((t) => (t.cases?.length ?? 0) > 0);

  return (
    <section className="organizer-page organizer-page-wide">
      <Reveal>
        <div className="page-toolbar">
          <h1>{item.title}</h1>
          <span className={`status-badge status-${item.status}`}>
            {statusLabel(item.status)}
          </span>
        </div>
      </Reveal>

      {message && <p className="form-success participate-banner">{message}</p>}
      {error && <p className="form-error participate-banner">{error}</p>}

      <div className="detail-actions">
        {isDraft && (
          <>
            <button
              type="button"
              className="btn-primary"
              disabled={busy || !canPublish}
              title={
                canPublish
                  ? undefined
                  : "Добавьте хотя бы один трек с кейсом"
              }
              onClick={handlePublish}
            >
              Опубликовать
            </button>
            <button
              type="button"
              className="btn-ghost"
              disabled={busy}
              onClick={handleDelete}
            >
              Удалить черновик
            </button>
          </>
        )}
        <Link to={`/hackathons/${item.id}`} className="btn-secondary">
          Публичная карточка
        </Link>
        <Link to="/organizer/hackathons" className="btn-ghost">
          К списку
        </Link>
      </div>

      <Reveal delay={80}>
        <section className="card glass detail-block">
          <h2>Треки и кейсы ({item.tracks.length})</h2>
          {isDraft && !canPublish && (
            <p className="form-hint organizer-publish-hint">
              Перед публикацией нужен минимум один трек и один кейс.
            </p>
          )}
          {id && (
            <TrackCaseManager
              hackathonId={id}
              tracks={item.tracks}
              isDraft={isDraft}
              busy={busy}
              setBusy={setBusy}
              onChanged={reload}
              onError={setError}
            />
          )}
        </section>
      </Reveal>

      {!isDraft && (
        <>
          <Reveal delay={120}>
            <section className="card glass detail-block">
              <h2>Участники ({registrations.length})</h2>
              {registrations.length === 0 ? (
                <div className="organizer-empty-block">
                  <p className="muted">Пока никто не зарегистрировался.</p>
                  <p className="form-hint">
                    Поделитесь ссылкой на карточку хакатона с участниками.
                  </p>
                </div>
              ) : (
                <ul className="organizer-list">
                  {registrations.map((r) => (
                    <li key={r.id}>
                      <strong>{r.user.full_name}</strong>
                      <span>{r.user.email}</span>
                      <time>{formatDate(r.registered_at)}</time>
                    </li>
                  ))}
                </ul>
              )}
            </section>
          </Reveal>

          <Reveal delay={160}>
            <section className="card glass detail-block">
              <h2>Сдачи ({submissions.length})</h2>
              {submissions.length === 0 ? (
                <div className="organizer-empty-block">
                  <p className="muted">Сдач пока нет.</p>
                  <p className="form-hint">
                    Команды появятся здесь после выбора кейса и отправки решения.
                  </p>
                </div>
              ) : (
                <ul className="organizer-submissions">
                  {submissions.map((s) => (
                    <li key={s.id} className="submission-row">
                      <div>
                        <strong>{s.team_name}</strong>
                        {s.title && <span> — {s.title}</span>}
                        {(s.track_title || s.case_title) && (
                          <p className="muted">
                            {[s.track_title, s.case_title].filter(Boolean).join(" · ")}
                          </p>
                        )}
                        {s.submitted_at && (
                          <p className="meta-sub">Сдано {formatDate(s.submitted_at)}</p>
                        )}
                      </div>
                      <div className="submission-links">
                        <a href={s.repo_url} target="_blank" rel="noreferrer">
                          Репозиторий
                        </a>
                        {s.demo_url && (
                          <a href={s.demo_url} target="_blank" rel="noreferrer">
                            Демо
                          </a>
                        )}
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </section>
          </Reveal>
        </>
      )}

      <Reveal delay={200}>
        <section className="card glass detail-block muted-block">
          <h2>Таймлайн</h2>
          <dl className="timeline-list">
            <div>
              <dt>Регистрация</dt>
              <dd>
                {formatDate(item.timeline.registration_opens_at)} —{" "}
                {formatDate(item.timeline.registration_closes_at)}
              </dd>
            </div>
            <div>
              <dt>Событие</dt>
              <dd>
                {formatDate(item.timeline.event_starts_at)} —{" "}
                {formatDate(item.timeline.event_ends_at)}
              </dd>
            </div>
            <div>
              <dt>Дедлайн сдачи</dt>
              <dd>{formatDate(item.timeline.submission_deadline_at)}</dd>
            </div>
          </dl>
        </section>
      </Reveal>
    </section>
  );
}
