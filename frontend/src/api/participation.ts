import { api } from './client';
import type {
  CreateTeamRequest,
  HackathonRegistration,
  HackathonRegistrationWithUser,
  JoinTeamRequest,
  ParticipationStatus,
  Submission,
  SubmissionWithTeam,
  Team,
  UpdateTeamMemberRoleRequest,
  UpdateTeamRequest,
  UpsertSubmissionRequest,
} from './participationTypes';

export function getParticipation(hackathonId: string): Promise<ParticipationStatus> {
  return api<ParticipationStatus>(`/hackathons/${hackathonId}/participation`);
}

export function registerForHackathon(hackathonId: string): Promise<HackathonRegistration> {
  return api<HackathonRegistration>(`/hackathons/${hackathonId}/register`, { method: 'POST' });
}

export function unregisterFromHackathon(hackathonId: string): Promise<void> {
  return api<void>(`/hackathons/${hackathonId}/register`, { method: 'DELETE' });
}

export function listTeams(hackathonId: string): Promise<{ items: Team[] }> {
  return api<{ items: Team[] }>(`/hackathons/${hackathonId}/teams`);
}

export function createTeam(hackathonId: string, body: CreateTeamRequest): Promise<Team> {
  return api<Team>(`/hackathons/${hackathonId}/teams`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function getTeam(teamId: string): Promise<Team> {
  return api<Team>(`/teams/${teamId}`);
}

export function updateTeam(teamId: string, body: UpdateTeamRequest): Promise<Team> {
  return api<Team>(`/teams/${teamId}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}

export function joinTeam(teamId: string, body?: JoinTeamRequest): Promise<Team> {
  return api<Team>(`/teams/${teamId}/join`, {
    method: 'POST',
    body: JSON.stringify(body ?? {}),
  });
}

export function leaveTeam(teamId: string): Promise<void> {
  return api<void>(`/teams/${teamId}/leave`, { method: 'POST' });
}

export function updateMemberRole(
  teamId: string,
  userId: string,
  body: UpdateTeamMemberRoleRequest,
): Promise<void> {
  return api<void>(`/teams/${teamId}/members/${userId}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}

export function getTeamSubmission(teamId: string): Promise<Submission> {
  return api<Submission>(`/teams/${teamId}/submission`);
}

export function upsertTeamSubmission(
  teamId: string,
  body: UpsertSubmissionRequest,
): Promise<Submission> {
  return api<Submission>(`/teams/${teamId}/submission`, {
    method: 'PUT',
    body: JSON.stringify(body),
  });
}

export function listOrganizerRegistrations(
  hackathonId: string,
): Promise<{ items: HackathonRegistrationWithUser[] }> {
  return api<{ items: HackathonRegistrationWithUser[] }>(
    `/organizer/hackathons/${hackathonId}/registrations`,
  );
}

export function listOrganizerSubmissions(
  hackathonId: string,
): Promise<{ items: SubmissionWithTeam[] }> {
  return api<{ items: SubmissionWithTeam[] }>(
    `/organizer/hackathons/${hackathonId}/submissions`,
  );
}
