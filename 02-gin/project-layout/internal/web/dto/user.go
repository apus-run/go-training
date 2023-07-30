package dto

// UserRequest 请求数据结果
type UserRequest struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse 返回数据结构
type UserResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
