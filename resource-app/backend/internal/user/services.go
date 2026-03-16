package user

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUsers() ([]User, error) {
	return s.repo.GetUsers()
}

func (s *Service) GetUserByEmail(email string) (*User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *Service) GetUserByID(id string) (*User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) CreateUser(user *User) error {
	return s.repo.CreateUser(user)
}

func (s *Service) UpdateUserRole(userID string, role Role) error {
	return s.repo.UpdateUserRole(userID, role)
}
