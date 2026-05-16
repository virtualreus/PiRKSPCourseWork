import { api } from './client';
import type {
  Case,
  CreateCaseRequest,
  CreateHackathonRequest,
  CreateTrackRequest,
  HackathonDetail,
  HackathonListItem,
  Track,
  UpdateHackathonRequest,
} from './hackathonTypes';

export function listHackathons(): Promise<{ items: HackathonListItem[] }> {
  return api<{ items: HackathonListItem[] }>('/hackathons');
}

export function getHackathon(id: string): Promise<HackathonDetail> {
  return api<HackathonDetail>(`/hackathons/${id}`);
}

export function listOrganizerHackathons(): Promise<{ items: HackathonListItem[] }> {
  return api<{ items: HackathonListItem[] }>('/organizer/hackathons');
}

export function getOrganizerHackathon(id: string): Promise<HackathonDetail> {
  return api<HackathonDetail>(`/organizer/hackathons/${id}`);
}

export function createHackathon(body: CreateHackathonRequest): Promise<HackathonDetail> {
  return api<HackathonDetail>('/organizer/hackathons', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function updateHackathon(id: string, body: UpdateHackathonRequest): Promise<HackathonDetail> {
  return api<HackathonDetail>(`/organizer/hackathons/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}

export function deleteHackathon(id: string): Promise<void> {
  return api<void>(`/organizer/hackathons/${id}`, { method: 'DELETE' });
}

export function publishHackathon(id: string): Promise<HackathonDetail> {
  return api<HackathonDetail>(`/organizer/hackathons/${id}/publish`, { method: 'POST' });
}

export function createTrack(hackathonId: string, body: CreateTrackRequest): Promise<Track> {
  return api<Track>(`/organizer/hackathons/${hackathonId}/tracks`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function createCase(trackId: string, body: CreateCaseRequest): Promise<Case> {
  return api<Case>(`/organizer/tracks/${trackId}/cases`, {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export function updateTrack(trackId: string, body: CreateTrackRequest): Promise<Track> {
  return api<Track>(`/organizer/tracks/${trackId}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}

export function deleteTrack(trackId: string): Promise<void> {
  return api<void>(`/organizer/tracks/${trackId}`, { method: 'DELETE' });
}

export function updateCase(caseId: string, body: CreateCaseRequest): Promise<Case> {
  return api<Case>(`/organizer/cases/${caseId}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
  });
}

export function deleteCase(caseId: string): Promise<void> {
  return api<void>(`/organizer/cases/${caseId}`, { method: 'DELETE' });
}
