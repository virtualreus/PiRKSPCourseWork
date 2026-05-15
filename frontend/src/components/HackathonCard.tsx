import { Link } from 'react-router-dom';

import type { HackathonListItem } from '../api/hackathonTypes';
import { formatDate, statusLabel } from '../utils/hackathon';

export function HackathonCard({ item }: { item: HackathonListItem }) {
  return (
    <Link to={`/hackathons/${item.id}`} className="hackathon-card card glass">
      <div className="hackathon-card-head">
        <span className={`status-badge status-${item.status}`}>{statusLabel(item.status)}</span>
        {item.format && <span className="format-badge">{item.format}</span>}
      </div>
      <h3>{item.title}</h3>
      {item.short_description && <p className="card-desc">{item.short_description}</p>}
      <dl className="hackathon-card-dates">
        <div>
          <dt>Регистрация</dt>
          <dd>{formatDate(item.registration_opens_at)}</dd>
        </div>
        <div>
          <dt>Дедлайн сдачи</dt>
          <dd>{formatDate(item.submission_deadline_at)}</dd>
        </div>
      </dl>
    </Link>
  );
}
