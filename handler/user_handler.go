package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/luminous479/TechMart/helper"
	"github.com/luminous479/TechMart/service"
)

type UserHandler struct{

	serv *service.UserService
}
func NewUserHandler(sr *service.UserService) *UserHandler{
	return &UserHandler{
		serv: sr,
	}
}
 
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id , err := strconv.Atoi(r.PathValue("id"))

	if err != nil{
		helper.WriteJSONError(w,"Invalid ID",
		http.StatusBadRequest)
		return
	}

	user,  err := h.serv.GetByID(id)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows){
			helper.
			WriteJSONError(w,
	        "User Not Found!",
			http.StatusNotFound)
            return
		}
		helper.
		WriteJSONError(w,
		"Failed to get User",
		 http.StatusInternalServerError)
		 return
		
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}