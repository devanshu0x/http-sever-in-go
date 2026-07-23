# HTTP From TCP

A minimal HTTP/1.1 server built directly on top of raw TCP sockets in Go without using `net/http`.

The goal of this project is to understand how HTTP works under the hood by implementing the protocol manually, from parsing requests and headers to writing responses over a TCP connection.

## Features

* HTTP/1.1 request parsing
* Response writer
* Routing
* Header parsing
* `Content-Length` support
* Chunked transfer encoding
* Trailer headers

## Run

```bash
go run ./cmd/server
```

## Why?

This project was built for learning low-level networking, TCP sockets, and the HTTP protocol by removing the abstractions provided by high-level frameworks.
