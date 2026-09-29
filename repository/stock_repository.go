package repository

import "database/sql"


type stockRepository struct{
	db *sql.DB
}

func NewstockRepository(db *sql.DB) *stockRepository{
	return &stockRepository{
		db: db,
	}
}

