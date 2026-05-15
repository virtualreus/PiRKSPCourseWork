package errs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrDuplicate     = errors.New("duplicate")
	ErrEmailTaken    = errors.New("email already registered")
	ErrInvalidJSON   = errors.New("invalid json body")
	ErrValidation    = errors.New("validation failed")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrInvalidRole   = errors.New("invalid platform role")
	ErrWeakPassword  = errors.New("password must be at least 8 characters")
	ErrEmptyEmail    = errors.New("email is required")
	ErrEmptyName     = errors.New("full name is required")
	ErrInvalidToken  = errors.New("invalid token")
	ErrForbidden     = errors.New("forbidden")
	ErrNotOwner      = errors.New("not owner of hackathon")
	ErrCannotDelete  = errors.New("cannot delete published hackathon")
	ErrPublishNotReady = errors.New("at least one track with a case is required")
)
