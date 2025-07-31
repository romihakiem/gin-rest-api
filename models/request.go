package models

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
}

type PostRequest struct {
	Title      string `json:"title" binding:"required,min=2,max=200"`
	Body       string `json:"body" binding:"required"`
	CategoryId uint   `json:"categoryId" binding:"required,min=1"`
}

type CategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type CommentAddRequest struct {
	PostId uint   `json:"postId" binding:"required,min=1"`
	Body   string `json:"body" binding:"required,min=1"`
}

type CommentEditRequest struct {
	Body string `json:"body" binding:"required,min=1"`
}
