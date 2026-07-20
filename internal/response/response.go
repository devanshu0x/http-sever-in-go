package response

import "io"

type StatusCode int

const(
	StatusOK StatusCode=200
	StatusBadRequest StatusCode=400
	StatusInternalServerError StatusCode=500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error{

}