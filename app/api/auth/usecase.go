package auth

// Service godoc
type Service interface {
	IssueToken(LoginRequest, *TokenResponse) error
	RefreshToken(AccessTokenByRefreshRequest, *TokenResponse) error
}

// LoginRequest godoc
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}

// AccessTokenByRefreshRequest godoc
type AccessTokenByRefreshRequest struct {
	Token string `json:"token" validate:"required"`
}

// TokenResponse godoc
type TokenResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}
