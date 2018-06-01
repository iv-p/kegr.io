package main

import (
	"github.com/gin-gonic/gin"
	"kegr.io/rest_api/controllers"
	"kegr.io/storage_client"
)

func main() {
	client := storage_client.NewClient("localhost:24471")

	r := gin.Default()

	kegController := controllers.NewKegController(client)
	kegController.Register(r)

	liquidController := controllers.NewLiquidController(client)
	liquidController.Register(r)

	r.Run()
}
