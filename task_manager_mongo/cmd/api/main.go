package main

import (
	"fmt"
	"task_manager/routes"
)

func main() {
	route := routes.SetupRouter()
	route.Run(":8080")
	fmt.Println("Server running at http://localhost:8080")
}
