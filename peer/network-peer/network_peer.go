package networkpeer

import ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"

func MainLoop(ctx ftp_context.Context) {
	defer ctx.Finished()
}
