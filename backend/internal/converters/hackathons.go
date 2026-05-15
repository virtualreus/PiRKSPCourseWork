package converters

import (
	"database/sql"
	"time"

	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
)

type HackathonsConverter struct{}

func NewHackathonsConverter() *HackathonsConverter {
	return &HackathonsConverter{}
}

func (c *HackathonsConverter) ToListItem(h entities.Hackathon) dto.HackathonListItem {
	item := dto.HackathonListItem{
		ID:                   h.ID.String(),
		Title:                h.Title,
		Status:               h.Status,
		Format:               h.Format,
		RegistrationOpensAt:  formatTime(h.RegistrationOpensAt),
		SubmissionDeadlineAt: formatTime(h.SubmissionDeadlineAt),
	}
	if h.ShortDescription.Valid {
		item.ShortDescription = h.ShortDescription.String
	}
	return item
}

func (c *HackathonsConverter) ToDetail(h entities.Hackathon, tracks []entities.Track, cases []entities.Case) dto.HackathonDetail {
	detail := dto.HackathonDetail{
		HackathonListItem: c.ToListItem(h),
		Description:       h.Description,
		Timeline:          c.ToTimeline(h),
		MaxTeamSize:       h.MaxTeamSize,
		Tracks:            c.BuildTracksWithCases(tracks, cases),
	}
	if h.PrizesInfo.Valid {
		detail.PrizesInfo = h.PrizesInfo.String
	}
	return detail
}

func (c *HackathonsConverter) ToTimeline(h entities.Hackathon) dto.HackathonTimeline {
	return dto.HackathonTimeline{
		RegistrationOpensAt:  formatTime(h.RegistrationOpensAt),
		RegistrationClosesAt: formatTime(h.RegistrationClosesAt),
		EventStartsAt:        formatTime(h.EventStartsAt),
		EventEndsAt:          formatTime(h.EventEndsAt),
		SubmissionDeadlineAt: formatTime(h.SubmissionDeadlineAt),
	}
}

func (c *HackathonsConverter) BuildTracksWithCases(tracks []entities.Track, cases []entities.Case) []dto.TrackWithCases {
	byTrack := make(map[string][]dto.Case)
	for _, item := range cases {
		byTrack[item.TrackID.String()] = append(byTrack[item.TrackID.String()], c.ToCase(item))
	}

	out := make([]dto.TrackWithCases, 0, len(tracks))
	for _, t := range tracks {
		out = append(out, dto.TrackWithCases{
			Track: c.ToTrack(t),
			Cases: byTrack[t.ID.String()],
		})
	}

	return out
}

func (c *HackathonsConverter) ToTrack(t entities.Track) dto.Track {
	track := dto.Track{
		ID:          t.ID.String(),
		HackathonID: t.HackathonID.String(),
		Title:       t.Title,
	}
	if t.Description.Valid {
		track.Description = t.Description.String
	}
	return track
}

func (c *HackathonsConverter) ToCase(item entities.Case) dto.Case {
	caze := dto.Case{
		ID:      item.ID.String(),
		TrackID: item.TrackID.String(),
		Title:   item.Title,
	}
	if item.Description != "" {
		caze.Description = item.Description
	}
	if item.CustomerName.Valid {
		caze.CustomerName = item.CustomerName.String
	}
	if item.ResourcesURL.Valid {
		caze.ResourcesURL = item.ResourcesURL.String
	}
	return caze
}

func (c *HackathonsConverter) TimelineToEntity(t dto.HackathonTimeline) (entities.Hackathon, error) {
	opens, err := parseTime(t.RegistrationOpensAt)
	if err != nil {
		return entities.Hackathon{}, err
	}
	closes, err := parseTime(t.RegistrationClosesAt)
	if err != nil {
		return entities.Hackathon{}, err
	}
	starts, err := parseTime(t.EventStartsAt)
	if err != nil {
		return entities.Hackathon{}, err
	}
	ends, err := parseTime(t.EventEndsAt)
	if err != nil {
		return entities.Hackathon{}, err
	}
	deadline, err := parseTime(t.SubmissionDeadlineAt)
	if err != nil {
		return entities.Hackathon{}, err
	}

	return entities.Hackathon{
		RegistrationOpensAt:  opens,
		RegistrationClosesAt: closes,
		EventStartsAt:        starts,
		EventEndsAt:          ends,
		SubmissionDeadlineAt: deadline,
	}, nil
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func NullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
