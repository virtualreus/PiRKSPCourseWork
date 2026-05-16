import type { HackathonListItem } from './hackathonTypes';
import type { SubmitBlockReason } from './participationTypes';
import type { User } from './types';

export interface UserDashboardStats {
  registrations_count: number;
  active_hackathons: number;
  teams_count: number;
  submissions_count: number;
  organized_count?: number;
}

export interface UserParticipationTeam {
  id: string;
  name: string;
  is_captain: boolean;
  team_role: string;
  member_count: number;
  track_title?: string;
  case_title?: string;
}

export interface UserParticipationSubmission {
  id: string;
  title?: string;
  repo_url: string;
  submitted_at?: string | null;
}

export interface UserParticipation {
  hackathon: HackathonListItem;
  registered_at: string;
  team?: UserParticipationTeam;
  submission?: UserParticipationSubmission;
  can_submit: boolean;
  submit_block_reason?: SubmitBlockReason;
}

export interface UserDashboard {
  user: User;
  stats: UserDashboardStats;
  participations: UserParticipation[];
  organized_hackathons?: HackathonListItem[];
}
