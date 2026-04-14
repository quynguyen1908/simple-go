package constants

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

const (
	UserStatusActive = "active"
)

const (
	ProviderSystem = "system"
	ProviderGoogle = "google"
)

const (
	TokenTypeEmailConfirmation = "email_confirmation"
	TokenTypeRefresh           = "refresh_token"
)

const (
	ErrBadRequest     = "Bad Request"
	ErrConflict       = "Conflict"
	ErrNotFound       = "Not Found"
	ErrInternal       = "Internal Server Error"
	ErrUnauthorized   = "Unauthorized"
	ErrForbidden      = "Forbidden"
	ErrRequestTimeout = "Request Timeout"
)
