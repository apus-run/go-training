package main

import (
	"fmt"
	"net/http"
	"time"
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
	url := fmt.Sprintf("%s/ping", "http://localhost:9000")

	if err := Ping(url, 5); err != nil {
		fmt.Printf("server no response: %s", err)
	}
	fmt.Printf("server started success!")

}

// Ping 用来检查是否程序正常启动
func Ping(addr string, maxCount int) error {
	seconds := 1
	fmt.Printf("地址: %s\n", addr)
	url := fmt.Sprintf("%s/ping", addr)
	for i := 0; i < maxCount; i++ {
		resp, err := http.Get(url)
		if nil == err && resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		fmt.Printf("等待服务在线, 已等待 %d 秒，最多等待 %d 秒 \n", seconds, maxCount)
		time.Sleep(time.Second * 1)
		seconds++
	}
	return fmt.Errorf("服务启动失败: %s", addr)
}
