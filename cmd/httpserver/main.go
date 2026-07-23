package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/devanshu0x/http-sever-in-go/internal/headers"
	"github.com/devanshu0x/http-sever-in-go/internal/request"
	"github.com/devanshu0x/http-sever-in-go/internal/response"
	"github.com/devanshu0x/http-sever-in-go/internal/server"
)

const port = 42069

func main() {

	server, err := server.Serve(port, func(w response.Writer, req *request.Request) {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(400)
			content := []byte(`<html>
  				<head>
    				<title>400 Bad Request</title>
  				</head>
  				<body>
    				<h1>Bad Request</h1>
    				<p>Your request honestly kinda sucked.</p>
  				</body>
				</html>`)
			headers := response.GetDefaultHeaders(len(content))
			headers.Set("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody(content)
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			w.WriteStatusLine(500)
			content := []byte(`<html>
  			<head>
    			<title>500 Internal Server Error</title>
  			</head>
  			<body>
    			<h1>Internal Server Error</h1>
    			<p>Okay, you know what? This one is on me.</p>
  			</body>
			</html>`)
			headers := response.GetDefaultHeaders(len(content))
			headers.Set("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody(content)

		}else if req.RequestLine.RequestTarget=="/video"{
			video,err:=os.ReadFile("assets/vim.mp4")
			if err!=nil{
				fmt.Println("Eror: ",err)
				return
			}
			h:=response.GetDefaultHeaders(len(video))
			h.Set("Content-Type","video/mp4")
			w.WriteStatusLine(200)
			w.WriteHeaders(h)
			w.WriteBody(video)
		}else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {

			path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

			resp, err := http.Get("https://httpbin.org/" + path)
			if err != nil {
				fmt.Printf("Error occured: %s\n", err)
			}

			w.WriteStatusLine(200)
			h:=response.GetDefaultHeaders(0)
			h.Remove("Content-length")
			h.Set("transfer-encoding","chunked")
			h.Add("trailer","x-content-sha256")
			h.Add("trailer","x-content-length")
			w.WriteHeaders(h)
			var msg []byte
			buf := make([]byte, 32)
			for {
				readN, err := resp.Body.Read(buf)
				if err != nil {
					if err == io.EOF {
						w.WriteChunkedBodyDone()
						break
					}
					fmt.Println("Error: ", err)
					break
				}
				msg=append(msg,buf[:readN]...)
				w.WriteChunkedBody(buf[:readN])
			}

			sum:=sha256.Sum256(msg)
			hex:=hex.EncodeToString(sum[:])
			trailer:= headers.NewHeaders()
			trailer.Add("x-content-sha256",hex)
			trailer.Add("x-content-length",fmt.Sprintf("%d",len(msg)))
			w.WriteTrailers(trailer)


		} else {
			w.WriteStatusLine(200)
			content := []byte(`<html>
  				<head>
    				<title>200 OK</title>
  				</head>
  				<body>
    				<h1>Success!</h1>
    				<p>Your request was an absolute banger.</p>
  				</body>
			</html>`)
			headers := response.GetDefaultHeaders(len(content))
			headers.Set("Content-Type", "text/html")
			w.WriteHeaders(headers)
			w.WriteBody(content)
		}

	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
