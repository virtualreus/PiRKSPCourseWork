import type { HackathonStatus } from './hackathonTypes';
import type { User } from './types';

export type SubmitBlockReason =
  | 'no_team'
  | 'no_case'
  | 'deadline_passed'
  | 'hackathon_finished'
  | 'hackathon_not_active';

export type TeamMemberRole =
  | 'team_lead'
  | 'developer'
  | 'designer'
  | 'data_scientist'
  | 'devops_qa'
  | 'other';

export interface HackathonRegistration {
  id: string;
  hackathon_id: string;
  user_id: string;
  registered_at: string;
}

export interface TeamMember {
  user_id: string;
  full_name: string;
  team_role: TeamMemberRole;
}

export interface Team {
  id: string;
  hackathon_id: string;
  name: string;
  captain_id: string;
  track_id?: string | null;
  case_id?: string | null;
  members: TeamMember[];
}

export interface ParticipationStatus {
  hackathon_id: string;
  is_registered: boolean;
  registration?: HackathonRegistration;
  team?: Team;
  has_submission: boolean;
  can_register: boolean;
  can_create_team: boolean;
  can_submit: boolean;
  submit_block_reason?: SubmitBlockReason;
  hackathon_status: HackathonStatus;
  submission_deadline_at?: string;
}

export interface Submission {
  id: string;
  team_id: string;
  hackathon_id: string;
  title?: string;
  summary?: string;
  repo_url: string;
  demo_url?: string;
  pitch_url?: string;
  video_url?: string;
  submitted_at?: string | null;
  updated_at: string;
}

export interface SubmissionWithTeam extends Submission {
  team_name: string;
  case_title?: string;
  track_title?: string;
}

export interface HackathonRegistrationWithUser extends HackathonRegistration {
  user: User;
}

export interface CreateTeamRequest {
  name: string;
  track_id?: string;
  case_id?: string;
  team_role?: TeamMemberRole;
}

export interface UpdateTeamRequest {
  name?: string;
  track_id?: string;
  case_id?: string;
}

export interface JoinTeamRequest {
  team_role?: TeamMemberRole;
}

export interface UpsertSubmissionRequest {
  title?: string;
  summary?: string;
  repo_url: string;
  demo_url?: string;
  pitch_url?: string;
  video_url?: string;
}

export interface UpdateTeamMemberRoleRequest {
  team_role: TeamMemberRole;
}
