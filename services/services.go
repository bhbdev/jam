package services

type Services struct {
	UserService   *UserService
	JobAppService *JobAppService
}

func NewServices(
	userService *UserService,
	jobAppService *JobAppService,
) *Services {
	return &Services{
		UserService:   userService,
		JobAppService: jobAppService,
	}
}
