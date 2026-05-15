import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import * as hackathonsApi from '../../api/hackathons';
import type { HackathonFormat } from '../../api/hackathonTypes';
import { ApiError } from '../../api/client';
import { defaultTimeline, toDatetimeLocal, toRFC3339 } from '../../utils/hackathon';

export function OrganizerHackathonNewPage() {
  const navigate = useNavigate();
  const timelineDefaults = defaultTimeline();

  const [title, setTitle] = useState('');
  const [shortDescription, setShortDescription] = useState('');
  const [description, setDescription] = useState('');
  const [format, setFormat] = useState<HackathonFormat>('online');
  const [maxTeamSize, setMaxTeamSize] = useState(5);
  const [prizesInfo, setPrizesInfo] = useState('');
  const [regOpens, setRegOpens] = useState(toDatetimeLocal(timelineDefaults.registration_opens_at));
  const [regCloses, setRegCloses] = useState(toDatetimeLocal(timelineDefaults.registration_closes_at));
  const [eventStarts, setEventStarts] = useState(toDatetimeLocal(timelineDefaults.event_starts_at));
  const [eventEnds, setEventEnds] = useState(toDatetimeLocal(timelineDefaults.event_ends_at));
  const [submissionDeadline, setSubmissionDeadline] = useState(
    toDatetimeLocal(timelineDefaults.submission_deadline_at),
  );
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      const created = await hackathonsApi.createHackathon({
        title,
        short_description: shortDescription,
        description,
        format,
        max_team_size: maxTeamSize,
        prizes_info: prizesInfo,
        timeline: {
          registration_opens_at: toRFC3339(regOpens),
          registration_closes_at: toRFC3339(regCloses),
          event_starts_at: toRFC3339(eventStarts),
          event_ends_at: toRFC3339(eventEnds),
          submission_deadline_at: toRFC3339(submissionDeadline),
        },
      });
      navigate(`/organizer/hackathons/${created.id}`);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Не удалось создать хакатон');
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="organizer-page">
      <h1>Новый хакатон</h1>
      <form className="card glass form-card wide-form" onSubmit={handleSubmit}>
        {error && <p className="form-error">{error}</p>}

        <label>
          Название
          <input required value={title} onChange={(e) => setTitle(e.target.value)} />
        </label>
        <label>
          Краткое описание
          <input value={shortDescription} onChange={(e) => setShortDescription(e.target.value)} />
        </label>
        <label>
          Описание
          <textarea required rows={4} value={description} onChange={(e) => setDescription(e.target.value)} />
        </label>
        <label>
          Формат
          <select value={format} onChange={(e) => setFormat(e.target.value as HackathonFormat)}>
            <option value="online">online</option>
            <option value="offline">offline</option>
            <option value="hybrid">hybrid</option>
          </select>
        </label>
        <label>
          Размер команды
          <input
            type="number"
            min={2}
            max={8}
            value={maxTeamSize}
            onChange={(e) => setMaxTeamSize(Number(e.target.value))}
          />
        </label>
        <label>
          Призы
          <textarea rows={2} value={prizesInfo} onChange={(e) => setPrizesInfo(e.target.value)} />
        </label>

        <h2 className="form-section-title">Даты</h2>
        <label>
          Открытие регистрации
          <input type="datetime-local" required value={regOpens} onChange={(e) => setRegOpens(e.target.value)} />
        </label>
        <label>
          Закрытие регистрации
          <input type="datetime-local" required value={regCloses} onChange={(e) => setRegCloses(e.target.value)} />
        </label>
        <label>
          Старт
          <input type="datetime-local" required value={eventStarts} onChange={(e) => setEventStarts(e.target.value)} />
        </label>
        <label>
          Конец кодинга
          <input type="datetime-local" required value={eventEnds} onChange={(e) => setEventEnds(e.target.value)} />
        </label>
        <label>
          Дедлайн сдачи
          <input
            type="datetime-local"
            required
            value={submissionDeadline}
            onChange={(e) => setSubmissionDeadline(e.target.value)}
          />
        </label>

        <div className="form-actions">
          <button type="submit" className="btn-primary" disabled={submitting}>
            {submitting ? 'Создание…' : 'Создать черновик'}
          </button>
          <Link to="/organizer/hackathons" className="btn-ghost">
            Отмена
          </Link>
        </div>
      </form>
    </section>
  );
}
