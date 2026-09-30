package service

import "github.com/luminous479/TechMart/model"

type UserRepository interface{

	GetByID(id int) (*model.User,error)
}


type UserService struct{

	repo UserRepository
} 

func NewUserService(repo UserRepository) *UserService{
	return &UserService{
		repo: repo,
	}
}

func ( sr UserService) GetByID (id int) (*model.User, error){

	return sr.repo.GetByID(id)

}