package main

import (
	api "github.com/gitKashish/EcomServer/cmd/api"
)

func main() {
	server := api.NewAPIServer(":8080", nil)
	server.Run()
}
