import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { ApiError } from '../api/client';
import { useAuth } from '../context/AuthContext';
import type { PlatformRole } from '../api/types';

export function RegisterPage() {
  const { register } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [fullName, setFullName] = useState('');
  const [platformRole, setPlatformRole] = useState<PlatformRole>('participant');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await register({
        email,
        password,
        full_name: fullName,
        platform_role: platformRole,
      });
      navigate('/', { replace: true });
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Не удалось зарегистрироваться');
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="auth-page">
      <form className="card glass form-card" onSubmit={handleSubmit}>
        <h1>Регистрация</h1>
        <p className="form-hint">Участник или организатор хакатонов</p>

        {error && <p className="form-error">{error}</p>}

        <label>
          ФИО
          <input
            type="text"
            autoComplete="name"
            required
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
          />
        </label>

        <label>
          Email
          <input
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>

        <label>
          Пароль
          <input
            type="password"
            autoComplete="new-password"
            required
            minLength={8}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>

        <label>
          Роль на платформе
          <select
            value={platformRole}
            onChange={(e) => setPlatformRole(e.target.value as PlatformRole)}
          >
            <option value="participant">Участник</option>
            <option value="organizer">Организатор</option>
          </select>
        </label>

        <button type="submit" className="btn-primary btn-block" disabled={submitting}>
          {submitting ? 'Создание…' : 'Создать аккаунт'}
        </button>

        <p className="form-footer">
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </p>
      </form>
    </div>
  );
}
