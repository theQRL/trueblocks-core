package scrapePkg

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/notify"
)

func NotifyChunkWritten(chunk index.Chunk, chunkPath string) (err error) {
	var cidString string
	if ok, _ := NotifyConfigured(); ok {
		if cidString, err = index.ChunkCid(chunkPath); err != nil {
			return err
		}
	}
	// Generate range from path, as chunks sometimes don't have Range set
	chunkRange := base.RangeFromFilename(index.ToIndexPath(chunkPath))
	return Notify(notify.Notification[[]notify.NotificationPayloadChunkWritten]{
		Msg:  notify.MessageChunkWritten,
		Meta: nil,
		Payload: []notify.NotificationPayloadChunkWritten{
			{
				Cid:    cidString,
				Range:  chunkRange.String(),
				Author: config.GetRootConfig().Settings.Notify.Author,
			},
		},
	})
}
