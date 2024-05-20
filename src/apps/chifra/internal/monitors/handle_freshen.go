package monitorsPkg

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/monitor"
)

func (opts *MonitorsOptions) FreshenMonitorsForWatch(addrs []base.Address) (bool, error) {
	strs := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		strs = append(strs, addr.Hex())
	}

	unusedMonitors := make([]monitor.Monitor, 0, len(addrs))
	var updater = monitor.NewUpdater(opts.Globals.Chain, opts.Globals.TestMode, true, strs)
	return updater.FreshenMonitors(&unusedMonitors)
}
