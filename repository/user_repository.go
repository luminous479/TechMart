package repository

import (
	"database/sql"

	"github.com/luminous479/TechMart/model"
)


type UserRepository struct{
	db *sql.DB
} 


func NewUserRepository(db *sql.DB) *UserRepository{
	return &UserRepository{
		db: db,
	}
}

func ( userRepo *UserRepository) GetByID(id int) (*model.User,error){

	var user model.User

	err := userRepo.db.QueryRow("SELECT id,name,email,created_at FROM users WHERE id = $1",
	id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	    &user.CreatedAt)

		if err != nil{
			return nil, err
		}
	return &user, nil

}