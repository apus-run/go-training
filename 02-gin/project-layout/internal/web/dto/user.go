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

// RegisterRequest 注册请求数据结构
type RegisterRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
