package dtos

type UpdateTodoRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Complete    *bool  `json:"complete,omitempty"`
}