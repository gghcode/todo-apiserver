package auth

// UseCase godoc
type UseCase interface {
	IssueToken(LoginRequest) (TokenResponse, error)
	RefreshToken(AccessTokenByRefreshRequest) (TokenResponse, error)
}
