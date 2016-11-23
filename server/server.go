package server

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

type Server struct {
	network  string
	listener net.Listener
	codec    Codec
	handler  Handler
	udpConn  *net.UDPConn
}

type Codec interface {
	Read(rw io.ReadWriter) (interface{}, error)
	Write(rw io.ReadWriter, v interface{}) error
}

type Handler interface {
	HandleConn(net.Conn, *net.UDPConn, []byte)
}

type HandlerFunc func(net.Conn, *net.UDPConn, []byte)

func (hf HandlerFunc) HandleConn(conn net.Conn, udpConn *net.UDPConn, udpData []byte) {
	hf(conn, udpConn, udpData)
}

func newServer(network string, listener net.Listener, udpConn *net.UDPConn, codec Codec, handler Handler) *Server {
	return &Server{
		network:  network,
		listener: listener,
		udpConn:  udpConn,
		codec:    codec,
		handler:  handler,
	}
}

func (server *Server) Network() string {
	return server.network
}

func (server *Server) Listener() net.Listener {
	return server.listener
}

func (server *Server) UDPConn() *net.UDPConn {
	return server.udpConn
}

func (server *Server) Serve() error {
	switch server.network {
	case "tcp":
		for {
			conn, err := Accept(server.listener)
			if err != nil {
				return err
			}

			go func() {
				server.handler.HandleConn(conn, nil, nil)
			}()
		}
	case "udp":
		for {
			buf := make([]byte, 10240)
			n, _, err := server.udpConn.ReadFromUDP(buf)
			if err != nil {
				return err
			}
			go func() {
				server.handler.HandleConn(nil, server.udpConn, buf[:n])
			}()
		}
	default:
		return errors.New("Network type not supported")
	}
	return nil
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
	switch network {
	case "tcp":
		listener, err := net.Listen(network, address)
		if err != nil {
			return nil, err
		}
		return newServer(network, listener, nil, codec, handler), nil
	case "udp":
		addr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		conn, err := net.ListenUDP(network, addr)
		if err != nil {
			return nil, err
		}
		return newServer(network, nil, conn, codec, handler), nil
	default:
		return nil, errors.New("Network type not supported")
	}
}
