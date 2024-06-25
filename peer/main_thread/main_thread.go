package mainthread

import (
	"time"

	ftp_context "github.com/it-shiloheye/ftp_system_v2/lib/context"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging"
	"github.com/it-shiloheye/ftp_system_v2/lib/logging/log_item"
)

var Logger = logging.Logger

func Loop(ctx ftp_context.Context) error {
	loc := log_item.Loc(`Loop(ctx ftp_context.Context) error`)
	defer ctx.Finished()

	tc := time.NewTicker(time.Minute * 5)
	for ok := true; ok; {

		select {
		case <-ctx.Done():
			return Logger.LogErr(loc, ctx.Err())
		case <-tc.C:
		}

	}

	return nil
}
