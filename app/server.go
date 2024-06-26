package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/http"
)

var directory = "/tmp"

func main() {
	if len(os.Args) > 2 {
		directory = os.Args[2]
		log.Println("Serving files from", directory)
	}

	server := http.NewServer("0.0.0.0:4221")

	server.HandleStrict(http.GET, "/", func(request http.Request) http.Response {
		return http.OKResponse()
	})

	server.Handle(http.GET, "/echo/", func(request http.Request) http.Response {
		echo := request.Path[6:]
		response := http.NewResponse()
		for header, value := range request.Headers {
			if strings.ToLower(header) == "accept-encoding" {
				encodings := strings.Split(strings.TrimSpace(value), ", ")
				for _, encoding := range encodings {
					if encoding == "gzip" {
						var buf bytes.Buffer
						w := gzip.NewWriter(&buf)
						w.Write([]byte(echo))
						w.Close()
						response.AddHeader("Content-Encoding", "gzip").
							AddHeader("Content-Type", "text/plain").
							SetBody(buf.String()).
							SetStatus(http.OK)
						return response
					}
				}
			}
		}
		response.AddHeader("Content-Type", "text/plain")
		response.SetBody(echo)
		return response
	})

	server.Handle(http.GET, "/user-agent", func(request http.Request) http.Response {
		response := http.OKResponse()
		response.SetBody(request.Headers["User-Agent"])
		return response
	})

	server.Handle(http.GET, "/files/", func(request http.Request) http.Response {
		response := http.NewResponse()
		fileName := request.Path[6:]
		fullPath := filepath.Join(directory, fileName)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			response.AddHeader("Content-Type", "text/plain")
			response.Status = http.NotFound
			return response
		}
		response.AddHeader("Content-Type", "application/octet-stream")
		response.SetBodyFile(fullPath)
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

		return http.CreatedResponse()
	})

	server.ListenAndServe()
}
