import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

import * as hackathonsApi from '../../api/hackathons';
import type { HackathonListItem } from '../../api/hackathonTypes';
import { ApiError } from '../../api/client';
import { Reveal } from '../../components/Reveal';
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
    <section className="organizer-page organizer-page-wide">
      <Reveal>
        <div className="page-toolbar">
          <div>
            <h1>Мои хакатоны</h1>
            <p className="page-toolbar-lead">Черновики, публикация и мониторинг участников</p>
          </div>
          <Link to="/organizer/hackathons/new" className="btn-primary">
            Создать
          </Link>
        </div>
      </Reveal>

      {loading && (
        <div className="page-loading" aria-busy>
          <div className="spinner" />
        </div>
      )}
      {error && <p className="form-error participate-banner">{error}</p>}

      {!loading && items.length === 0 && (
        <Reveal delay={80}>
          <div className="card glass profile-empty organizer-empty-block">
            <h2>Пока нет хакатонов</h2>
            <p className="muted">
              Создайте черновик, добавьте треки и кейсы, затем опубликуйте событие в каталог.
            </p>
            <Link to="/organizer/hackathons/new" className="btn-primary">
              Создать первый хакатон
            </Link>
          </div>
        </Reveal>
      )}

      {!loading && items.length > 0 && (
        <div className="organizer-list">
          {items.map((item, i) => (
            <Reveal key={item.id} delay={i * 50}>
              <Link
                to={`/organizer/hackathons/${item.id}`}
                className="card glass organizer-row"
              >
                <div>
                  <span className={`status-badge status-${item.status}`}>
                    {statusLabel(item.status)}
                  </span>
                  <h3>{item.title}</h3>
                  {item.short_description && (
                    <p className="organizer-row-desc">{item.short_description}</p>
                  )}
                  <p className="muted-text">Дедлайн: {formatDate(item.submission_deadline_at)}</p>
                </div>
                <span className="organizer-row-chevron" aria-hidden>
                  →
                </span>
              </Link>
            </Reveal>
          ))}
        </div>
      )}
    </section>
  );
}
