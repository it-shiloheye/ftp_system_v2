package server

import (
	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_v2/lib/tls_handler/v2"
	"log"
	"net"
	"net/http"
)

func ServerLoop(ctx ftp_context.Context) (ftp_err error) {
	Srvr := ServerType{
		Port: ServerConfig.PeerPort,
	}
	Srvr.InitServer(certs_loc.tlsCert)
	go Srvr.ServerRun(ctx.Add())
	defer ctx.Finished()
	select {
	case ftp_err = <-Srvr.ErrC:
		break
	case <-ctx.Done():
	}

	return
}
