package auth

// UsecaseInteractor godoc
type UsecaseInteractor interface {
	IssueToken(LoginRequest) (TokenResponse, error)
	RefreshToken(AccessTokenByRefreshRequest) (TokenResponse, error)
}
