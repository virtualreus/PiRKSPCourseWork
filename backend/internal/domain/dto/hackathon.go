package dto

type HackathonTimeline struct {
	RegistrationOpensAt   string `json:"registration_opens_at"`
	RegistrationClosesAt  string `json:"registration_closes_at"`
	EventStartsAt         string `json:"event_starts_at"`
	EventEndsAt           string `json:"event_ends_at"`
	SubmissionDeadlineAt  string `json:"submission_deadline_at"`
}

type HackathonListItem struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	ShortDescription     string `json:"short_description,omitempty"`
	Status               string `json:"status"`
	Format               string `json:"format,omitempty"`
	RegistrationOpensAt  string `json:"registration_opens_at"`
	SubmissionDeadlineAt string `json:"submission_deadline_at"`
}

type Case struct {
	ID           string `json:"id"`
	TrackID      string `json:"track_id"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	ResourcesURL string `json:"resources_url,omitempty"`
}

type Track struct {
	ID          string `json:"id"`
	HackathonID string `json:"hackathon_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type TrackWithCases struct {
	Track
	Cases []Case `json:"cases"`
}

type HackathonDetail struct {
	HackathonListItem
	Description  string             `json:"description"`
	Timeline     HackathonTimeline  `json:"timeline"`
	MaxTeamSize  int                `json:"max_team_size"`
	PrizesInfo   string             `json:"prizes_info,omitempty"`
	Tracks       []TrackWithCases   `json:"tracks"`
}

type CreateHackathonRequest struct {
	Title            string             `json:"title"`
	Description      string             `json:"description"`
	ShortDescription string             `json:"short_description"`
	Format           string             `json:"format"`
	Timeline         HackathonTimeline  `json:"timeline"`
	MaxTeamSize      int                `json:"max_team_size"`
	PrizesInfo       string             `json:"prizes_info"`
}

type UpdateHackathonRequest struct {
	Title            *string            `json:"title"`
	Description      *string            `json:"description"`
	ShortDescription *string            `json:"short_description"`
	Format           *string            `json:"format"`
	Timeline         *HackathonTimeline `json:"timeline"`
	MaxTeamSize      *int               `json:"max_team_size"`
	PrizesInfo       *string            `json:"prizes_info"`
	Status           *string            `json:"status"`
}

type CreateTrackRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateCaseRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	CustomerName string `json:"customer_name"`
	ResourcesURL string `json:"resources_url"`
}

type HackathonListResponse struct {
	Items []HackathonListItem `json:"items"`
}

type TrackListResponse struct {
	Items []Track `json:"items"`
}

type CaseListResponse struct {
	Items []Case `json:"items"`
}
