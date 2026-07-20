package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct{
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server,error){
	listener,err:= net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err!=nil{
		return nil,err
	}

	server:= &Server{
		listener: listener,
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

	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n"))
}