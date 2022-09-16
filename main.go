package main

// @APITitle Main
// @APIDescription Main API for Microservices in Go!

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joshsoftware/sparkode-core/config"
	"github.com/joshsoftware/sparkode-core/service"

	logger "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"github.com/urfave/negroni"
)

func main() {
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	})

	config.Load()

	cliApp := cli.NewApp()
	cliApp.Name = config.AppName()
	cliApp.Version = "1.0.0"
	cliApp.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "start server",
			Action: func(c *cli.Context) error {
				return startApp()
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func startApp() (err error) {

	// mux router
	router := service.InitRouter()

	// init web server
	server := negroni.Classic()
	server.UseHandler(router)

	port := config.AppPort() // This can be changed to the service port number via environment variable.
	addr := fmt.Sprintf(":%s", strconv.Itoa(port))

	server.Run(addr)
	return
}
