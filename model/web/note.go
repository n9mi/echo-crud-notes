package web

type NoteRequest struct {
	ID         int    `json:"id"`
	Title      string `json:"title" validate:"required,min=2,max=100"`
	Body       string `json:"body" validate:"required,min=2,max=255"`
	CategoryId int    `json:"id_category" validate:"required,gte=0"`
}

type NoteResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Category string `json:"category"`
}
