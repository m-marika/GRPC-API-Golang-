package main

import (
	"os"

	server "calcLab2/calcLab/server"

	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	//"ngs-core/libs/log"
)

func main() {

	app := cli.NewApp()

	app.Name = "GoTool"
	// Описание цели программы
	app.Usage = "To save the world"
	// номер версии программы

	app.Version = "1.0.0"

	var host string
	var port int
	//var hostS string

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "host",          // имя параметра
			Value:       "0.0.0.0",       // Значение параметра по умолчанию
			Usage:       "Адрес сервера", // описание функции параметра
			Destination: &host,           // Переменные, которые получают значения
		},
		&cli.IntFlag{
			Name:        "port",
			Value:       8083,
			Usage:       "Server port",
			Destination: &port,
		},
	}

	app.Action = func(context *cli.Context) error {

		log.WithFields(log.Fields{
			"package": "main",
			"func":    "calcLab_server",
			"port":    port,
			"host":    host,
		}).Info("Server to start")

		return server.StartServer(host, port, context.Args().Slice())
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
