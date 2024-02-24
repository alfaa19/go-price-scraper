package main

import (
	"github.com/alfaa19/go-scraper/handler"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	//Route
	r.GET("/pricelist", handler.GetItemHandler)

	r.Run(":8080")
}
