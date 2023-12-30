package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

var directory = "/tmp"

func main() {
	if len(os.Args) > 2 {
		directory = os.Args[2]
		fmt.Println("Serving files from", directory)
	}

	server := http.NewServer("0.0.0.0:4221")

	server.HandleStrict(http.GET, "/", func(request http.Request) http.Response {
		response := http.NewResponse()
		response.AddHeader("Content-Type", "text/plain")
		return response
	})

	server.Handle(http.GET, "/echo/", func(request http.Request) http.Response {
		response := http.NewResponse()
		response.AddHeader("Content-Type", "text/plain")
		response.SetBody(request.Path[6:])
		return response
	})

	server.Handle(http.GET, "/user-agent", func(request http.Request) http.Response {
		response := http.NewResponse()
		response.AddHeader("Content-Type", "text/plain")
		response.SetBody(request.Headers["User-Agent"])
		return response
	})

	server.Handle(http.GET, "/files/", func(request http.Request) http.Response {
		response := http.NewResponse()
		filePath := directory + request.Path[6:]
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			response.AddHeader("Content-Type", "text/plain")
			response.Status = http.NotFound
			return response
		}
		response.AddHeader("Content-Type", "application/octet-stream")
		response.SetBodyFile(filePath)
		return response
	})

	server.Handle(http.POST, "/files/", func(request http.Request) http.Response {
		fileName := request.Path[7:]
		fullPath := filepath.Join(directory, fileName)
		file, err := os.Create(fullPath)
		if err != nil {
			response := http.NewResponse()
			fmt.Println("Error creating file:", err)
			response.Status = http.IntError
			return response
		}
		defer file.Close()

		file.Write([]byte(request.Body))

		response := http.NewResponse()
		response.AddHeader("Content-Type", "text/plain")
		response.Status = http.Created
		return response
	})

	server.ListenAndServe()
}
