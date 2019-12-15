package auth

import "github.com/gghcode/apas-todo-apiserver/domain/auth"

type (
	tokenResponseDTO struct {
		Type         string `json:"type"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int64  `json:"expires_in"`
	}

	tokenResponseSerializer struct {
		model auth.TokenResponse
	}
)

// newTokenResponseSerializer serialize token response
func newTokenResponseSerializer(model auth.TokenResponse) *tokenResponseSerializer {
	return &tokenResponseSerializer{
		model: model,
	}
}

func (s *tokenResponseSerializer) Response() tokenResponseDTO {
	return tokenResponseDTO{
		Type:         s.model.Type,
		AccessToken:  s.model.AccessToken,
		RefreshToken: s.model.RefreshToken,
		ExpiresIn:    s.model.ExpiresIn,
	}
}
