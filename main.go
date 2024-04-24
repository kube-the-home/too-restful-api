package main

import (
	"kube-the-home/too-restful-api/config"
	"kube-the-home/too-restful-api/webserver"
)

func main() {
	config.InitLogger()
	config.InitConfig()

	webserver.Execute()
}
