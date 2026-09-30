package service

import "github.com/luminous479/TechMart/model"

type UserRepository interface{

	GetByID(id string) (*model.User,error)
}