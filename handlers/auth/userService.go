package auth

import "db-service/internal/database"

type UserService struct {
	db *database.Queries
}

func NewService() *UserService {
	return &UserService{}
}

func (us *UserService) WithDb(db *database.Queries) *UserService {
	us.db = db
	return us
}
