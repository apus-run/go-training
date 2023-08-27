//go:build wireinject

package main

import (
	"github.com/google/wire"

	"project-layout/pkg/ginx"
)

// wireApp init web application.
func wireApp(logger *log.Logger) (*ginx.HttpServer, func(), error) {
	panic(wire.Build())
}
