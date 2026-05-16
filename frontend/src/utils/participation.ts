import type { HackathonStatus } from '../api/hackathonTypes';
import type { ParticipationStatus, SubmitBlockReason } from '../api/participationTypes';
import { statusLabel } from './hackathon';

export type SubmitBlockInfo = {
  variant: 'info' | 'warning' | 'success' | 'danger';
  title: string;
  body: string;
  actionLabel?: string;
  actionTo?: string;
};

export function formatDeadlineRemaining(deadlineIso: string): string | null {
  const deadline = new Date(deadlineIso).getTime();
  const diff = deadline - Date.now();
  if (diff <= 0) {
    return null;
  }
  const hours = Math.floor(diff / (1000 * 60 * 60));
  const days = Math.floor(hours / 24);
  if (days > 0) {
    return `осталось ${days} дн. ${hours % 24} ч.`;
  }
  if (hours > 0) {
    return `осталось ${hours} ч.`;
  }
  const minutes = Math.floor(diff / (1000 * 60));
  return `осталось ${minutes} мин.`;
}

export function getSubmitBlockInfo(
  reason: SubmitBlockReason | undefined,
  hackathonId: string,
  isCaptain: boolean,
): SubmitBlockInfo | null {
  switch (reason) {
    case 'no_team':
      return {
        variant: 'warning',
        title: 'Нужна команда',
        body: 'Сначала создайте команду или вступите в существующую — без этого сдача недоступна.',
        actionLabel: 'Перейти к команде',
        actionTo: `/hackathons/${hackathonId}/team`,
      };
    case 'no_case':
      return {
        variant: 'warning',
        title: 'Не выбран кейс',
        body: isCaptain
          ? 'Капитан команды должен выбрать трек и кейс на странице команды. После сохранения откроется форма сдачи.'
          : 'Попросите капитана команды выбрать трек и кейс — без этого сдача решения недоступна.',
        actionLabel: isCaptain ? 'Выбрать кейс' : 'Открыть команду',
        actionTo: `/hackathons/${hackathonId}/team`,
      };
    case 'deadline_passed':
      return {
        variant: 'danger',
        title: 'Дедлайн сдачи прошёл',
        body: 'Редактирование решения больше недоступно. Если вы уже отправляли работу, организатор видит последнюю версию.',
      };
    case 'hackathon_finished':
      return {
        variant: 'danger',
        title: 'Хакатон завершён',
        body: 'Статус мероприятия — «Завершён». Новые сдачи и правки не принимаются.',
      };
    case 'hackathon_not_active':
      return {
        variant: 'info',
        title: 'Сдача пока недоступна',
        body: 'Хакатон ещё не в фазе регистрации или проведения. Следите за статусом на карточке мероприятия.',
        actionLabel: 'К хакатону',
        actionTo: `/hackathons/${hackathonId}`,
      };
    default:
      return null;
  }
}

export function getStep3Hint(status: ParticipationStatus): string {
  if (!status.team) {
    return 'Сначала соберите команду на шаге 2.';
  }
  if (status.submit_block_reason === 'no_case') {
    return 'Выберите кейс в настройках команды — это обязательный шаг перед сдачей.';
  }
  if (status.has_submission) {
    return 'Решение уже отправлено. Можно обновить ссылки до дедлайна.';
  }
  if (status.can_submit) {
    return 'Заполните ссылки на репозиторий, демо и материалы питча.';
  }
  const info = getSubmitBlockInfo(status.submit_block_reason, status.hackathon_id, false);
  return info?.body ?? 'Сдача временно недоступна.';
}

export function participationProgress(status: ParticipationStatus): {
  done: number;
  total: number;
  percent: number;
} {
  let done = 0;
  if (status.is_registered) {
    done += 1;
  }
  if (status.team) {
    done += 1;
  }
  if (status.has_submission) {
    done += 1;
  }
  const total = 3;
  return { done, total, percent: Math.round((done / total) * 100) };
}

export function hackathonStatusChip(status: HackathonStatus): string {
  return statusLabel(status);
}
