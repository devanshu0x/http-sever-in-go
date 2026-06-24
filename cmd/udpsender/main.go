package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main(){
	addr,err:=net.ResolveUDPAddr("udp","localhost:42069")
	if err!=nil{
		log.Fatal("Error: ",err)
	}
	conn,err:=net.DialUDP("udp",nil,addr)
	if err!=nil{
		log.Fatal("Error: ",err)
	}

	defer conn.Close()

	reader:=bufio.NewReader(os.Stdin)

	for{
		fmt.Print(">")
		line,err:=reader.ReadString('\n')
		if err!=nil{
		fmt.Println("Error: ",err)
		}
		_,err=conn.Write([]byte(line))
		if err!=nil{
		fmt.Println("Error: ",err)
		}
	}

}