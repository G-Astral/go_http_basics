package models

type UpdateUserInput struct {
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}
