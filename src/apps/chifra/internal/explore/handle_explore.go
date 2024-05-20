package explorePkg

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/utils"
)

func (opts *ExploreOptions) HandleExplore() error {
	for _, url := range urls {
		ret := url.getUrl(opts)
		if !opts.Globals.TestMode {
			logger.Info("Opening", ret)
			utils.OpenBrowser(ret)
		} else {
			logger.Info("Not opening", ret, "in test mode")
		}
	}

	return nil
}
