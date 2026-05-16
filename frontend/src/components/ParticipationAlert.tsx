import type { ReactNode } from 'react';
import { Link } from 'react-router-dom';

type ParticipationAlertProps = {
  variant: 'info' | 'warning' | 'success' | 'danger';
  title: string;
  children: ReactNode;
  actionLabel?: string;
  actionTo?: string;
};

const ICONS: Record<ParticipationAlertProps['variant'], string> = {
  info: 'ℹ',
  warning: '!',
  success: '✓',
  danger: '×',
};

export function ParticipationAlert({
  variant,
  title,
  children,
  actionLabel,
  actionTo,
}: ParticipationAlertProps) {
  return (
    <aside className={`participation-alert participation-alert-${variant}`} role="status">
      <span className="participation-alert-icon" aria-hidden>
        {ICONS[variant]}
      </span>
      <div className="participation-alert-body">
        <strong>{title}</strong>
        <div className="participation-alert-text">{children}</div>
        {actionLabel && actionTo && (
          <Link to={actionTo} className="btn-secondary btn-sm">
            {actionLabel}
          </Link>
        )}
      </div>
    </aside>
  );
}
