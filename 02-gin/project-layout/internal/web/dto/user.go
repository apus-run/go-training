package dto

// UserRequest 请求数据结果
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse 返回数据结构
type UserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required, password"`
	ConfirmPassword string `json:"confirmPassword" binding:"required, password"`
}
