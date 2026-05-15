package hackathon_usecase

import (
	"context"
	"strings"
	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/converters"
	"github.com/nikitatisenko/pirksp/internal/domain/dto"
	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
	"github.com/nikitatisenko/pirksp/internal/errs"
	"github.com/nikitatisenko/pirksp/internal/usecase"
)

var defaultPublicStatuses = []string{"registration", "running", "finished"}

type hackathonUseCase struct {
	repo      repository.HackathonsRepository
	converter *converters.HackathonsConverter
}

func NewHackathonUseCase(repo repository.HackathonsRepository) usecase.HackathonUseCase {
	return &hackathonUseCase{
		repo:      repo,
		converter: converters.NewHackathonsConverter(),
	}
}

func (u *hackathonUseCase) ListPublic(ctx context.Context, statuses []string, query string) ([]dto.HackathonListItem, error) {
	if len(statuses) == 0 {
		statuses = defaultPublicStatuses
	}

	items, err := u.repo.List(ctx, repository.HackathonListFilter{
		Statuses:     statuses,
		Query:        strings.TrimSpace(query),
		ExcludeDraft: true,
	})
	if err != nil {
		return nil, err
	}

	out := make([]dto.HackathonListItem, 0, len(items))
	for _, h := range items {
		out = append(out, u.converter.ToListItem(h))
	}

	return out, nil
}

func (u *hackathonUseCase) GetPublic(ctx context.Context, id uuid.UUID) (*dto.HackathonDetail, error) {
	h, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if h.Status == "draft" {
		return nil, errs.ErrNotFound
	}

	return u.loadDetail(ctx, h)
}

func (u *hackathonUseCase) GetOrganizer(ctx context.Context, organizerID, hackathonID uuid.UUID) (*dto.HackathonDetail, error) {
	h, err := u.ensureOwner(ctx, organizerID, hackathonID)
	if err != nil {
		return nil, err
	}
	return u.loadDetail(ctx, h)
}

func (u *hackathonUseCase) ListOrganizer(ctx context.Context, organizerID uuid.UUID) ([]dto.HackathonListItem, error) {
	items, err := u.repo.List(ctx, repository.HackathonListFilter{
		OrganizerID: &organizerID,
	})
	if err != nil {
		return nil, err
	}

	out := make([]dto.HackathonListItem, 0, len(items))
	for _, h := range items {
		out = append(out, u.converter.ToListItem(h))
	}

	return out, nil
}

func (u *hackathonUseCase) Create(ctx context.Context, organizerID uuid.UUID, req dto.CreateHackathonRequest) (*dto.HackathonDetail, error) {
	h, err := u.buildHackathonFromCreate(organizerID, req)
	if err != nil {
		return nil, err
	}

	created, err := u.repo.Create(ctx, h)
	if err != nil {
		return nil, err
	}

	return u.loadDetail(ctx, created)
}

func (u *hackathonUseCase) Update(ctx context.Context, organizerID, hackathonID uuid.UUID, req dto.UpdateHackathonRequest) (*dto.HackathonDetail, error) {
	h, err := u.ensureOwner(ctx, organizerID, hackathonID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		h.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		h.Description = strings.TrimSpace(*req.Description)
	}
	if req.ShortDescription != nil {
		h.ShortDescription = converters.NullString(strings.TrimSpace(*req.ShortDescription))
	}
	if req.Format != nil {
		format := strings.TrimSpace(*req.Format)
		if format != "" {
			h.Format = format
		}
	}
	if req.PrizesInfo != nil {
		h.PrizesInfo = converters.NullString(strings.TrimSpace(*req.PrizesInfo))
	}
	if req.MaxTeamSize != nil && *req.MaxTeamSize >= 2 && *req.MaxTeamSize <= 8 {
		h.MaxTeamSize = *req.MaxTeamSize
	}
	if req.Timeline != nil {
		timeline, err := u.converter.TimelineToEntity(*req.Timeline)
		if err != nil {
			return nil, errs.ErrValidation
		}
		h.RegistrationOpensAt = timeline.RegistrationOpensAt
		h.RegistrationClosesAt = timeline.RegistrationClosesAt
		h.EventStartsAt = timeline.EventStartsAt
		h.EventEndsAt = timeline.EventEndsAt
		h.SubmissionDeadlineAt = timeline.SubmissionDeadlineAt
	}
	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)
		if status == "running" || status == "finished" {
			h.Status = status
		}
	}

	if h.Title == "" || h.Description == "" {
		return nil, errs.ErrValidation
	}

	updated, err := u.repo.Update(ctx, h)
	if err != nil {
		return nil, err
	}

	return u.loadDetail(ctx, updated)
}

func (u *hackathonUseCase) Delete(ctx context.Context, organizerID, hackathonID uuid.UUID) error {
	h, err := u.ensureOwner(ctx, organizerID, hackathonID)
	if err != nil {
		return err
	}
	if h.Status != "draft" {
		return errs.ErrCannotDelete
	}

	return u.repo.Delete(ctx, hackathonID)
}

func (u *hackathonUseCase) Publish(ctx context.Context, organizerID, hackathonID uuid.UUID) (*dto.HackathonDetail, error) {
	h, err := u.ensureOwner(ctx, organizerID, hackathonID)
	if err != nil {
		return nil, err
	}
	if h.Status != "draft" {
		return nil, errs.ErrValidation
	}

	count, err := u.repo.CountCasesByHackathon(ctx, hackathonID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errs.ErrPublishNotReady
	}

	updated, err := u.repo.UpdateStatus(ctx, hackathonID, "registration")
	if err != nil {
		return nil, err
	}

	return u.loadDetail(ctx, updated)
}

func (u *hackathonUseCase) ListTracks(ctx context.Context, organizerID, hackathonID uuid.UUID) ([]dto.Track, error) {
	if _, err := u.ensureOwner(ctx, organizerID, hackathonID); err != nil {
		return nil, err
	}

	tracks, err := u.repo.ListTracks(ctx, hackathonID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.Track, 0, len(tracks))
	for _, t := range tracks {
		out = append(out, u.converter.ToTrack(t))
	}

	return out, nil
}

func (u *hackathonUseCase) CreateTrack(ctx context.Context, organizerID, hackathonID uuid.UUID, req dto.CreateTrackRequest) (*dto.Track, error) {
	if _, err := u.ensureOwner(ctx, organizerID, hackathonID); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errs.ErrValidation
	}

	created, err := u.repo.CreateTrack(ctx, entities.Track{
		HackathonID: hackathonID,
		Title:       title,
		Description: converters.NullString(strings.TrimSpace(req.Description)),
	})
	if err != nil {
		return nil, err
	}

	track := u.converter.ToTrack(created)
	return &track, nil
}

func (u *hackathonUseCase) UpdateTrack(ctx context.Context, organizerID, trackID uuid.UUID, req dto.CreateTrackRequest) (*dto.Track, error) {
	track, err := u.repo.GetTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}
	if _, err := u.ensureOwner(ctx, organizerID, track.HackathonID); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errs.ErrValidation
	}

	track.Title = title
	track.Description = converters.NullString(strings.TrimSpace(req.Description))

	updated, err := u.repo.UpdateTrack(ctx, track)
	if err != nil {
		return nil, err
	}

	out := u.converter.ToTrack(updated)
	return &out, nil
}

func (u *hackathonUseCase) DeleteTrack(ctx context.Context, organizerID, trackID uuid.UUID) error {
	track, err := u.repo.GetTrack(ctx, trackID)
	if err != nil {
		return err
	}
	if _, err := u.ensureOwner(ctx, organizerID, track.HackathonID); err != nil {
		return err
	}

	return u.repo.DeleteTrack(ctx, trackID)
}

func (u *hackathonUseCase) ListCases(ctx context.Context, organizerID, trackID uuid.UUID) ([]dto.Case, error) {
	track, err := u.repo.GetTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}
	if _, err := u.ensureOwner(ctx, organizerID, track.HackathonID); err != nil {
		return nil, err
	}

	cases, err := u.repo.ListCasesByTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.Case, 0, len(cases))
	for _, c := range cases {
		out = append(out, u.converter.ToCase(c))
	}

	return out, nil
}

func (u *hackathonUseCase) CreateCase(ctx context.Context, organizerID, trackID uuid.UUID, req dto.CreateCaseRequest) (*dto.Case, error) {
	track, err := u.repo.GetTrack(ctx, trackID)
	if err != nil {
		return nil, err
	}
	if _, err := u.ensureOwner(ctx, organizerID, track.HackathonID); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	if title == "" || description == "" {
		return nil, errs.ErrValidation
	}

	created, err := u.repo.CreateCase(ctx, entities.Case{
		TrackID:      trackID,
		Title:        title,
		Description:  description,
		CustomerName: converters.NullString(strings.TrimSpace(req.CustomerName)),
		ResourcesURL: converters.NullString(strings.TrimSpace(req.ResourcesURL)),
	})
	if err != nil {
		return nil, err
	}

	out := u.converter.ToCase(created)
	return &out, nil
}

func (u *hackathonUseCase) UpdateCase(ctx context.Context, organizerID, caseID uuid.UUID, req dto.CreateCaseRequest) (*dto.Case, error) {
	caze, err := u.repo.GetCase(ctx, caseID)
	if err != nil {
		return nil, err
	}
	track, err := u.repo.GetTrack(ctx, caze.TrackID)
	if err != nil {
		return nil, err
	}
	if _, err := u.ensureOwner(ctx, organizerID, track.HackathonID); err != nil {
		return nil, err
	}

	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	if title == "" || description == "" {
		return nil, errs.ErrValidation
	}

	caze.Title = title
	caze.Description = description
	caze.CustomerName = converters.NullString(strings.TrimSpace(req.CustomerName))
	caze.ResourcesURL = converters.NullString(strings.TrimSpace(req.ResourcesURL))

	updated, err := u.repo.UpdateCase(ctx, caze)
	if err != nil {
		return nil, err
	}

	out := u.converter.ToCase(updated)
	return &out, nil
}

func (u *hackathonUseCase) DeleteCase(ctx context.Context, organizerID, caseID uuid.UUID) error {
	caze, err := u.repo.GetCase(ctx, caseID)
	if err != nil {
		return err
	}
	track, err := u.repo.GetTrack(ctx, caze.TrackID)
	if err != nil {
		return err
	}
	if _, err := u.ensureOwner(ctx, organizerID, track.HackathonID); err != nil {
		return err
	}

	return u.repo.DeleteCase(ctx, caseID)
}

func (u *hackathonUseCase) ensureOwner(ctx context.Context, organizerID, hackathonID uuid.UUID) (entities.Hackathon, error) {
	h, err := u.repo.GetByID(ctx, hackathonID)
	if err != nil {
		return entities.Hackathon{}, err
	}
	if h.OrganizerID != organizerID {
		return entities.Hackathon{}, errs.ErrNotOwner
	}
	return h, nil
}

func (u *hackathonUseCase) loadDetail(ctx context.Context, h entities.Hackathon) (*dto.HackathonDetail, error) {
	tracks, err := u.repo.ListTracks(ctx, h.ID)
	if err != nil {
		return nil, err
	}
	cases, err := u.repo.ListCasesByHackathon(ctx, h.ID)
	if err != nil {
		return nil, err
	}

	detail := u.converter.ToDetail(h, tracks, cases)
	return &detail, nil
}

func (u *hackathonUseCase) buildHackathonFromCreate(organizerID uuid.UUID, req dto.CreateHackathonRequest) (entities.Hackathon, error) {
	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	if title == "" || description == "" {
		return entities.Hackathon{}, errs.ErrValidation
	}

	timeline, err := u.converter.TimelineToEntity(req.Timeline)
	if err != nil {
		return entities.Hackathon{}, errs.ErrValidation
	}

	format := strings.TrimSpace(req.Format)
	if format == "" {
		format = "online"
	}

	maxSize := req.MaxTeamSize
	if maxSize == 0 {
		maxSize = 5
	}
	if maxSize < 2 || maxSize > 8 {
		return entities.Hackathon{}, errs.ErrValidation
	}

	return entities.Hackathon{
		OrganizerID:            organizerID,
		Title:                  title,
		ShortDescription:       converters.NullString(strings.TrimSpace(req.ShortDescription)),
		Description:            description,
		Format:                 format,
		Status:                 "draft",
		MaxTeamSize:            maxSize,
		PrizesInfo:             converters.NullString(strings.TrimSpace(req.PrizesInfo)),
		RegistrationOpensAt:    timeline.RegistrationOpensAt,
		RegistrationClosesAt:   timeline.RegistrationClosesAt,
		EventStartsAt:          timeline.EventStartsAt,
		EventEndsAt:            timeline.EventEndsAt,
		SubmissionDeadlineAt:   timeline.SubmissionDeadlineAt,
	}, nil
}
