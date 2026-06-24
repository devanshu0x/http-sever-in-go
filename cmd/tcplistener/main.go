package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	s := make(chan string)

	go func() {
		defer close(s)
		defer f.Close()
		buff := make([]byte, 8)
		var currentLine string
		for {
			length, err := f.Read(buff)
			if err != nil {
				if err == io.EOF {
					if currentLine != "" {
						s <- currentLine
					}
					break
				}
				fmt.Println("Error: ", err)
				break
			}
			parts := strings.Split(string(buff[:length]), "\n")
			for i, part := range parts {
				if i+1 == len(parts) {
					currentLine += part
				} else {
					currentLine += part
					s <- currentLine
					currentLine = ""
				}
			}
		}
	}()

	return s
}

func main() {
	listener,err:=net.Listen("tcp","127.0.0.1:42069")
	if err!=nil{
		log.Fatalf("Failed to start listener: %v",err)
	}
	defer listener.Close()
	fmt.Println("TCP server started at: ",listener.Addr())
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println("Failed to accept connection: ",err)
			continue
		}
		fmt.Println("Connection is accepted")
		ch:= getLinesChannel(conn)

		for val := range ch {
		fmt.Println(val)
		}
	}

	// s := getLinesChannel(file)


	

}
