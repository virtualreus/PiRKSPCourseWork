package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	PlatformRole string    `db:"platform_role"`
	CreatedAt    time.Time `db:"created_at"`
}
