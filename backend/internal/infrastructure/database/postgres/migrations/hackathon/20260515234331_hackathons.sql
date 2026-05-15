-- +goose Up
-- +goose StatementBegin
CREATE TABLE hackathons (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organizer_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title                   VARCHAR(255) NOT NULL,
    short_description       TEXT,
    description             TEXT NOT NULL,
    format                  VARCHAR(16) NOT NULL DEFAULT 'online'
        CHECK (format IN ('online', 'offline', 'hybrid')),
    status                  VARCHAR(32) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'registration', 'running', 'finished')),
    max_team_size           INT NOT NULL DEFAULT 5
        CHECK (max_team_size >= 2 AND max_team_size <= 8),
    prizes_info             TEXT,
    registration_opens_at   TIMESTAMPTZ NOT NULL,
    registration_closes_at  TIMESTAMPTZ NOT NULL,
    event_starts_at         TIMESTAMPTZ NOT NULL,
    event_ends_at           TIMESTAMPTZ NOT NULL,
    submission_deadline_at  TIMESTAMPTZ NOT NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tracks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hackathon_id    UUID NOT NULL REFERENCES hackathons(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    track_id        UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT NOT NULL,
    customer_name   VARCHAR(255),
    resources_url   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_hackathons_status ON hackathons(status);
CREATE INDEX idx_hackathons_organizer ON hackathons(organizer_id);
CREATE INDEX idx_tracks_hackathon ON tracks(hackathon_id);
CREATE INDEX idx_cases_track ON cases(track_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cases;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS hackathons;
-- +goose StatementEnd
