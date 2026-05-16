import { useState, type FormEvent } from 'react';

import * as hackathonsApi from '../../api/hackathons';
import type { Case, CreateCaseRequest, CreateTrackRequest, TrackWithCases } from '../../api/hackathonTypes';
import { ApiError } from '../../api/client';

type TrackCaseManagerProps = {
  hackathonId: string;
  tracks: TrackWithCases[];
  isDraft: boolean;
  busy: boolean;
  setBusy: (v: boolean) => void;
  onChanged: () => Promise<void>;
  onError: (msg: string) => void;
};

export function TrackCaseManager({
  hackathonId,
  tracks,
  isDraft,
  busy,
  setBusy,
  onChanged,
  onError,
}: TrackCaseManagerProps) {
  const [trackTitle, setTrackTitle] = useState('');
  const [trackDescription, setTrackDescription] = useState('');
  const [caseTitle, setCaseTitle] = useState('');
  const [caseDescription, setCaseDescription] = useState('');
  const [caseCustomer, setCaseCustomer] = useState('');
  const [caseResources, setCaseResources] = useState('');
  const [selectedTrackId, setSelectedTrackId] = useState(tracks[0]?.id ?? '');
  const [editingTrackId, setEditingTrackId] = useState<string | null>(null);
  const [editingCaseId, setEditingCaseId] = useState<string | null>(null);
  const [editTrackTitle, setEditTrackTitle] = useState('');
  const [editTrackDesc, setEditTrackDesc] = useState('');
  const [editCaseForm, setEditCaseForm] = useState<CreateCaseRequest>({
    title: '',
    description: '',
  });

  async function runAction(action: () => Promise<void>) {
    setBusy(true);
    onError('');
    try {
      await action();
      await onChanged();
    } catch (err) {
      if (err instanceof ApiError) {
        onError(err.message);
      } else {
        onError('Не удалось выполнить действие');
      }
    } finally {
      setBusy(false);
    }
  }

  async function handleAddTrack(e: FormEvent) {
    e.preventDefault();
    await runAction(async () => {
      await hackathonsApi.createTrack(hackathonId, {
        title: trackTitle.trim(),
        description: trackDescription.trim() || undefined,
      });
      setTrackTitle('');
      setTrackDescription('');
    });
  }

  async function handleAddCase(e: FormEvent) {
    e.preventDefault();
    if (!selectedTrackId) {
      return;
    }
    await runAction(async () => {
      await hackathonsApi.createCase(selectedTrackId, {
        title: caseTitle.trim(),
        description: caseDescription.trim(),
        customer_name: caseCustomer.trim() || undefined,
        resources_url: caseResources.trim() || undefined,
      });
      setCaseTitle('');
      setCaseDescription('');
      setCaseCustomer('');
      setCaseResources('');
    });
  }

  function startEditTrack(track: TrackWithCases) {
    setEditingTrackId(track.id);
    setEditingCaseId(null);
    setEditTrackTitle(track.title);
    setEditTrackDesc(track.description ?? '');
  }

  function startEditCase(c: Case) {
    setEditingCaseId(c.id);
    setEditingTrackId(null);
    setEditCaseForm({
      title: c.title,
      description: c.description ?? '',
      customer_name: c.customer_name,
      resources_url: c.resources_url,
    });
  }

  async function saveTrack(trackId: string) {
    const body: CreateTrackRequest = {
      title: editTrackTitle.trim(),
      description: editTrackDesc.trim() || undefined,
    };
    await runAction(async () => {
      await hackathonsApi.updateTrack(trackId, body);
      setEditingTrackId(null);
    });
  }

  async function saveCase(caseId: string) {
    await runAction(async () => {
      await hackathonsApi.updateCase(caseId, {
        title: editCaseForm.title.trim(),
        description: editCaseForm.description.trim(),
        customer_name: editCaseForm.customer_name?.trim() || undefined,
        resources_url: editCaseForm.resources_url?.trim() || undefined,
      });
      setEditingCaseId(null);
    });
  }

  async function removeTrack(trackId: string, title: string) {
    if (!confirm(`Удалить трек «${title}» и все его кейсы?`)) {
      return;
    }
    await runAction(async () => {
      await hackathonsApi.deleteTrack(trackId);
      if (editingTrackId === trackId) {
        setEditingTrackId(null);
      }
    });
  }

  async function removeCase(caseId: string, title: string) {
    if (!confirm(`Удалить кейс «${title}»?`)) {
      return;
    }
    await runAction(async () => {
      await hackathonsApi.deleteCase(caseId);
      if (editingCaseId === caseId) {
        setEditingCaseId(null);
      }
    });
  }

  return (
    <>
      <div className="organizer-tracks-list">
        {tracks.length === 0 ? (
          <p className="muted organizer-empty-hint">
            Добавьте хотя бы один трек и кейс — без них участники не смогут сдать решение.
          </p>
        ) : (
          tracks.map((track) => (
            <article key={track.id} className="organizer-track-card">
              {editingTrackId === track.id && isDraft ? (
                <div className="inline-form">
                  <label>
                    Название трека
                    <input
                      required
                      value={editTrackTitle}
                      onChange={(e) => setEditTrackTitle(e.target.value)}
                    />
                  </label>
                  <label>
                    Описание
                    <input
                      value={editTrackDesc}
                      onChange={(e) => setEditTrackDesc(e.target.value)}
                    />
                  </label>
                  <div className="organizer-inline-actions">
                    <button
                      type="button"
                      className="btn-primary btn-sm"
                      disabled={busy}
                      onClick={() => saveTrack(track.id)}
                    >
                      Сохранить
                    </button>
                    <button
                      type="button"
                      className="btn-ghost btn-sm"
                      disabled={busy}
                      onClick={() => setEditingTrackId(null)}
                    >
                      Отмена
                    </button>
                  </div>
                </div>
              ) : (
                <div className="organizer-track-head">
                  <div>
                    <h3>{track.title}</h3>
                    {track.description && <p className="muted">{track.description}</p>}
                  </div>
                  {isDraft && (
                    <div className="organizer-item-actions">
                      <button
                        type="button"
                        className="btn-ghost btn-sm"
                        disabled={busy}
                        onClick={() => startEditTrack(track)}
                      >
                        Изменить
                      </button>
                      <button
                        type="button"
                        className="btn-ghost btn-sm btn-danger-text"
                        disabled={busy}
                        onClick={() => removeTrack(track.id, track.title)}
                      >
                        Удалить
                      </button>
                    </div>
                  )}
                </div>
              )}

              <ul className="organizer-cases-list">
                {(track.cases ?? []).map((c) => (
                  <li key={c.id}>
                    {editingCaseId === c.id && isDraft ? (
                      <div className="inline-form organizer-case-edit">
                        <label>
                          Название
                          <input
                            required
                            value={editCaseForm.title}
                            onChange={(e) =>
                              setEditCaseForm({ ...editCaseForm, title: e.target.value })
                            }
                          />
                        </label>
                        <label>
                          Описание
                          <textarea
                            required
                            rows={2}
                            value={editCaseForm.description}
                            onChange={(e) =>
                              setEditCaseForm({ ...editCaseForm, description: e.target.value })
                            }
                          />
                        </label>
                        <label>
                          Заказчик
                          <input
                            value={editCaseForm.customer_name ?? ''}
                            onChange={(e) =>
                              setEditCaseForm({ ...editCaseForm, customer_name: e.target.value })
                            }
                          />
                        </label>
                        <div className="organizer-inline-actions">
                          <button
                            type="button"
                            className="btn-primary btn-sm"
                            disabled={busy}
                            onClick={() => saveCase(c.id)}
                          >
                            Сохранить
                          </button>
                          <button
                            type="button"
                            className="btn-ghost btn-sm"
                            disabled={busy}
                            onClick={() => setEditingCaseId(null)}
                          >
                            Отмена
                          </button>
                        </div>
                      </div>
                    ) : (
                      <div className="organizer-case-row">
                        <div>
                          <strong>{c.title}</strong>
                          {c.customer_name && (
                            <span className="muted"> · {c.customer_name}</span>
                          )}
                          {c.description && <p className="muted case-desc">{c.description}</p>}
                        </div>
                        {isDraft && (
                          <div className="organizer-item-actions">
                            <button
                              type="button"
                              className="btn-ghost btn-sm"
                              disabled={busy}
                              onClick={() => startEditCase(c)}
                            >
                              Изменить
                            </button>
                            <button
                              type="button"
                              className="btn-ghost btn-sm btn-danger-text"
                              disabled={busy}
                              onClick={() => removeCase(c.id, c.title)}
                            >
                              Удалить
                            </button>
                          </div>
                        )}
                      </div>
                    )}
                  </li>
                ))}
              </ul>
            </article>
          ))
        )}
      </div>

      {isDraft && (
        <>
          <form className="inline-form organizer-add-form" onSubmit={handleAddTrack}>
            <h3>Добавить трек</h3>
            <label>
              Название
              <input
                required
                value={trackTitle}
                onChange={(e) => setTrackTitle(e.target.value)}
                placeholder="Например, FinTech"
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

          {tracks.length > 0 && (
            <form className="inline-form organizer-add-form" onSubmit={handleAddCase}>
              <h3>Добавить кейс</h3>
              <label>
                Трек
                <select
                  value={selectedTrackId || tracks[0]?.id}
                  onChange={(e) => setSelectedTrackId(e.target.value)}
                >
                  {tracks.map((t) => (
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
                <input value={caseCustomer} onChange={(e) => setCaseCustomer(e.target.value)} />
              </label>
              <label>
                Ссылка на ресурсы
                <input
                  type="url"
                  value={caseResources}
                  onChange={(e) => setCaseResources(e.target.value)}
                />
              </label>
              <button type="submit" className="btn-secondary" disabled={busy}>
                Добавить кейс
              </button>
            </form>
          )}
        </>
      )}
    </>
  );
}
