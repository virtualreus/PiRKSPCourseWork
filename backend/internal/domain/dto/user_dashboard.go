package dto

type UserDashboard struct {
	User           User                `json:"user"`
	Stats          UserDashboardStats  `json:"stats"`
	Participations []UserParticipation `json:"participations"`
	Organized      []HackathonListItem `json:"organized_hackathons,omitempty"`
}

type UserDashboardStats struct {
	RegistrationsCount int `json:"registrations_count"`
	ActiveHackathons   int `json:"active_hackathons"`
	TeamsCount         int `json:"teams_count"`
	SubmissionsCount   int `json:"submissions_count"`
	OrganizedCount     int `json:"organized_count,omitempty"`
}

type UserParticipation struct {
	Hackathon         HackathonListItem           `json:"hackathon"`
	RegisteredAt      string                      `json:"registered_at"`
	Team              *UserParticipationTeam      `json:"team,omitempty"`
	Submission        *UserParticipationSubmission `json:"submission,omitempty"`
	CanSubmit         bool                        `json:"can_submit"`
	SubmitBlockReason string                      `json:"submit_block_reason,omitempty"`
}

type UserParticipationTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsCaptain   bool   `json:"is_captain"`
	TeamRole    string `json:"team_role"`
	MemberCount int    `json:"member_count"`
	TrackTitle  string `json:"track_title,omitempty"`
	CaseTitle   string `json:"case_title,omitempty"`
}

type UserParticipationSubmission struct {
	ID          string  `json:"id"`
	Title       string  `json:"title,omitempty"`
	RepoURL     string  `json:"repo_url"`
	SubmittedAt *string `json:"submitted_at,omitempty"`
}
