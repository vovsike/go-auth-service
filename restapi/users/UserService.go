package users

type UserService struct {
	Store UserStore
}

func NewUserService(store UserStore) *UserService {
	return &UserService{Store: store}
}

func (s *UserService) GetAllUsers() []User {
	return s.Store.GetAll()
}
