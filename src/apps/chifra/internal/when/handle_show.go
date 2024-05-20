// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package whenPkg

import (
	"context"
	"errors"

	"github.com/theQRL/go-zond"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/identifiers"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
)

func (opts *WhenOptions) HandleShow() error {
	chain := opts.Globals.Chain

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawNamedBlock], errorChan chan error) {
		for _, br := range opts.BlockIds {
			blockNums, err := br.ResolveBlocks(chain)
			if err != nil {
				errorChan <- err
				if errors.Is(err, zond.NotFound) {
					continue
				}
				cancel()
				return
			}

			for _, bn := range blockNums {
				block, err := opts.Conn.GetBlockHeaderByNumber(bn)
				if err != nil {
					errorChan <- err
					if errors.Is(err, zond.NotFound) {
						continue
					}
					cancel()
					return
				}
				if br.StartType == identifiers.BlockHash && base.HexToHash(br.Orig) != block.Hash {
					errorChan <- errors.New("block hash not found")
					continue
				}

				nb, _ := tslib.FromBnToNamedBlock(chain, block.BlockNumber)
				if nb == nil {
					modelChan <- &types.SimpleNamedBlock{
						BlockNumber: block.BlockNumber,
						Timestamp:   block.Timestamp,
					}
				} else {
					modelChan <- nb
				}
			}
		}
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
