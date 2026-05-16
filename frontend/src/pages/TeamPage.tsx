import { useEffect, useMemo, useState, type FormEvent } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';

import * as hackathonsApi from '../api/hackathons';
import type { HackathonDetail } from '../api/hackathonTypes';
import * as participationApi from '../api/participation';
import type { ParticipationStatus, Team, TeamMemberRole } from '../api/participationTypes';
import { ApiError } from '../api/client';
import { ParticipationAlert } from '../components/ParticipationAlert';
import { Reveal } from '../components/Reveal';
import { useAuth } from '../context/AuthContext';
import { teamRoleLabel } from '../utils/team';

export function TeamPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [hackathon, setHackathon] = useState<HackathonDetail | null>(null);
  const [status, setStatus] = useState<ParticipationStatus | null>(null);
  const [teams, setTeams] = useState<Team[]>([]);
  const [teamName, setTeamName] = useState('');
  const [trackId, setTrackId] = useState('');
  const [caseId, setCaseId] = useState('');
  const [joinRole, setJoinRole] = useState<TeamMemberRole>('developer');
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);

  const myTeam = status?.team ?? null;
  const isCaptain = myTeam && user?.id === myTeam.captain_id;
  const selectedTrack = hackathon?.tracks.find((t) => t.id === myTeam?.track_id);
  const selectedCase = selectedTrack?.cases?.find((c) => c.id === myTeam?.case_id);
  const needsCase = Boolean(myTeam && !myTeam.case_id && hackathon && hackathon.tracks.length > 0);

  const casesForTrack = useMemo(() => {
    if (!hackathon || !trackId) {
      return [];
    }
    const track = hackathon.tracks.find((t) => t.id === trackId);
    return track?.cases ?? [];
  }, [hackathon, trackId]);

  async function reload() {
    if (!id) {
      return;
    }
    const [h, p, list] = await Promise.all([
      hackathonsApi.getHackathon(id),
      participationApi.getParticipation(id),
      participationApi.listTeams(id),
    ]);
    setHackathon(h);
    setStatus(p);
    setTeams(list.items);
    if (p.team?.track_id) {
      setTrackId(p.team.track_id);
    }
    if (p.team?.case_id) {
      setCaseId(p.team.case_id);
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
          setError('Не удалось загрузить данные');
        }
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  useEffect(() => {
    if (!status?.is_registered && !loading && id) {
      navigate(`/hackathons/${id}/participate`, { replace: true });
    }
  }, [status, loading, id, navigate]);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    if (!id) {
      return;
    }
    setBusy(true);
    setError('');
    setMessage('');
    try {
      await participationApi.createTeam(id, {
        name: teamName.trim(),
        track_id: trackId || undefined,
        case_id: caseId || undefined,
        team_role: 'team_lead',
      });
      setMessage('Команда создана');
      setTeamName('');
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleJoin(teamId: string) {
    setBusy(true);
    setError('');
    setMessage('');
    try {
      await participationApi.joinTeam(teamId, { team_role: joinRole });
      setMessage('Вы вступили в команду');
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleLeave() {
    if (!myTeam || !confirm('Выйти из команды?')) {
      return;
    }
    setBusy(true);
    setError('');
    try {
      await participationApi.leaveTeam(myTeam.id);
      setMessage('Вы вышли из команды');
      await reload();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleUpdateTrackCase(e: FormEvent) {
    e.preventDefault();
    if (!myTeam) {
      return;
    }
    setBusy(true);
    setError('');
    setMessage('');
    try {
      await participationApi.updateTeam(myTeam.id, {
        track_id: trackId || undefined,
        case_id: caseId || undefined,
      });
      setMessage('Трек и кейс обновлены');
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

  if (!hackathon || !id) {
    return null;
  }

  return (
    <div className="participate-page team-page participate-page-wide">
      <Reveal>
        <nav className="breadcrumb">
          <Link to={`/hackathons/${id}`}>{hackathon.title}</Link>
          <Link to={`/hackathons/${id}/participate`}>Участие</Link>
          <span>Команда</span>
        </nav>
        <header className="participate-hero">
          <div>
            <h1 className="participate-title">Команда</h1>
            <p className="participate-lead">
              До {hackathon.max_team_size} человек · кейс обязателен перед сдачей решения
            </p>
          </div>
        </header>
      </Reveal>

      {error && <p className="form-error participate-banner">{error}</p>}
      {message && <p className="form-success participate-banner">{message}</p>}

      {myTeam ? (
        <Reveal delay={80}>
          {needsCase && (
            <ParticipationAlert
              variant="warning"
              title={isCaptain ? 'Выберите кейс для сдачи' : 'Ожидается выбор кейса'}
            >
              {isCaptain
                ? 'Без выбранного трека и кейса форма сдачи будет заблокирована. Укажите их ниже и нажмите «Сохранить выбор».'
                : 'Капитан команды должен выбрать трек и кейс — после этого откроется сдача решения.'}
            </ParticipationAlert>
          )}

          <section className="card glass detail-block team-panel">
            <div className="team-header">
              <h2>{myTeam.name}</h2>
              {isCaptain && <span className="captain-badge">Капитан</span>}
            </div>

            <div className="team-case-summary">
              <div>
                <span className="meta-label">Трек</span>
                <strong>{selectedTrack?.title ?? 'Не выбран'}</strong>
              </div>
              <div>
                <span className="meta-label">Кейс</span>
                <strong>{selectedCase?.title ?? 'Не выбран'}</strong>
              </div>
            </div>

            <ul className="team-members">
              {myTeam.members.map((m) => (
                <li key={m.user_id}>
                  <span className="member-name">{m.full_name}</span>
                  <span className="member-role">{teamRoleLabel(m.team_role)}</span>
                  {m.user_id === myTeam.captain_id && (
                    <span className="member-captain">капитан</span>
                  )}
                </li>
              ))}
            </ul>

            {isCaptain && hackathon.tracks.length > 0 && (
              <form className="inline-form team-track-form" onSubmit={handleUpdateTrackCase}>
                <h3>Трек и кейс {needsCase && <span className="required-mark">*</span>}</h3>
                <p className="form-section-hint">
                  Выбор кейса открывает возможность сдать решение на следующем шаге.
                </p>
                <label>
                  Трек
                  <select
                    value={trackId}
                    onChange={(e) => {
                      setTrackId(e.target.value);
                      setCaseId('');
                    }}
                  >
                    <option value="">Не выбран</option>
                    {hackathon.tracks.map((t) => (
                      <option key={t.id} value={t.id}>
                        {t.title}
                      </option>
                    ))}
                  </select>
                </label>
                {casesForTrack.length > 0 && (
                  <label>
                    Кейс
                    <select value={caseId} onChange={(e) => setCaseId(e.target.value)}>
                      <option value="">Не выбран</option>
                      {casesForTrack.map((c) => (
                        <option key={c.id} value={c.id}>
                          {c.title}
                        </option>
                      ))}
                    </select>
                  </label>
                )}
                <button type="submit" className="btn-secondary" disabled={busy}>
                  Сохранить выбор
                </button>
              </form>
            )}

            {!isCaptain && hackathon.tracks.length > 0 && !myTeam.case_id && (
              <p className="form-hint team-hint">
                Только капитан может выбрать трек и кейс для команды.
              </p>
            )}

            <div className="team-actions">
              <Link
                to={`/hackathons/${id}/submission`}
                className={`btn-primary ${needsCase ? 'btn-disabled' : ''}`}
                title={needsCase ? 'Сначала выберите кейс' : undefined}
              >
                К сдаче решения
              </Link>
              <button type="button" className="btn-ghost" disabled={busy} onClick={handleLeave}>
                Выйти из команды
              </button>
            </div>
          </section>
        </Reveal>
      ) : (
        <>
          <Reveal delay={80}>
            <section className="card glass form-card wide-form">
              <h2>Создать команду</h2>
              <form className="inline-form" onSubmit={handleCreate}>
                <label>
                  Название
                  <input
                    required
                    minLength={2}
                    maxLength={64}
                    value={teamName}
                    onChange={(e) => setTeamName(e.target.value)}
                    placeholder="Например, QuantumBits"
                  />
                </label>
                {hackathon.tracks.length > 0 && (
                  <>
                    <label>
                      Трек (опционально)
                      <select
                        value={trackId}
                        onChange={(e) => {
                          setTrackId(e.target.value);
                          setCaseId('');
                        }}
                      >
                        <option value="">Позже</option>
                        {hackathon.tracks.map((t) => (
                          <option key={t.id} value={t.id}>
                            {t.title}
                          </option>
                        ))}
                      </select>
                    </label>
                    {casesForTrack.length > 0 && (
                      <label>
                        Кейс
                        <select value={caseId} onChange={(e) => setCaseId(e.target.value)}>
                          <option value="">Позже</option>
                          {casesForTrack.map((c) => (
                            <option key={c.id} value={c.id}>
                              {c.title}
                            </option>
                          ))}
                        </select>
                      </label>
                    )}
                  </>
                )}
                <button type="submit" className="btn-primary" disabled={busy || !status?.can_create_team}>
                  Создать
                </button>
              </form>
            </section>
          </Reveal>

          <Reveal delay={160}>
            <section className="card glass detail-block">
              <h2>Присоединиться</h2>
              <label className="join-role-label">
                Ваша роль в команде
                <select value={joinRole} onChange={(e) => setJoinRole(e.target.value as TeamMemberRole)}>
                  <option value="developer">Разработчик</option>
                  <option value="designer">Дизайнер</option>
                  <option value="data_scientist">Data Scientist</option>
                  <option value="devops_qa">DevOps / QA</option>
                  <option value="other">Другое</option>
                </select>
              </label>
              {teams.length === 0 ? (
                <p className="muted">Пока нет открытых команд — создайте свою.</p>
              ) : (
                <ul className="team-join-list">
                  {teams.map((t) => {
                    const full = t.members.length >= hackathon.max_team_size;
                    const alreadyMember = t.members.some((m) => m.user_id === user?.id);
                    return (
                      <li key={t.id} className="team-join-item">
                        <div>
                          <strong>{t.name}</strong>
                          <span className="team-size">
                            {t.members.length} / {hackathon.max_team_size}
                          </span>
                        </div>
                        <button
                          type="button"
                          className="btn-secondary"
                          disabled={busy || full || alreadyMember}
                          onClick={() => handleJoin(t.id)}
                        >
                          {alreadyMember ? 'Вы в команде' : full ? 'Мест нет' : 'Вступить'}
                        </button>
                      </li>
                    );
                  })}
                </ul>
              )}
            </section>
          </Reveal>
        </>
      )}

      <Reveal delay={240}>
        <div className="participate-footer">
          <Link to={`/hackathons/${id}/participate`} className="btn-ghost">
            К участию
          </Link>
        </div>
      </Reveal>
    </div>
  );
}
