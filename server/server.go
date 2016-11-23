package server

import (
	"io"
	"net"
	"strings"
	"time"
)

type Server struct {
	listener net.Listener
	codec    Codec
	handler  Handler
}

type Codec interface {
	Read(rw io.ReadWriter) (interface{}, error)
	Write(rw io.ReadWriter, v interface{}) error
}

type Handler interface {
	HandleConn(net.Conn)
}

type HandlerFunc func(net.Conn)

func (hf HandlerFunc) HandleConn(conn net.Conn) {
	hf(conn)
}

func newServer(listener net.Listener, codec Codec, handler Handler) *Server {
	return &Server{
		listener: listener,
		codec:    codec,
		handler:  handler,
	}
}

func (server *Server) Listener() net.Listener {
	return server.listener
}

func (server *Server) Serve() error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			return err
		}

		go func() {
			server.handler.HandleConn(conn)
		}()
	}
}

func (server *Server) Close() {
	server.listener.Close()
}

func Dial(network string, address string, codec Codec) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DialTimeout(network string, address string, timeout time.Duration, codec Codec) (net.Conn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}
		return conn, nil
	}
}

func Listen(network string, address string, codec Codec, handler Handler) (*Server, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return newServer(listener, codec, handler), nil
}
