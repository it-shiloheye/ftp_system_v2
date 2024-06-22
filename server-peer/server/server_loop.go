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

func ServerLoop(ctx ftp_context.Context) (ftp_err log_item.LogErr) {
	defer ctx.Finished()
	select {
	case ftp_err = <-gin_server_main_thread(ctx, certs_loc.tlsCert):
		break
	case <-ctx.Done():
	}

	return
}

func gin_server_main_thread(ctx ftp_context.Context, server_cert *ftp_tlshandler.TLSCert) <-chan log_item.LogErr {
	loc := log_item.Loc("gin_server_main_thread(ctx ftp_context.Context) (err log_item.LogErr)")

	err_c := make(chan log_item.LogErr, 1)

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

	valid_ip := func(ip string) bool {
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

	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		req_ip := ctx.RemoteIP()

		if valid_ip(req_ip) {

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

	RegisterRoutes(r)

	ctx.Add()
	go func() {
		defer ctx.Finished()
		log.Println("Starting: gin_server_main_thread")

		srv := http.Server{
			Addr:      ServerConfig.ServerPort,
			Handler:   r,
			TLSConfig: ftp_tlshandler.ServerTLSConf(server_cert.TlsCert),
		}

		log.Println("https://127.0.0.1", srv.Addr)

		if err_ := srv.ListenAndServeTLS("", ""); err_ != nil {
			err_c <- log_item.NewLogItem(loc, log_item.LogLevelError01).
				SetAfter(`srv.ListenAndServeTLS("","")`).
				SetMessage("server failed").
				AppendParentError(err_)
		}
		close(err_c)
	}()

	return err_c
}
