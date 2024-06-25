package main

import (
	"log"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	db "github.com/it-shiloheye/ftp_system_v2/lib/db_access"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	browserserver "github.com/it-shiloheye/ftp_system_v2/peer/browser-server"
	server_config "github.com/it-shiloheye/ftp_system_v2/peer/config"
	networkpeer "github.com/it-shiloheye/ftp_system_v2/peer/network-peer"
)

var ServerConfig = server_config.ServerConfig

func init() {
	logging.InitialiseLogging(ServerConfig.StorageDirectory)

}

func main() {

	ctx := ftp_context.CreateNewContext()
	log.Println("hello world from server")

	db.ConnectToDB(ctx)

	go logging.Logger.Engine(ctx, ServerConfig.StorageDirectory)

	PeerSrv := networkpeer.CreatePeerServer(ctx)
	BrowserSrv := browserserver.CreateBrowserServer()
	err_c := make(chan error)
	go PeerSrv.ServerRun(ctx.Add(), err_c)
	go BrowserSrv.ServerRun(ctx.Add(), err_c)

	log.Println("\nBrowser: http://127.0.0.1"+ServerConfig.BrowserPort, "\nPeer: https://"+ServerConfig.LocalIp().String()+ServerConfig.PeerPort)
	ctx.Wait()
}
