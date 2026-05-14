import { useEffect, useState, type FormEvent } from 'react';

import { ApiError } from '../api/client';
import * as authApi from '../api/auth';
import { useAuth } from '../context/AuthContext';

export function ProfilePage() {
  const { user, refreshUser } = useAuth();
  const [fullName, setFullName] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (user) {
      setFullName(user.full_name);
    }
  }, [user]);

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

  return (
    <div className="profile-page">
      <div className="card form-card">
        <h1>Профиль</h1>

        <dl className="profile-meta">
          <div>
            <dt>Email</dt>
            <dd>{user.email}</dd>
          </div>
          <div>
            <dt>Роль</dt>
            <dd>{user.platform_role === 'organizer' ? 'Организатор' : 'Участник'}</dd>
          </div>
          <div>
            <dt>Регистрация</dt>
            <dd>{new Date(user.created_at).toLocaleString('ru-RU')}</dd>
          </div>
        </dl>

        <form onSubmit={handleSubmit}>
          {message && <p className="form-success">{message}</p>}
          {error && <p className="form-error">{error}</p>}

          <label>
            ФИО
            <input
              type="text"
              required
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
            />
          </label>

          <button type="submit" className="btn-primary" disabled={saving}>
            {saving ? 'Сохранение…' : 'Сохранить'}
          </button>
        </form>
      </div>
    </div>
  );
}
