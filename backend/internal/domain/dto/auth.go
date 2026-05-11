package dto

type RegisterRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FullName     string `json:"full_name"`
	PlatformRole string `json:"platform_role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	FullName *string `json:"full_name"`
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	PlatformRole string `json:"platform_role"`
	CreatedAt    string `json:"created_at"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}
