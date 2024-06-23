package server

import (
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
	ftp_tlshandler "github.com/it-shiloheye/ftp_system_v2/lib/tls_handler/v2"
)

var Logger = logging.Logger

type ServerType struct {
	Port string `json:"port"`

	HttpsServer *http.Server
	*gin.Engine

	ErrC      chan error
	ServerRun func(ctx ftp_context.Context)
}

func (st *ServerType) InitServer(server_cert *ftp_tlshandler.TLSCert) <-chan error {

	loc := log_item.Loc("gin_server_main_thread(ctx ftp_context.Context) (err log_item.LogErr)")

	st.Engine = gin.Default()
	r := st.Engine
	st.ErrC = make(chan error, 1)

	server_ip, ip_net, err1 := net.ParseCIDR("192.168.0.0/24")
	if err1 != nil {
		log.Fatalln(&log_item.LogItem{
			Location:  loc,
			After:     `ip, ip_net, err1  := net.ParseCIDR("192.168.0.0/24")`,
			Message:   err1.Error(),
			Level:     log_item.LogLevelError02,
			CallStack: []error{err1},
		})
	}

	r.Use(func(ctx *gin.Context) {
		req_ip := ctx.RemoteIP()

		if valid_ip(req_ip, ip_net, server_ip) {

			ctx.Next()
			return
		}
		ctx.Status(400)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	st.ServerRun = func(ctx ftp_context.Context) {
		defer ctx.Finished()
		log.Println("Starting: gin_server_main_thread")

		port := st.Port
		st.HttpsServer = &http.Server{
			Addr:      st.Port,
			Handler:   r,
			TLSConfig: ftp_tlshandler.ServerTLSConf(server_cert.TlsCert),
		}

		log.Println("\nhttps://127.0.0.1"+port, "\nhttps://"+ServerConfig.LocalIp().String()+port)

		if err_ := st.HttpsServer.ListenAndServeTLS("", ""); err_ != nil {
			st.ErrC <- Logger.LogErr(loc, &log_item.LogItem{
				Message:   "server failed",
				Level:     log_item.LogLevelError01,
				CallStack: []error{err_},
			})
		}
		close(st.ErrC)
	}

	return st.ErrC
}

func valid_ip(ip string, ip_net *net.IPNet, server_ip net.IP) bool {
	req_ip := net.ParseIP(ip)

	if req_ip.Equal(server_ip) || net.IPv6loopback.Equal(req_ip) {
		return true
	}

	if req_ip.IsLoopback() {
		return true
	}

	if ip_net.Contains(req_ip) {
		return true
	}

	return false
}
