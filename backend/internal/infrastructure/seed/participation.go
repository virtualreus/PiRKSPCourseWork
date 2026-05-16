package seed

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/nikitatisenko/pirksp/internal/domain/entities"
	"github.com/nikitatisenko/pirksp/internal/domain/ports/repository"
)

func seedParticipation(
	ctx context.Context,
	users repository.UsersRepository,
	hackathons repository.HackathonsRepository,
	participation repository.ParticipationRepository,
	log *slog.Logger,
) {
	participant, err := users.GetByEmail(ctx, "user@user.ru")
	if err != nil {
		return
	}

	rows, err := participation.ListUserParticipations(ctx, participant.ID)
	if err != nil || len(rows) > 0 {
		return
	}

	list, err := hackathons.List(ctx, repository.HackathonListFilter{ExcludeDraft: true})
	if err != nil {
		log.Warn("seed: list hackathons for participation", "err", err)
		return
	}

	for _, h := range list {
		var teamName string
		var withSubmission bool
		switch h.Title {
		case "Цифровой город 2026":
			teamName = "UrbanMinds"
			withSubmission = false
		case "FinTech Product Sprint":
			teamName = "PayFlow"
			withSubmission = true
		default:
			continue
		}

		if err := seedParticipantOnHackathon(ctx, hackathons, participation, participant.ID, h, teamName, withSubmission); err != nil {
			log.Warn("seed: participation", "hackathon", h.Title, "err", err)
		}
	}

	log.Info("seed: demo participation ready")
}

func seedParticipantOnHackathon(
	ctx context.Context,
	hackathons repository.HackathonsRepository,
	participation repository.ParticipationRepository,
	userID uuid.UUID,
	h entities.Hackathon,
	teamName string,
	withSubmission bool,
) error {
	_, err := participation.CreateRegistration(ctx, entities.HackathonRegistration{
		HackathonID: h.ID,
		UserID:      userID,
	})
	if err != nil {
		return err
	}

	tracks, err := hackathons.ListTracks(ctx, h.ID)
	if err != nil || len(tracks) == 0 {
		return err
	}

	track := tracks[0]
	cases, err := hackathons.ListCasesByTrack(ctx, track.ID)
	if err != nil || len(cases) == 0 {
		return err
	}
	caze := cases[0]

	team, err := participation.CreateTeam(ctx, entities.Team{
		HackathonID: h.ID,
		Name:        teamName,
		CaptainID:   userID,
		TrackID:     sql.NullString{String: track.ID.String(), Valid: true},
		CaseID:      sql.NullString{String: caze.ID.String(), Valid: true},
	})
	if err != nil {
		return err
	}

	if err := participation.AddMember(ctx, team.ID, userID, "team_lead"); err != nil {
		return err
	}

	if !withSubmission {
		return nil
	}

	now := time.Now().UTC()
	_, err = participation.UpsertSubmission(ctx, entities.Submission{
		TeamID:      team.ID,
		HackathonID: h.ID,
		Title:       sql.NullString{String: teamName + " MVP", Valid: true},
		Summary:     sql.NullString{String: "Демо-решение для сценария P2P-перевода.", Valid: true},
		RepoURL:     "https://github.com/demo/payflow-hackathon",
		DemoURL:     sql.NullString{String: "https://payflow-demo.example.com", Valid: true},
		SubmittedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt:   now,
	})
	return err
}
