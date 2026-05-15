import type { HackathonStatus } from '../api/hackathonTypes';

export function statusLabel(status: HackathonStatus): string {
  switch (status) {
    case 'draft':
      return 'Черновик';
    case 'registration':
      return 'Регистрация';
    case 'running':
      return 'Идёт хакатон';
    case 'finished':
      return 'Завершён';
    default:
      return status;
  }
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function toRFC3339(localValue: string): string {
  return new Date(localValue).toISOString();
}

export function toDatetimeLocal(iso: string): string {
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function defaultTimeline() {
  const now = new Date();
  now.setMinutes(0, 0, 0);
  const addDays = (days: number) => {
    const d = new Date(now);
    d.setDate(d.getDate() + days);
    return d.toISOString();
  };
  return {
    registration_opens_at: now.toISOString(),
    registration_closes_at: addDays(14),
    event_starts_at: addDays(15),
    event_ends_at: addDays(17),
    submission_deadline_at: addDays(17),
  };
}
