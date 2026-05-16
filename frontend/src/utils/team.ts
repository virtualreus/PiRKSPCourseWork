import type { TeamMemberRole } from '../api/participationTypes';

const ROLE_LABELS: Record<TeamMemberRole, string> = {
  team_lead: 'Тимлид',
  developer: 'Разработчик',
  designer: 'Дизайнер',
  data_scientist: 'Data Scientist',
  devops_qa: 'DevOps / QA',
  other: 'Другое',
};

export function teamRoleLabel(role: TeamMemberRole): string {
  return ROLE_LABELS[role] ?? role;
}
