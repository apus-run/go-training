package main

import (
	"time"
)

type User struct {
	ID       uint64
	Name     string // 账户名
	Email    string
	Password string
	Phone    string

	CreatedTime time.Time  // 创建时间
	UpdatedTime time.Time  // 更新时间
	DeletedTime *time.Time // 删除时间
}

func main() {

	// user := &User{}
	//u := &User{} // u := new(User)
	//fmt.Println(user)
	//fmt.Println(u)
	//user.ID = 1
	//fmt.Println(user)
	//u.ID = 2
	//fmt.Println(u)

	//userModel := User{}
	//userModel, _ = userModel.FromEntity().(User)
	//fmt.Println(userModel)

	//var user User
	//user.ID = 1
	//fmt.Printf("%+v\n", user)
	//userModel := User{}
	//userModel.ID = 2
	//fmt.Printf("%+v\n", userModel)
	//if strings.HasPrefix(":9000", ":") {
	//	fmt.Printf(":9000")
	//} else {
	//	fmt.Printf("9000")
	//}

	//hash, err := user.GenerateHashPassword("123456")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(hash)
	//fmt.Println(user.VerifyPassword("$2a$10$18ojWf3/xEPxU0vlcjQjv.EMle.5EepSaw66G9fR0c0zlRBSnsW3e", "123456"))

	//isEmail, err := regexp.
	//	MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, regexp.None).
	//	MatchString("moocss@163.com")
	//if !isEmail || err != nil {
	//	fmt.Printf("邮箱不合法\n")
	//} else {
	//	fmt.Printf("邮箱合法\n")
	//}
	//
	//isPassword, err := regexp.
	//	MustCompile(`^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`, regexp.None).
	//	MatchString("www#123456")
	//if !isPassword || err != nil {
	//	fmt.Printf("密码不合法\n")
	//} else {
	//	fmt.Printf("密码合法\n")
	//}

}

//func (u *User) GenerateHashPassword(password string) (string, error) {
//	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		return "", err
//	}
//	pass := string(hash)
//	return pass, nil
//}
//
//// VerifyPassword 验证密码
//func (u *User) VerifyPassword(hashedPassword, password string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//	return err == nil
//}
