package whenPkg

import (
	"context"
	"errors"

	"github.com/theQRL/go-zond"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
)

func (opts *WhenOptions) HandleList() error {
	chain := opts.Globals.Chain

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawNamedBlock], errorChan chan error) {
		results, err := tslib.GetSpecials(chain)
		if err != nil {
			errorChan <- err
			if errors.Is(err, zond.NotFound) {
				return
			}
			cancel()
			return
		}

		for _, result := range results {
			if opts.Globals.Verbose || result.Component == "execution" {
				modelChan <- &result
			}
		}
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
