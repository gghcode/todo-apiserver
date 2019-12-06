package fake

import (
	"github.com/stretchr/testify/mock"
)

// AppService is fake app service
type AppService struct {
	mock.Mock
}

// NewAppService return new app service
func NewAppService() *AppService {
	return &AppService{}
}

// Version return stub version
func (srv *AppService) Version() string {
	args := srv.Called()
	return args.String(0)
}
