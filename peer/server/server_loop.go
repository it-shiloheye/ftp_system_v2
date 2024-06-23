package server

import (
	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
)

func TestServerLoop(ctx ftp_context.Context, port string) (ftp_err error) {
	Srvr := ServerType{
		Port: port,
	}
	Srvr.InitServer(C_loc.tlsCert)
	err_c := make(chan error)
	go Srvr.ServerRun(ctx.Add(), err_c)
	defer ctx.Finished()
	select {
	case ftp_err = <-err_c:
		break
	case <-ctx.Done():
	}

	return
}
