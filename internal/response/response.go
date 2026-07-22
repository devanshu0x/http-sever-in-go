package response

import (
	"fmt"
	"io"

	"github.com/devanshu0x/http-sever-in-go/internal/headers"
)

type StatusCode int

const(
	StatusOK StatusCode=200
	StatusBadRequest StatusCode=400
	StatusInternalServerError StatusCode=500
	WriterStateInit WriterState="init"
	WriterStateStatusLine WriterState="status_line_done"
	WriterStateHeaders WriterState="headers_done"
)

func getReasonPhrase(s StatusCode) string{
	switch s{
	case StatusOK:
		return "OK"
	case StatusBadRequest:
		return "Bad Request"
	case StatusInternalServerError:
		return "Internal Server Error"
	default:
		return ""			
	}
}

type WriterState string

type Writer struct{
	WriterState WriterState
	res io.Writer
}

func NewResponseWriter(w io.Writer) *Writer{
	return &Writer{
		res: w,
		WriterState: WriterStateInit,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error{
	if w.WriterState!=WriterStateInit{
		return fmt.Errorf("In wrong state cannot write status line in this order")
	}
	_,err:=fmt.Fprintf(w.res,"HTTP/1.1 %d %s\r\n",statusCode,getReasonPhrase(statusCode))
	if err!=nil{
		return err
	}

	w.WriterState=WriterStateStatusLine
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error{
	if w.WriterState!=WriterStateStatusLine{
		return fmt.Errorf("In wrong state cannot write headers in this order")
	}

	for key,val:=range headers{
		_,err:=fmt.Fprintf(w.res,"%s: %s\r\n",key,val)
		if err!=nil{
			return err
		}
	}

	w.res.Write([]byte("\r\n"))
	w.WriterState=WriterStateHeaders
	return nil

}

func (w *Writer) WriteBody(p []byte) (int ,error){
	if w.WriterState!=WriterStateHeaders{
		return 0,fmt.Errorf("In wrong state cannot write body in this order")
	}
	return w.res.Write(p)
}


func GetDefaultHeaders(contentLen int) headers.Headers{
	headers:=headers.NewHeaders()
	headers["Content-Length"]=fmt.Sprintf("%d",contentLen)
	headers["Connection"]="close"
	// headers["Content-Type"]="text/plain"
	headers.Set("Content-Type","text/plain")

	return headers
}

