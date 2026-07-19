package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/devanshu0x/http-sever-in-go/internal/headers"
)

type parserState string

const(
	StateInit parserState = "init"
	StateHeadersParsing  parserState="headers_parsing"
	StateBodyParsing parserState="body_parsing"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	state parserState
	Body []byte
}

func newRequest() *Request{
	return &Request{
		state: StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error){
	read:=0
outer:
	for{
		switch r.state{
		case StateInit:
			rl,n,err:=parseRequestLine(data[read:])
			if err!=nil{
				return 0,err
			}
			if n==0{
				break outer
			}
			r.RequestLine=*rl
			read+=n
			r.state=StateHeadersParsing
		case StateHeadersParsing:
			n,done,err:=r.Headers.Parse(data[read:])
			if err!=nil{
				return read,err
			}
			if n==0 && !done{
				break outer
			}
			read+=n
			if done{
				if r.Headers.Get("content-length")!=""{
					r.state=StateBodyParsing
				}else{
					r.state=StateDone
				}
			}	
		case StateBodyParsing:
			r.Body = append(r.Body, data[read:]...)
			read+=len(data[read:])
			contentLen:=r.Headers.Get("content-length")

			parsedLen,err:=strconv.Atoi(contentLen)
			if err!=nil{
				return read,err
			}

			if len(r.Body)>parsedLen{
				return read,fmt.Errorf("Body is longer than content length")
			}else if len(r.Body)==parsedLen{
				r.state=StateDone
			}else{
				break outer
			}

		case StateDone:
			break outer
		}
	}

	return read,nil
}

func (r *Request) done() bool{
	return r.state==StateDone
}


var SEPARATOR=[]byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine,int,error){
	idx:=bytes.Index(b,SEPARATOR)
	if idx==-1{
		return nil,0,nil
	}
	reqLine:=b[:idx]
	read:= idx+len(SEPARATOR)
  
	parts:=bytes.Split(reqLine,[]byte(" "))

	if len(parts)!=3{
		return nil,0,fmt.Errorf("Wrong request line format")
	}

	httpVer:=bytes.Split(parts[2],[]byte("/"))
	if len(httpVer)!=2 || string(httpVer[0])!="HTTP" || string(httpVer[1])!="1.1"{
		return nil,0,fmt.Errorf("Unsupported http version")
	}

	if strings.ToUpper(string(parts[0]))!=string(parts[0]){
		return nil,0,fmt.Errorf("Unsupported method format")
	}

	return &RequestLine{
		Method: string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:string(httpVer[1]) ,
		
	},read,nil
}

func RequestFromReader(reader io.Reader) (*Request, error){
	request:=newRequest()

	buf:= make([]byte,1024)
	// What if bufLen exceeds??? How to take care of it
	bufLen:=0
	for !request.done(){
		n,err:=reader.Read(buf[bufLen:])
		if err!=nil{
			return nil,err
		}
		bufLen+=n

		readN,err:=request.parse(buf[:bufLen])
		if err!=nil{
			return nil,err
		}

		copy(buf,buf[readN:bufLen])
		bufLen-=readN

	}
	
	return request,nil

}