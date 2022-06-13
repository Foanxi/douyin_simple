package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/initialize"
	"github.com/gin-gonic/gin"
)

func main() {
	initialize.LoadConfig()
	initialize.Mysql()
	r := gin.Default()
	initRouter(r)

	err := r.Run()
	if err != nil {
		fmt.Print(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
