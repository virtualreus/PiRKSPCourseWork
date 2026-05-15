import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import * as hackathonsApi from "../api/hackathons";
import type { HackathonDetail } from "../api/hackathonTypes";
import { ApiError } from "../api/client";
import { useAuth } from "../context/AuthContext";
import { formatDate, statusLabel } from "../utils/hackathon";

export function HackathonDetailPage() {
  const { id } = useParams();
  const { user } = useAuth();
  const [item, setItem] = useState<HackathonDetail | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) {
      return;
    }
    (async () => {
      try {
        const data = await hackathonsApi.getHackathon(id);
        setItem(data);
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

  if (loading) {
    return (
      <div className="page-loading" aria-busy>
        <div className="spinner" />
      </div>
    );
  }

  if (error || !item) {
    return (
      <div className="card glass form-card">
        <p className="form-error">{error || "Хакатон не найден"}</p>
        <Link to="/">На главную</Link>
      </div>
    );
  }

  const isOwner = user?.platform_role === "organizer";

  return (
    <article className="hackathon-detail">
      <div className="detail-header">
        <span className={`status-badge status-${item.status}`}>
          {statusLabel(item.status)}
        </span>
        {item.format && <span className="format-badge">{item.format}</span>}
      </div>
      <h1>{item.title}</h1>
      {item.short_description && (
        <p className="lead">{item.short_description}</p>
      )}
      <p>{item.description}</p>

      {item.prizes_info && (
        <section className="card glass detail-block">
          <h2>Призы</h2>
          <p>{item.prizes_info}</p>
        </section>
      )}

      <section className="card glass detail-block">
        <h2>Таймлайн</h2>
        <dl className="timeline-list">
          <div>
            <dt>Регистрация открыта</dt>
            <dd>{formatDate(item.timeline.registration_opens_at)}</dd>
          </div>
          <div>
            <dt>Регистрация закрыта</dt>
            <dd>{formatDate(item.timeline.registration_closes_at)}</dd>
          </div>
          <div>
            <dt>Старт</dt>
            <dd>{formatDate(item.timeline.event_starts_at)}</dd>
          </div>
          <div>
            <dt>Конец кодинга</dt>
            <dd>{formatDate(item.timeline.event_ends_at)}</dd>
          </div>
          <div>
            <dt>Дедлайн сдачи</dt>
            <dd>{formatDate(item.timeline.submission_deadline_at)}</dd>
          </div>
        </dl>
      </section>

      <section className="card glass detail-block">
        <h2>Треки и кейсы</h2>
        {item.tracks.map((track) => (
          <div key={track.id} className="track-block">
            <h3>{track.title}</h3>
            {track.description && <p>{track.description}</p>}
            <ul className="case-list">
              {track.cases?.map((c) => (
                <li key={c.id}>
                  <strong>{c.title}</strong>
                  {c.customer_name && <span> - {c.customer_name}</span>}
                  {c.description && <p>{c.description}</p>}
                  {c.resources_url && (
                    <a href={c.resources_url} target="_blank" rel="noreferrer">
                      Ресурсы
                    </a>
                  )}
                </li>
              ))}
            </ul>
          </div>
        ))}
      </section>

      <div className="detail-actions">
        {!user && (
          <Link to="/login" className="btn-primary">
            Войти, чтобы участвовать
          </Link>
        )}
        {isOwner && (
          <Link
            to={`/organizer/hackathons/${item.id}`}
            className="btn-secondary"
          >
            Управление
          </Link>
        )}
        <Link to="/" className="btn-ghost">
          К каталогу
        </Link>
      </div>
    </article>
  );
}
