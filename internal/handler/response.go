package handler

import "github.com/rtsoy/todo-app/internal/model"

type resourceResponse struct {
	Count      int              `json:"count"`
	Results    any              `json:"results"`
	Pagination model.Pagination `json:"pagination"`
}
