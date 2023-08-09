package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID   uint64
	Name string
}

type EUser interface{}

type TUser struct {
	ID   uint64
	Name string
}

func (u *User) ToEntity() {
	fmt.Println("ToEntity", u == nil)
}

func (u *User) FromEntity() any {
	return User{
		ID:   1,
		Name: "kami",
	}
}

func main() {
	user := &User{}
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

	hash, err := user.GenerateHashPassword("111111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hash)
	fmt.Println(user.VerifyPassword("$2a$10$aQG9ikUEMb.lOHSz8gZC7u2otmykBebB0C2iRT.WssZZmgtdzRnfq", "111111"))
}

func (u *User) GenerateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	pass := string(hash)
	return pass, nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
