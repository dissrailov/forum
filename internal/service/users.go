package service

func (s *service) CreateUser(name, email, password string) error {
	return s.repo.CreateUser(name, email, password)
}

func (s *service) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}
func (s *service) Authenticate(email, password string) (int, error) {
	return s.repo.Authenticate(email, password)
}
