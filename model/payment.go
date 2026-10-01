package model

type Payment interface{

	Pay(amount float64) error 
}