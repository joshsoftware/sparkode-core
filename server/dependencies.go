package server

import (
	"github.com/joshsoftware/sparkode-core/app"
	"github.com/joshsoftware/sparkode-core/service"
)

type dependencies struct {
	RunService service.Service
}

func InitDependencies() (dependencies, error) {
	logger := app.GetLogger()

	// call new service
	runService := service.NewService(logger)

	return dependencies{
		RunService: runService,
	}, nil
}
