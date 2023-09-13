package svc

type UserService interface {
	Login() error
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (svc *userService) Login() error {
	return nil
}
