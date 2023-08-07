package uuid

import (
	"github.com/google/uuid"
)

type GenUUID struct{}

func New() *GenUUID {
	return &GenUUID{}
}

func (u *GenUUID) New() string {
	return uuid.New().String()
}

func Parse(value string) (string, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
