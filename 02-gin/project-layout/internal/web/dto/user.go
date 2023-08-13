package dto

// UserRequest 请求数据结果
type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse 返回数据结构
type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// RegisterRequest 注册请求数据结构
type RegisterRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UpdateProfileRequest struct {
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	RealName string `json:"real_name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}

type UpdateProfileResponse struct {
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	RealName string `json:"real_name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}

type ProfileResponse struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	RealName string `json:"real_name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}
