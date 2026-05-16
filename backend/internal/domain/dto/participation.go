package dto

type HackathonRegistration struct {
	ID           string `json:"id"`
	HackathonID  string `json:"hackathon_id"`
	UserID       string `json:"user_id"`
	RegisteredAt string `json:"registered_at"`
}

type HackathonRegistrationWithUser struct {
	HackathonRegistration
	User User `json:"user"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	TeamRole string `json:"team_role"`
}

type Team struct {
	ID          string       `json:"id"`
	HackathonID string       `json:"hackathon_id"`
	Name        string       `json:"name"`
	CaptainID   string       `json:"captain_id"`
	TrackID     *string      `json:"track_id,omitempty"`
	CaseID      *string      `json:"case_id,omitempty"`
	Members     []TeamMember `json:"members"`
}

type ParticipationStatus struct {
	HackathonID          string                 `json:"hackathon_id"`
	IsRegistered         bool                   `json:"is_registered"`
	Registration         *HackathonRegistration `json:"registration,omitempty"`
	Team                 *Team                  `json:"team,omitempty"`
	HasSubmission        bool                   `json:"has_submission"`
	CanRegister          bool                   `json:"can_register"`
	CanCreateTeam        bool                   `json:"can_create_team"`
	CanSubmit            bool                   `json:"can_submit"`
	SubmitBlockReason    string                 `json:"submit_block_reason,omitempty"`
	HackathonStatus      string                 `json:"hackathon_status"`
	SubmissionDeadlineAt string                 `json:"submission_deadline_at"`
}

type CreateTeamRequest struct {
	Name     string  `json:"name"`
	TrackID  *string `json:"track_id"`
	CaseID   *string `json:"case_id"`
	TeamRole string  `json:"team_role"`
}

type UpdateTeamRequest struct {
	Name    *string `json:"name"`
	TrackID *string `json:"track_id"`
	CaseID  *string `json:"case_id"`
}

type JoinTeamRequest struct {
	TeamRole string `json:"team_role"`
}

type UpdateTeamMemberRoleRequest struct {
	TeamRole string `json:"team_role"`
}

type Submission struct {
	ID          string  `json:"id"`
	TeamID      string  `json:"team_id"`
	HackathonID string  `json:"hackathon_id"`
	Title       string  `json:"title,omitempty"`
	Summary     string  `json:"summary,omitempty"`
	RepoURL     string  `json:"repo_url"`
	DemoURL     string  `json:"demo_url,omitempty"`
	PitchURL    string  `json:"pitch_url,omitempty"`
	VideoURL    string  `json:"video_url,omitempty"`
	SubmittedAt *string `json:"submitted_at,omitempty"`
	UpdatedAt   string  `json:"updated_at"`
}

type SubmissionWithTeam struct {
	Submission
	TeamName   string `json:"team_name"`
	CaseTitle  string `json:"case_title,omitempty"`
	TrackTitle string `json:"track_title,omitempty"`
}

type UpsertSubmissionRequest struct {
	Title    string `json:"title"`
	Summary  string `json:"summary"`
	RepoURL  string `json:"repo_url"`
	DemoURL  string `json:"demo_url"`
	PitchURL string `json:"pitch_url"`
	VideoURL string `json:"video_url"`
}

type TeamListResponse struct {
	Items []Team `json:"items"`
}

type RegistrationListResponse struct {
	Items []HackathonRegistrationWithUser `json:"items"`
}

type SubmissionListResponse struct {
	Items []SubmissionWithTeam `json:"items"`
}
