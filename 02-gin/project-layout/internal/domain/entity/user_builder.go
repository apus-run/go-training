// 推荐后面用工具自动生成
// 使用方式
// userEntity := NewUserBuilder().SetID(id).SetName('name').Build()

package entity

import "time"

type UserBuilder struct {
	id       uint64
	name     string // 账户名
	avatar   string
	email    string
	password string
	phone    string

	gender   int       // 性别
	nickName string    // 昵称
	realName string    // 真实姓名
	birthday time.Time // 生日
	profile  string    // 个人简介

	createdTime time.Time  // 创建时间
	updatedTime time.Time  // 更新时间
	deletedTime *time.Time // 删除时间
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (b *UserBuilder) ID(id uint64) *UserBuilder {
	b.id = id
	return b
}

func (b *UserBuilder) Name(name string) *UserBuilder {
	b.name = name
	return b
}

func (b *UserBuilder) Avatar(avatar string) *UserBuilder {
	b.avatar = avatar
	return b
}

func (b *UserBuilder) Email(email string) *UserBuilder {
	b.email = email
	return b
}

func (b *UserBuilder) Password(password string) *UserBuilder {
	b.password = password
	return b
}

func (b *UserBuilder) Phone(phone string) *UserBuilder {
	b.phone = phone
	return b
}

func (b *UserBuilder) Gender(getter int) *UserBuilder {
	b.gender = getter
	return b
}

func (b *UserBuilder) NickName(nickName string) *UserBuilder {
	b.nickName = nickName
	return b
}

func (b *UserBuilder) RealName(realName string) *UserBuilder {
	b.realName = realName
	return b
}

func (b *UserBuilder) Birthday(birthday time.Time) *UserBuilder {
	b.birthday = birthday
	return b
}

func (b *UserBuilder) Profile(profile string) *UserBuilder {
	b.profile = profile
	return b
}

func (b *UserBuilder) CreatedTime(createdTime time.Time) *UserBuilder {
	b.createdTime = createdTime
	return b
}

func (b *UserBuilder) Build() *User {
	return &User{
		id:          b.id,
		name:        b.name,
		avatar:      b.avatar,
		email:       b.email,
		password:    b.password,
		phone:       b.phone,
		gender:      b.gender,
		nickName:    b.nickName,
		realName:    b.realName,
		birthday:    b.birthday,
		profile:     b.profile,
		createdTime: b.createdTime,
		updatedTime: b.updatedTime,
		deletedTime: b.deletedTime,
	}
}
