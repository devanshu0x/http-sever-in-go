package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/devanshu0x/http-sever-in-go/internal/request"
	"github.com/devanshu0x/http-sever-in-go/internal/response"
)

type Server struct{
	listener net.Listener
	isClosed atomic.Bool
	handler Handler
}


type Handler func(w response.Writer, req *request.Request) 

func Serve(port int,handler Handler) (*Server,error){
	listener,err:= net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err!=nil{
		return nil,err
	}

	server:= &Server{
		listener: listener,
		handler: handler,
	}

	go server.Listen()

	return server,nil
}

func (s* Server) Close() error{
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s* Server) Listen(){
	for {
		conn,err:=s.listener.Accept()
		if err!=nil{
			if s.isClosed.Load(){
				return
			}

			continue
		}
		go s.handle(conn)
	}

}

func (s* Server) handle(conn net.Conn){
	defer conn.Close()

	req,err:=request.RequestFromReader(conn)
	if err!=nil{
		fmt.Printf("Error: %v",err)
		return 
	}

	rw:=response.NewResponseWriter(conn)

	s.handler(*rw,req)



	// b:=&bytes.Buffer{}

	// handlerErr:=s.handler(b,req)

	// if(handlerErr!=nil){
	// 	err=response.WriteStatusLine(conn,handlerErr.StatusCode)
	// if err!=nil{
	// 	fmt.Printf("ErrorW: %v",err)
	// }
	// err=response.WriteHeaders(conn,response.GetDefaultHeaders(len(handlerErr.Message)))
	// if err!=nil{
	// 	fmt.Printf("ErrorW: %v",err)
	// }
	// conn.Write([]byte(handlerErr.Message))
	// return
	// }

	// // conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n"))
	// err=response.WriteStatusLine(conn,response.StatusOK)
	// if err!=nil{
	// 	fmt.Printf("ErrorW: %v",err)
	// }
	// err=response.WriteHeaders(conn,response.GetDefaultHeaders(b.Len()))
	// if err!=nil{
	// 	fmt.Printf("ErrorW: %v",err)
	// }

	// conn.Write(b.Bytes())

}