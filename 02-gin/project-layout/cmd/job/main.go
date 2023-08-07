package main

import (
	"project-layout/cmd/job/cron"
)

func main() {
	if err := cron.Run(); err != nil {

	}
}
