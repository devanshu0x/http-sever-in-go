package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

var crlf= []byte("\r\n")

func isTokenChar(c byte) bool {
    switch {
    case 'a' <= c && c <= 'z':
        return true
    case 'A' <= c && c <= 'Z':
        return true
    case '0' <= c && c <= '9':
        return true
    }

    switch c {
    case '!', '#', '$', '%', '&', '\'', '*',
        '+', '-', '.', '^', '_', '`', '|', '~':
        return true
    }

    return false
}

func validHeaderKey(key []byte) bool {
	for _,c:=range key{
		if !isTokenChar(c){
			return false
		}
	}

	return true
}

func (h Headers) Get(key string) string{
	key=strings.ToLower(key)

	val,ok:= h[key]
	if !ok{
		return ""
	}
	return val
}

func (h Headers) Set(key, val string){
	key=strings.ToLower(key)
	h[key]=val
}

func (h Headers) Parse(data []byte) (n int, done bool, err error){
	bytesRead:=0;
	idx:=bytes.Index(data,crlf)
	switch idx{
	case -1:
		return bytesRead,false,nil
	case 0:
		return len(crlf),true,nil
	default:
		bytesRead=idx+len(crlf)
		fieldLine:=data[:idx]
		splitBytes:=bytes.SplitN(fieldLine,[]byte(":"),2)

		if len(splitBytes)!=2{
			return 0,false,fmt.Errorf("Field line in wrong format")
		}
		
		if keyLen:=len(splitBytes[0]); splitBytes[0][keyLen-1]==byte(' ') || splitBytes[0][0]==byte(' '){
			return 0,false, fmt.Errorf("Space present between key and colon")
		}
		key:=bytes.Trim(splitBytes[0]," ")
		value:=bytes.Trim(splitBytes[1]," ")

		if(!validHeaderKey(key)){
			return 0,false,fmt.Errorf("Wrong token in key")
		}

		ParsedKey:=strings.ToLower(string(key))
		
		if v,ok:=h[ParsedKey]; ok{
			v=fmt.Sprintf("%s, %s",v,value)
			h[ParsedKey]=v
		}else{
			h[ParsedKey]=string(value)
		}

		return bytesRead,false,nil
	}
}