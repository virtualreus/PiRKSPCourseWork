import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import * as hackathonsApi from '../api/hackathons';
import type { HackathonDetail } from '../api/hackathonTypes';
import * as participationApi from '../api/participation';
import type { ParticipationStatus, Submission } from '../api/participationTypes';
import { ApiError } from '../api/client';
import { ParticipationAlert } from '../components/ParticipationAlert';
import { Reveal } from '../components/Reveal';
import { useAuth } from '../context/AuthContext';
import { formatDate, statusLabel } from '../utils/hackathon';
import {
  formatDeadlineRemaining,
  getSubmitBlockInfo,
} from '../utils/participation';

export function SubmissionPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [hackathon, setHackathon] = useState<HackathonDetail | null>(null);
  const [status, setStatus] = useState<ParticipationStatus | null>(null);
  const [submission, setSubmission] = useState<Submission | null>(null);
  const [title, setTitle] = useState('');
  const [summary, setSummary] = useState('');
  const [repoUrl, setRepoUrl] = useState('');
  const [demoUrl, setDemoUrl] = useState('');
  const [pitchUrl, setPitchUrl] = useState('');
  const [videoUrl, setVideoUrl] = useState('');
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
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

    if (p.team) {
      try {
        const sub = await participationApi.getTeamSubmission(p.team.id);
        setSubmission(sub);
        setTitle(sub.title ?? '');
        setSummary(sub.summary ?? '');
        setRepoUrl(sub.repo_url);
        setDemoUrl(sub.demo_url ?? '');
        setPitchUrl(sub.pitch_url ?? '');
        setVideoUrl(sub.video_url ?? '');
      } catch (err) {
        if (!(err instanceof ApiError && err.status === 404)) {
          throw err;
        }
        setSubmission(null);
      }
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
          setError('Не удалось загрузить сдачу');
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  useEffect(() => {
    if (!loading && status && !status.team && id) {
      navigate(`/hackathons/${id}/team`, { replace: true });
    }
  }, [status, loading, id, navigate]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!status?.team || !status.can_submit) {
      return;
    }
    setBusy(true);
    setError('');
    setMessage('');
    try {
      const saved = await participationApi.upsertTeamSubmission(status.team.id, {
        title: title.trim() || undefined,
        summary: summary.trim() || undefined,
        repo_url: repoUrl.trim(),
        demo_url: demoUrl.trim() || undefined,
        pitch_url: pitchUrl.trim() || undefined,
        video_url: videoUrl.trim() || undefined,
      });
      setSubmission(saved);
      setMessage(submission?.submitted_at ? 'Сдача обновлена' : 'Решение отправлено');
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

  if (!hackathon || !status?.team || !id) {
    return null;
  }

  const isCaptain = user?.id === status.team.captain_id;
  const canEdit = status.can_submit;
  const blockInfo = getSubmitBlockInfo(status.submit_block_reason, id, isCaptain);
  const deadlineLeft = status.submission_deadline_at
    ? formatDeadlineRemaining(status.submission_deadline_at)
    : null;

  const selectedTrack = hackathon.tracks.find((t) => t.id === status.team?.track_id);
  const selectedCase = selectedTrack?.cases?.find((c) => c.id === status.team?.case_id);

  return (
    <div className="participate-page submission-page participate-page-wide">
      <Reveal>
        <nav className="breadcrumb">
          <Link to={`/hackathons/${id}`}>{hackathon.title}</Link>
          <Link to={`/hackathons/${id}/participate`}>Участие</Link>
          <span>Сдача</span>
        </nav>
        <header className="participate-hero">
          <div>
            <h1 className="participate-title">Сдача решения</h1>
            <p className="participate-lead">
              Команда «{status.team.name}» · статус хакатона:{' '}
              {statusLabel(status.hackathon_status)}
            </p>
          </div>
          {status.submission_deadline_at && (
            <div className="deadline-card">
              <span className="deadline-label">Дедлайн сдачи</span>
              <strong>{formatDate(status.submission_deadline_at)}</strong>
              {deadlineLeft && canEdit && (
                <span className="deadline-remaining">{deadlineLeft}</span>
              )}
            </div>
          )}
        </header>
      </Reveal>

      {error && <p className="form-error participate-banner">{error}</p>}
      {message && <p className="form-success participate-banner">{message}</p>}

      <Reveal delay={60}>
        <div className="submission-meta-grid">
          <div className="meta-tile">
            <span className="meta-label">Кейс</span>
            <strong>{selectedCase?.title ?? 'Не выбран'}</strong>
            {selectedTrack && <span className="meta-sub">{selectedTrack.title}</span>}
          </div>
          <div className="meta-tile">
            <span className="meta-label">Статус сдачи</span>
            <strong>{submission?.submitted_at ? 'Отправлено' : 'Черновик'}</strong>
            {submission?.submitted_at && (
              <span className="meta-sub">{formatDate(submission.submitted_at)}</span>
            )}
          </div>
          <div className="meta-tile">
            <span className="meta-label">Участников</span>
            <strong>{status.team.members.length}</strong>
          </div>
        </div>
      </Reveal>

      {canEdit ? (
        <Reveal delay={100}>
          <ParticipationAlert variant="success" title="Можно отправлять решение">
            Заполните обязательную ссылку на репозиторий и при необходимости добавьте демо,
            презентацию и видео. После первого сохранения фиксируется время первичной сдачи.
          </ParticipationAlert>
        </Reveal>
      ) : (
        blockInfo && (
          <Reveal delay={100}>
            <ParticipationAlert
              variant={blockInfo.variant}
              title={blockInfo.title}
              actionLabel={blockInfo.actionLabel}
              actionTo={blockInfo.actionTo}
            >
              {blockInfo.body}
            </ParticipationAlert>
          </Reveal>
        )
      )}

      <Reveal delay={140}>
        <form className="card glass form-card wide-form submission-form" onSubmit={handleSubmit}>
          <section className="form-section">
            <h2 className="form-section-title">О проекте</h2>
            <p className="form-section-hint">
              Название и описание увидят организаторы в списке сдач.
            </p>
            <label>
              Название проекта
              <input
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                disabled={!canEdit || busy}
                placeholder="Например, SmartCity Dashboard"
              />
            </label>
            <label>
              Краткое описание
              <textarea
                rows={4}
                value={summary}
                onChange={(e) => setSummary(e.target.value)}
                disabled={!canEdit || busy}
                placeholder="Суть решения, стек, ключевая ценность для заказчика"
              />
            </label>
          </section>

          <section className="form-section">
            <h2 className="form-section-title">Артефакты</h2>
            <p className="form-section-hint">
              Репозиторий обязателен. Остальные ссылки помогут жюри быстрее оценить работу.
            </p>
            <label>
              Репозиторий <span className="required-mark">*</span>
              <input
                type="url"
                required
                value={repoUrl}
                onChange={(e) => setRepoUrl(e.target.value)}
                disabled={!canEdit || busy}
                placeholder="https://github.com/org/project"
              />
            </label>
            <label>
              Демо / стенд
              <input
                type="url"
                value={demoUrl}
                onChange={(e) => setDemoUrl(e.target.value)}
                disabled={!canEdit || busy}
                placeholder="https://demo.example.com"
              />
            </label>
            <label>
              Презентация (PDF / Slides)
              <input
                type="url"
                value={pitchUrl}
                onChange={(e) => setPitchUrl(e.target.value)}
                disabled={!canEdit || busy}
                placeholder="https://..."
              />
            </label>
            <label>
              Видео питча
              <input
                type="url"
                value={videoUrl}
                onChange={(e) => setVideoUrl(e.target.value)}
                disabled={!canEdit || busy}
                placeholder="https://..."
              />
            </label>
          </section>

          <div className="form-actions-row">
            <button type="submit" className="btn-primary" disabled={!canEdit || busy}>
              {busy ? 'Сохранение…' : submission?.submitted_at ? 'Обновить сдачу' : 'Отправить решение'}
            </button>
            <Link to={`/hackathons/${id}/participate`} className="btn-ghost">
              К участию
            </Link>
            <Link to={`/hackathons/${id}/team`} className="btn-ghost">
              К команде
            </Link>
          </div>
        </form>
      </Reveal>
    </div>
  );
}
