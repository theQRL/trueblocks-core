// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package initPkg

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/history"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/validate"
)

func (opts *InitOptions) validateInit() error {
	chain := opts.Globals.Chain

	opts.testLog()

	if opts.BadFlag != nil {
		return opts.BadFlag
	}

	if !config.IsChainConfigured(chain) {
		return validate.Usage("chain {0} is not properly configured.", chain)
	}

	if opts.Globals.TestMode {
		return validate.Usage("integration testing was skipped for chifra init")
	}

	if len(opts.Publisher) > 0 {
		err := validate.ValidateExactlyOneAddr([]string{opts.Publisher})
		if err != nil {
			return err
		}
	}

	historyFile := config.PathToCache(chain) + "tmp/history.txt"
	if history.FromHistoryBool(historyFile, "init") && !opts.All {
		return validate.Usage("You previously called chifra init --all. You must continue to do so.")
	}

	return opts.Globals.Validate()
}
