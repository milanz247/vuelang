package requests

// CreateUserRequest is the body for POST /api/v1/users.
type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// UpdateUserRequest is the body for PUT /api/v1/users/:id.
type UpdateUserRequest struct {
	Name     string `json:"name"      binding:"required,min=2,max=100"`
	Email    string `json:"email"     binding:"required,email,max=150"`
	IsActive bool   `json:"is_active"`
}
