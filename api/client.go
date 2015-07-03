package api

import (
	"fmt"
	"sync"
	"github.com/h2object/rpc"
)

type Auth interface{
	rpc.PreRequest
}

type Logger interface{
	rpc.Logger
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{}) 
	Error(format string, args ...interface{}) 
	Critical(format string, args ...interface{})
}

var UserAgent = "Golang h2object/api package"

type Client struct{
	sync.RWMutex
	addr string
	conn *rpc.Client
}

func NewClient(host string, port int) *Client {
	connection := rpc.NewClient(rpc.H2OAnalyser{})
	clt := &Client{
		addr: fmt.Sprintf("%s:%d", host, port),	
		conn: connection,
	}
	return clt
}