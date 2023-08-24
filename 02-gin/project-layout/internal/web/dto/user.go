package dto

// LoginReq 请求数据结果
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResp 返回数据结构
type LoginResp struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

// RegisterReq 注册请求数据结构
type RegisterReq struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UpdateProfileReq struct {
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	RealName string `json:"real_name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}

type UpdateProfileResp struct {
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	RealName string `json:"real_name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}

type UserInfoResp struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Gender   int    `json:"gender"`
	NickName string `json:"nick_name"`
	RealName string `json:"real_name"`
	Birthday string `json:"birthday"`
	Profile  string `json:"profile"`
}
