import { useEffect, useState, type FormEvent } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import * as hackathonsApi from "../../api/hackathons";
import type { HackathonDetail } from "../../api/hackathonTypes";
import { ApiError } from "../../api/client";
import { formatDate, statusLabel } from "../../utils/hackathon";

export function OrganizerHackathonEditPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [item, setItem] = useState<HackathonDetail | null>(null);
  const [trackTitle, setTrackTitle] = useState("");
  const [trackDescription, setTrackDescription] = useState("");
  const [caseTitle, setCaseTitle] = useState("");
  const [caseDescription, setCaseDescription] = useState("");
  const [caseCustomer, setCaseCustomer] = useState("");
  const [caseResources, setCaseResources] = useState("");
  const [selectedTrackId, setSelectedTrackId] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);

  async function reload() {
    if (!id) {
      return;
    }
    const data = await hackathonsApi.getOrganizerHackathon(id);
    setItem(data);
    if (!selectedTrackId && data.tracks.length > 0) {
      setSelectedTrackId(data.tracks[0].id);
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

  async function handleAddTrack(e: FormEvent) {
    e.preventDefault();
    if (!id) {
      return;
    }
    setBusy(true);
    setError("");
    try {
      await hackathonsApi.createTrack(id, {
        title: trackTitle,
        description: trackDescription,
      });
      setTrackTitle("");
      setTrackDescription("");
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleAddCase(e: FormEvent) {
    e.preventDefault();
    if (!selectedTrackId) {
      return;
    }
    setBusy(true);
    setError("");
    try {
      await hackathonsApi.createCase(selectedTrackId, {
        title: caseTitle,
        description: caseDescription,
        customer_name: caseCustomer,
        resources_url: caseResources,
      });
      setCaseTitle("");
      setCaseDescription("");
      setCaseCustomer("");
      setCaseResources("");
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
    return <p className="page-loading">Загрузка…</p>;
  }

  if (!item) {
    return <p className="form-error">{error || "Не найдено"}</p>;
  }

  const isDraft = item.status === "draft";

  return (
    <section className="organizer-page">
      <div className="page-toolbar">
        <h1>{item.title}</h1>
        <span className={`status-badge status-${item.status}`}>
          {statusLabel(item.status)}
        </span>
      </div>

      {message && <p className="form-success">{message}</p>}
      {error && <p className="form-error">{error}</p>}

      <div className="detail-actions">
        {isDraft && (
          <>
            <button
              type="button"
              className="btn-primary"
              disabled={busy}
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

      <section className="card glass detail-block">
        <h2>Треки ({item.tracks.length})</h2>
        {item.tracks.map((track) => (
          <div key={track.id} className="track-block">
            <h3>{track.title}</h3>
            {track.description && <p>{track.description}</p>}
            <ul className="case-list">
              {(track.cases ?? []).map((c) => (
                <li key={c.id}>
                  <strong>{c.title}</strong> - {c.description}
                </li>
              ))}
            </ul>
          </div>
        ))}

        {isDraft && (
          <form className="inline-form" onSubmit={handleAddTrack}>
            <h3>Добавить трек</h3>
            <label>
              Название
              <input
                required
                value={trackTitle}
                onChange={(e) => setTrackTitle(e.target.value)}
              />
            </label>
            <label>
              Описание
              <input
                value={trackDescription}
                onChange={(e) => setTrackDescription(e.target.value)}
              />
            </label>
            <button type="submit" className="btn-secondary" disabled={busy}>
              Добавить трек
            </button>
          </form>
        )}
      </section>

      {isDraft && item.tracks.length > 0 && (
        <section className="card glass detail-block">
          <form className="inline-form" onSubmit={handleAddCase}>
            <h2>Добавить кейс</h2>
            <label>
              Трек
              <select
                value={selectedTrackId}
                onChange={(e) => setSelectedTrackId(e.target.value)}
              >
                {item.tracks.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.title}
                  </option>
                ))}
              </select>
            </label>
            <label>
              Название кейса
              <input
                required
                value={caseTitle}
                onChange={(e) => setCaseTitle(e.target.value)}
              />
            </label>
            <label>
              Описание
              <textarea
                required
                rows={3}
                value={caseDescription}
                onChange={(e) => setCaseDescription(e.target.value)}
              />
            </label>
            <label>
              Заказчик
              <input
                value={caseCustomer}
                onChange={(e) => setCaseCustomer(e.target.value)}
              />
            </label>
            <label>
              Ссылка на ресурсы
              <input
                value={caseResources}
                onChange={(e) => setCaseResources(e.target.value)}
              />
            </label>
            <button type="submit" className="btn-secondary" disabled={busy}>
              Добавить кейс
            </button>
          </form>
        </section>
      )}

      <section className="card glass detail-block muted-block">
        <h2>Таймлайн</h2>
        <p>
          Регистрация: {formatDate(item.timeline.registration_opens_at)} -{" "}
          {formatDate(item.timeline.registration_closes_at)}
        </p>
        <p>Дедлайн сдачи: {formatDate(item.timeline.submission_deadline_at)}</p>
      </section>
    </section>
  );
}
