import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import * as hackathonsApi from '../../api/hackathons';
import type { HackathonListItem } from '../../api/hackathonTypes';
import { ApiError } from '../../api/client';
import { formatDate, statusLabel } from '../../utils/hackathon';

export function OrganizerHackathonsPage() {
  const [items, setItems] = useState<HackathonListItem[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const resp = await hackathonsApi.listOrganizerHackathons();
        setItems(resp.items ?? []);
      } catch (err) {
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError('Не удалось загрузить список');
        }
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <section className="organizer-page">
      <div className="page-toolbar">
        <h1>Мои хакатоны</h1>
        <Link to="/organizer/hackathons/new" className="btn-primary">
          Создать
        </Link>
      </div>

      {loading && <p className="page-loading">Загрузка…</p>}
      {error && <p className="form-error">{error}</p>}

      <div className="organizer-list">
        {items.map((item) => (
          <Link key={item.id} to={`/organizer/hackathons/${item.id}`} className="card glass organizer-row">
            <div>
              <span className={`status-badge status-${item.status}`}>{statusLabel(item.status)}</span>
              <h3>{item.title}</h3>
              <p className="muted-text">Дедлайн: {formatDate(item.submission_deadline_at)}</p>
            </div>
          </Link>
        ))}
      </div>

    </section>
  );
}
