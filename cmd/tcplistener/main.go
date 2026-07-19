package main

import (
	"fmt"
	"log"
	"net"
	"github.com/devanshu0x/http-sever-in-go/internal/request"
)



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
		
		req,err:=request.RequestFromReader(conn)
		if err!=nil{
			fmt.Printf("Some error occured: %v",err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n",req.RequestLine.Method)
		fmt.Printf("- Target: %s\n",req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n",req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key,val:=range req.Headers{
			fmt.Printf("- %s: %s\n",key,val)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))
	}

	// s := getLinesChannel(file)



	

}
