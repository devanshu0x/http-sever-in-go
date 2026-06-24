package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var SEPARATOR="\r\n"

func parseRequestLine(s string) (*RequestLine,string,error){
	idx:=strings.Index(s,SEPARATOR)
	if idx==-1{
		return nil,"",fmt.Errorf("No strat line found")
	}
	reqLine:=s[:idx]
	remaining:=s[idx+len(SEPARATOR):]

	parts:=strings.Split(reqLine," ")

	return &RequestLine{
		Method: parts[0],
		
	},remaining,nil
}

func RequestFromReader(reader io.Reader) (*Request, error){
	bytes,err:=io.ReadAll(reader)
	if err!=nil{
		return nil,err
	}
	request:=&Request{}
	requestLine,next,err:=parseRequestLine(string(bytes))
	if err!=nil{
		return nil,err
	}
	request.RequestLine=*requestLine
	fmt.Println(next)

	return request,nil

}