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

type parserState string

const(
	StateInit parserState = "init"
	StateDone parserState = "done"
)

var SEPARATOR="\r\n"

func parseRequestLine(s string) (*RequestLine,string,error){
	idx:=strings.Index(s,SEPARATOR)
	if idx==-1{
		return nil,"",fmt.Errorf("No strat line found")
	}
	reqLine:=s[:idx]
	remaining:=s[idx+len(SEPARATOR):]
  
	parts:=strings.Split(reqLine," ")

	if len(parts)!=3{
		return nil,"",fmt.Errorf("Wrong request line format")
	}

	httpVer:=strings.Split(parts[2],"/")
	if len(httpVer)!=2 || httpVer[0]!="HTTP" || httpVer[1]!="1.1"{
		return nil,"",fmt.Errorf("Unsupported http version")
	}

	if strings.ToUpper(parts[0])!=parts[0]{
		return nil,"",fmt.Errorf("Unsupported method format")
	}

	return &RequestLine{
		Method: parts[0],
		RequestTarget: parts[1],
		HttpVersion:httpVer[1] ,
		
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