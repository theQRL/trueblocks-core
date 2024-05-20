package exportPkg

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/monitor"
)

func (opts *ExportOptions) FreshenMonitorsForExport(monitorArray *[]monitor.Monitor) (bool, error) {
	var updater = monitor.NewUpdater(opts.Globals.Chain, opts.Globals.TestMode, true, opts.Addrs)
	return updater.FreshenMonitors(monitorArray)
}
