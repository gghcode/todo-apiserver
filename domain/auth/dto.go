package auth

// LoginRequest godoc
type LoginRequest struct {
	Username string
	Password string
}

// AccessTokenByRefreshRequest godoc
type AccessTokenByRefreshRequest struct {
	Token string
}

// TokenResponse godoc
type TokenResponse struct {
	Type         string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}
