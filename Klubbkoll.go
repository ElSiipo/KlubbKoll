package main

import (
	"KlubbKoll/app"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := SetupRouter()

	log.Fatal(http.ListenAndServe(":1234", router))
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		v1.GET("/clubs", app.GetClubs)
		// v1.GET("/clubs/:id", app.GetClub)
	}

	return router
}
