import { useEffect, useMemo, useState, type FormEvent } from 'react';
import { Link } from 'react-router-dom';

import { ApiError } from '../api/client';
import * as authApi from '../api/auth';
import * as profileApi from '../api/profile';
import type { TeamMemberRole } from '../api/participationTypes';
import type { UserDashboard, UserParticipation } from '../api/profileTypes';
import { ParticipationAlert } from '../components/ParticipationAlert';
import { Reveal } from '../components/Reveal';
import { useAuth } from '../context/AuthContext';
import { formatDate, statusLabel } from '../utils/hackathon';
import { formatDeadlineRemaining, getSubmitBlockInfo } from '../utils/participation';
import { teamRoleLabel } from '../utils/team';

function userInitials(name: string): string {
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((p) => p[0]?.toUpperCase() ?? '')
    .join('');
}

function participationStatusLabel(p: UserParticipation): string {
  if (p.submission?.submitted_at) {
    return 'Сдача отправлена';
  }
  if (p.team) {
    return p.can_submit ? 'Готово к сдаче' : 'Команда собрана';
  }
  return 'Зарегистрирован';
}

function participationStatusClass(p: UserParticipation): string {
  if (p.submission?.submitted_at) {
    return 'profile-status-done';
  }
  if (p.can_submit) {
    return 'profile-status-ready';
  }
  if (p.team) {
    return 'profile-status-progress';
  }
  return 'profile-status-muted';
}

export function ProfilePage() {
  const { user, refreshUser } = useAuth();
  const [dashboard, setDashboard] = useState<UserDashboard | null>(null);
  const [fullName, setFullName] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (user) {
      setFullName(user.full_name);
    }
  }, [user]);

  useEffect(() => {
    (async () => {
      try {
        const data = await profileApi.getDashboard();
        setDashboard(data);
      } catch (err) {
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError('Не удалось загрузить профиль');
        }
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const activeParticipation = useMemo(
    () =>
      dashboard?.participations.find(
        (p) => p.hackathon.status === 'registration' || p.hackathon.status === 'running',
      ) ?? null,
    [dashboard],
  );

  const pastParticipations = useMemo(
    () =>
      dashboard?.participations.filter(
        (p) => p.hackathon.status !== 'registration' && p.hackathon.status !== 'running',
      ) ?? [],
    [dashboard],
  );

  if (!user) {
    return null;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setMessage('');
    setError('');
    setSaving(true);
    try {
      await authApi.updateMe({ full_name: fullName });
      await refreshUser();
      const data = await profileApi.getDashboard();
      setDashboard(data);
      setMessage('Профиль обновлён');
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Не удалось сохранить');
      }
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="page-loading profile-page" aria-busy>
        <div className="spinner" />
      </div>
    );
  }

  const stats = dashboard?.stats;
  const isOrganizer = user.platform_role === 'organizer';

  return (
    <div className="profile-page profile-page-wide">
      <Reveal>
        <header className="profile-hero card glass">
          <div className="profile-avatar" aria-hidden>
            {userInitials(dashboard?.user.full_name ?? user.full_name)}
          </div>
          <div className="profile-hero-body">
            <h1>{dashboard?.user.full_name ?? user.full_name}</h1>
            <p className="profile-email">{user.email}</p>
            <div className="profile-badges">
              <span className={`role-badge role-${user.platform_role}`}>
                {isOrganizer ? 'Организатор' : 'Участник'}
              </span>
              <span className="profile-since">
                На платформе с {formatDate(user.created_at)}
              </span>
            </div>
          </div>
          <Link to="/" className="btn-ghost profile-catalog-link">
            Каталог хакатонов
          </Link>
        </header>
      </Reveal>

      {error && !dashboard && <p className="form-error participate-banner">{error}</p>}

      {stats && (
        <Reveal delay={60}>
          <div className="profile-stats">
            <div className="stat-tile">
              <span className="stat-value">{stats.registrations_count}</span>
              <span className="stat-label">Регистраций</span>
            </div>
            <div className="stat-tile">
              <span className="stat-value">{stats.active_hackathons}</span>
              <span className="stat-label">Активных хакатонов</span>
            </div>
            <div className="stat-tile">
              <span className="stat-value">{stats.teams_count}</span>
              <span className="stat-label">Команд</span>
            </div>
            <div className="stat-tile">
              <span className="stat-value">{stats.submissions_count}</span>
              <span className="stat-label">Сдач</span>
            </div>
            {isOrganizer && stats.organized_count != null && (
              <div className="stat-tile stat-tile-accent">
                <span className="stat-value">{stats.organized_count}</span>
                <span className="stat-label">Организую</span>
              </div>
            )}
          </div>
        </Reveal>
      )}

      <div className="profile-layout">
        <div className="profile-main">
          {activeParticipation && (
            <Reveal delay={100}>
              <section className="card glass profile-spotlight">
                <div className="profile-spotlight-header">
                  <span className="spotlight-label">Текущий хакатон</span>
                  <span className={`status-badge status-${activeParticipation.hackathon.status}`}>
                    {statusLabel(activeParticipation.hackathon.status)}
                  </span>
                </div>
                <h2>
                  <Link to={`/hackathons/${activeParticipation.hackathon.id}`}>
                    {activeParticipation.hackathon.title}
                  </Link>
                </h2>
                {activeParticipation.hackathon.short_description && (
                  <p className="profile-spotlight-desc">
                    {activeParticipation.hackathon.short_description}
                  </p>
                )}

                <div className="profile-spotlight-grid">
                  <div>
                    <span className="meta-label">Регистрация</span>
                    <strong>{formatDate(activeParticipation.registered_at)}</strong>
                  </div>
                  <div>
                    <span className="meta-label">Дедлайн сдачи</span>
                    <strong>
                      {formatDate(activeParticipation.hackathon.submission_deadline_at)}
                    </strong>
                    {formatDeadlineRemaining(
                      activeParticipation.hackathon.submission_deadline_at,
                    ) && (
                      <span className="meta-sub">
                        {formatDeadlineRemaining(
                          activeParticipation.hackathon.submission_deadline_at,
                        )}
                      </span>
                    )}
                  </div>
                  <div>
                    <span className="meta-label">Команда</span>
                    <strong>
                      {activeParticipation.team?.name ?? 'Не в команде'}
                    </strong>
                    {activeParticipation.team && (
                      <span className="meta-sub">
                        {activeParticipation.team.member_count} участник
                        {activeParticipation.team.is_captain ? ' · вы капитан' : ''}
                        {activeParticipation.team.team_role
                          ? ` · ${teamRoleLabel(activeParticipation.team.team_role as TeamMemberRole)}`
                          : ''}
                      </span>
                    )}
                  </div>
                  <div>
                    <span className="meta-label">Кейс</span>
                    <strong>
                      {activeParticipation.team?.case_title ?? 'Не выбран'}
                    </strong>
                    {activeParticipation.team?.track_title && (
                      <span className="meta-sub">{activeParticipation.team.track_title}</span>
                    )}
                  </div>
                </div>

                {!activeParticipation.can_submit &&
                  activeParticipation.submit_block_reason &&
                  activeParticipation.team && (
                    <ParticipationAlert
                      variant="warning"
                      title={
                        getSubmitBlockInfo(
                          activeParticipation.submit_block_reason,
                          activeParticipation.hackathon.id,
                          activeParticipation.team.is_captain,
                        )?.title ?? 'Действие требуется'
                      }
                      actionLabel="Открыть участие"
                      actionTo={`/hackathons/${activeParticipation.hackathon.id}/participate`}
                    >
                      {
                        getSubmitBlockInfo(
                          activeParticipation.submit_block_reason,
                          activeParticipation.hackathon.id,
                          activeParticipation.team.is_captain,
                        )?.body
                      }
                    </ParticipationAlert>
                  )}

                <div className="profile-spotlight-actions">
                  <Link
                    to={`/hackathons/${activeParticipation.hackathon.id}/participate`}
                    className="btn-primary"
                  >
                    Участие
                  </Link>
                  {activeParticipation.team && (
                    <>
                      <Link
                        to={`/hackathons/${activeParticipation.hackathon.id}/team`}
                        className="btn-secondary"
                      >
                        Команда
                      </Link>
                      <Link
                        to={`/hackathons/${activeParticipation.hackathon.id}/submission`}
                        className="btn-secondary"
                      >
                        Сдача
                      </Link>
                    </>
                  )}
                </div>
              </section>
            </Reveal>
          )}

          <Reveal delay={activeParticipation ? 160 : 100}>
            <section className="profile-section">
              <div className="profile-section-head">
                <h2>Мои хакатоны</h2>
                <span className="muted">
                  {dashboard?.participations.length ?? 0} всего
                </span>
              </div>

              {!dashboard?.participations.length ? (
                <div className="card glass profile-empty">
                  <p>Вы ещё не зарегистрированы ни на один хакатон.</p>
                  <Link to="/" className="btn-primary">
                    Смотреть каталог
                  </Link>
                </div>
              ) : (
                <ul className="profile-hackathon-list">
                  {dashboard.participations.map((p) => (
                    <li key={p.hackathon.id} className="card glass profile-hackathon-card">
                      <div className="profile-hackathon-top">
                        <div>
                          <Link to={`/hackathons/${p.hackathon.id}`} className="profile-hack-title">
                            {p.hackathon.title}
                          </Link>
                          <p className="muted profile-hack-meta">
                            {statusLabel(p.hackathon.status)}
                            {p.hackathon.format && ` · ${p.hackathon.format}`}
                          </p>
                        </div>
                        <span className={`profile-status-pill ${participationStatusClass(p)}`}>
                          {participationStatusLabel(p)}
                        </span>
                      </div>

                      <dl className="profile-hack-details">
                        <div>
                          <dt>Зарегистрирован</dt>
                          <dd>{formatDate(p.registered_at)}</dd>
                        </div>
                        <div>
                          <dt>Команда</dt>
                          <dd>{p.team?.name ?? '—'}</dd>
                        </div>
                        <div>
                          <dt>Кейс</dt>
                          <dd>{p.team?.case_title ?? '—'}</dd>
                        </div>
                        <div>
                          <dt>Сдача</dt>
                          <dd>
                            {p.submission?.submitted_at
                              ? formatDate(p.submission.submitted_at)
                              : p.submission
                                ? 'Черновик'
                                : '—'}
                          </dd>
                        </div>
                      </dl>

                      {p.submission?.repo_url && (
                        <a
                          href={p.submission.repo_url}
                          target="_blank"
                          rel="noreferrer"
                          className="profile-repo-link"
                        >
                          Репозиторий
                        </a>
                      )}

                      <div className="profile-hack-actions">
                        <Link
                          to={`/hackathons/${p.hackathon.id}/participate`}
                          className="btn-ghost btn-sm"
                        >
                          Участие
                        </Link>
                        {p.team && (
                          <Link
                            to={`/hackathons/${p.hackathon.id}/team`}
                            className="btn-ghost btn-sm"
                          >
                            Команда
                          </Link>
                        )}
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </section>
          </Reveal>

          {pastParticipations.length > 0 && activeParticipation && (
            <p className="muted profile-archive-note">
              Завершённые и прошлые события отображаются в списке выше.
            </p>
          )}

          {isOrganizer && (dashboard?.organized_hackathons?.length ?? 0) > 0 && (
            <Reveal delay={200}>
              <section className="profile-section">
                <div className="profile-section-head">
                  <h2>Организую</h2>
                  <Link to="/organizer/hackathons" className="btn-ghost btn-sm">
                    Все хакатоны
                  </Link>
                </div>
                <ul className="profile-organizer-list">
                  {dashboard!.organized_hackathons!.map((h) => (
                    <li key={h.id} className="card glass profile-organizer-card">
                      <Link to={`/organizer/hackathons/${h.id}`}>{h.title}</Link>
                      <span className={`status-badge status-${h.status}`}>
                        {statusLabel(h.status)}
                      </span>
                      <span className="muted">
                        дедлайн {formatDate(h.submission_deadline_at)}
                      </span>
                    </li>
                  ))}
                </ul>
              </section>
            </Reveal>
          )}
        </div>

        <aside className="profile-sidebar">
          <Reveal delay={120}>
            <form className="card glass form-card profile-edit-card" onSubmit={handleSubmit}>
              <h2>Редактирование</h2>
              <p className="form-hint">Имя отображается в команде и для организаторов.</p>

              {message && <p className="form-success">{message}</p>}
              {error && dashboard && <p className="form-error">{error}</p>}

              <label>
                ФИО
                <input
                  type="text"
                  required
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                />
              </label>

              <dl className="profile-sidebar-meta">
                <div>
                  <dt>Email</dt>
                  <dd>{user.email}</dd>
                </div>
                <div>
                  <dt>ID</dt>
                  <dd className="profile-id">{user.id}</dd>
                </div>
              </dl>

              <button type="submit" className="btn-primary btn-block" disabled={saving}>
                {saving ? 'Сохранение…' : 'Сохранить'}
              </button>
            </form>
          </Reveal>
        </aside>
      </div>
    </div>
  );
}
