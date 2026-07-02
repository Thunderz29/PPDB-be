package request

type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"user_password" binding:"required"`
}

type CreateUserRequest struct {
	RoleID       *int64  `json:"role_id"`
	UserFullName string  `json:"user_full_name" binding:"required"`
	UserName     string  `json:"user_name" binding:"required"`
	UserPhone    *string `json:"user_phone"`
	UserEmail    string  `json:"user_email" binding:"required,email"`
	UserPassword string  `json:"user_password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	RoleID       *int64  `json:"role_id"`
	UserFullName string  `json:"user_full_name" binding:"required"`
	UserName     string  `json:"user_name" binding:"required"`
	UserPhone    *string `json:"user_phone"`
	UserEmail    string  `json:"user_email" binding:"required,email"`
	UserPassword string  `json:"user_password"` // Optional on update
}
