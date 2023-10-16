package web

type CategoryJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=2,max=100"`
}
