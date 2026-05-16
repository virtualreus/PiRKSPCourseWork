-- +goose Up
-- +goose StatementBegin
CREATE TABLE hackathon_registrations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hackathon_id    UUID NOT NULL REFERENCES hackathons(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    registered_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (hackathon_id, user_id)
);

CREATE TABLE teams (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hackathon_id    UUID NOT NULL REFERENCES hackathons(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    captain_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id        UUID REFERENCES tracks(id) ON DELETE SET NULL,
    case_id         UUID REFERENCES cases(id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE team_members (
    team_id         UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_role       VARCHAR(32) NOT NULL CHECK (team_role IN (
        'team_lead', 'developer', 'designer', 'data_scientist', 'devops_qa', 'other'
    )),
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (team_id, user_id)
);

CREATE TABLE submissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id         UUID NOT NULL UNIQUE REFERENCES teams(id) ON DELETE CASCADE,
    hackathon_id    UUID NOT NULL REFERENCES hackathons(id) ON DELETE CASCADE,
    title           VARCHAR(255),
    summary         TEXT,
    repo_url        TEXT NOT NULL,
    demo_url        TEXT,
    pitch_url       TEXT,
    video_url       TEXT,
    submitted_at    TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_registrations_hackathon ON hackathon_registrations(hackathon_id);
CREATE INDEX idx_registrations_user ON hackathon_registrations(user_id);
CREATE INDEX idx_teams_hackathon ON teams(hackathon_id);
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_submissions_hackathon ON submissions(hackathon_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS hackathon_registrations;
-- +goose StatementEnd
