export type HackathonStatus = 'draft' | 'registration' | 'running' | 'finished';
export type HackathonFormat = 'online' | 'offline' | 'hybrid';

export interface HackathonTimeline {
  registration_opens_at: string;
  registration_closes_at: string;
  event_starts_at: string;
  event_ends_at: string;
  submission_deadline_at: string;
}

export interface HackathonListItem {
  id: string;
  title: string;
  short_description?: string;
  status: HackathonStatus;
  format?: HackathonFormat;
  registration_opens_at: string;
  submission_deadline_at: string;
}

export interface Case {
  id: string;
  track_id: string;
  title: string;
  description?: string;
  customer_name?: string;
  resources_url?: string;
}

export interface Track {
  id: string;
  hackathon_id: string;
  title: string;
  description?: string;
}

export interface TrackWithCases extends Track {
  cases: Case[];
}

export interface HackathonDetail extends HackathonListItem {
  description: string;
  timeline: HackathonTimeline;
  max_team_size: number;
  prizes_info?: string;
  tracks: TrackWithCases[];
}

export interface CreateHackathonRequest {
  title: string;
  description: string;
  short_description?: string;
  format?: HackathonFormat;
  timeline: HackathonTimeline;
  max_team_size?: number;
  prizes_info?: string;
}

export interface UpdateHackathonRequest {
  title?: string;
  description?: string;
  short_description?: string;
  format?: HackathonFormat;
  timeline?: HackathonTimeline;
  max_team_size?: number;
  prizes_info?: string;
}

export interface CreateTrackRequest {
  title: string;
  description?: string;
}

export interface CreateCaseRequest {
  title: string;
  description: string;
  customer_name?: string;
  resources_url?: string;
}
