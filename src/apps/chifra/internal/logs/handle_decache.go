// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package logsPkg

import (
	"context"

	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/decache"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/walk"
)

func (opts *LogsOptions) HandleDecache() error {
	silent := opts.Globals.TestMode || len(opts.Globals.File) > 0

	itemsToRemove, err := decache.LocationsFromTransactionIds(opts.Conn, opts.TransactionIds)
	if err != nil {
		return err
	}

	ctx := context.Background()
	fetchData := func(modelChan chan types.Modeler[types.RawModeler], errorChan chan error) {
		if msg, err := decache.Decache(opts.Conn, itemsToRemove, silent, walk.Cache_Logs); err != nil {
			errorChan <- err
		} else {
			s := types.SimpleMessage{
				Msg: msg,
			}
			modelChan <- &s
		}
	}

	opts.Globals.NoHeader = true
	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
