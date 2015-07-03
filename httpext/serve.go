package ext

import (
 	"net"
 	"net/http"
 	"strings"
 	"github.com/h2object/h2object/log"
)

func Serve(listener net.Listener, handler http.Handler, proto_name string, logger log.Logger) {
	logger.Info("%s: listening on %s", proto_name, listener.Addr().String())

	server := &http.Server{
		Handler: handler,
	}
	err := server.Serve(listener)
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		logger.Error("ERROR: http.Serve() - %s", err.Error())
	}

	logger.Info("%s: closing %s", proto_name, listener.Addr().String())
}
