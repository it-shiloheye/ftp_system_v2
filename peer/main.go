package main

import (
	"log"
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	db "github.com/it-shiloheye/ftp_system_v2/lib/db_access"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
	browserserver "github.com/it-shiloheye/ftp_system_v2/peer/browser-server"
	server_config "github.com/it-shiloheye/ftp_system_v2/peer/config"
	mainthread "github.com/it-shiloheye/ftp_system_v2/peer/main_thread"
	db_helpers "github.com/it-shiloheye/ftp_system_v2/peer/main_thread/db_access"
	networkpeer "github.com/it-shiloheye/ftp_system_v2/peer/network-peer"
)

var ServerConfig = server_config.ServerConfig
var Logger = logging.Logger

func init() {
	logging.InitialiseLogging(ServerConfig.StorageDirectory)

}

func main() {
	loc := log_item.Loc(`func main()`)
	ctx := ftp_context.CreateNewContext()
	log.Println("hello world from server")

	db.ConnectToDB(ctx)

	go logging.Logger.Engine(ctx, ServerConfig.StorageDirectory)
	for i := 0; ; i += 1 {

		err1 := db_helpers.ConnectClient(ctx)
		if err1 == nil {
			break
		}
		Logger.LogErr(loc, err1)
		<-time.After(time.Second)
		if i > 2 {
			log.Fatalln("fatal termination: couldn't connect to db and get peer_id")
		}
	}
	PeerSrv := networkpeer.CreatePeerServer(ctx)
	BrowserSrv := browserserver.CreateBrowserServer()
	err_c := make(chan error)
	go PeerSrv.ServerRun(ctx.Add(), err_c)
	go BrowserSrv.ServerRun(ctx.Add(), err_c)

	log.Println("\nBrowser: http://127.0.0.1"+ServerConfig.BrowserPort, "\nPeer: https://"+ServerConfig.LocalIp().String()+ServerConfig.PeerPort)

	err1 := mainthread.Loop(ctx.NewChild())
	if err1 != nil {
		log.Fatalln(logging.Logger.LogErr(loc, err1))
	}
	ctx.Wait()
}
