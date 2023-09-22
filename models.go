package main

type Task struct {
	ID          int    `json:"id" validate:"nonzero"`
	Title       string `json:"title" validate:"nonzero"`
	Description string `json:"description" validate:"nonzero"`
	Status      string `json:"status" validate:"nonzero"`
}

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
