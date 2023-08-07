package main

import (
	"github.com/sirupsen/logrus"

	"project-layout/cmd/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		logrus.Fatal(err)
	}
}
