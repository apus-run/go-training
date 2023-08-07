package main

import "fmt"

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
	//user := User{}
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

	var user User
	user.ID = 1
	fmt.Printf("%+v\n", user)
	userModel := User{}
	userModel.ID = 2
	fmt.Printf("%+v\n", userModel)

}
