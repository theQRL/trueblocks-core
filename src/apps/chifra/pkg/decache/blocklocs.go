package decache

import (
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/identifiers"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/utils"
)

func LocationsFromBlockIds(conn *rpc.Connection, ids []identifiers.Identifier, logs, trace bool) ([]cache.Locator, error) {
	locations := make([]cache.Locator, 0)
	for _, br := range ids {
		blockNums, err := br.ResolveBlocks(conn.Chain)
		if err != nil {
			return nil, err
		}
		for _, bn := range blockNums {
			if logs {
				logGroup := &types.SimpleLogGroup{
					BlockNumber:      bn,
					TransactionIndex: utils.NOPOS,
				}
				locations = append(locations, logGroup)

			} else if trace {
				traceGroup := &types.SimpleTraceGroup{
					BlockNumber:      bn,
					TransactionIndex: utils.NOPOS,
				}
				locations = append(locations, traceGroup)

			} else {
				rawBlock, err := conn.GetBlockHeaderByNumber(bn)
				if err != nil {
					return nil, err
				}
				locations = append(locations, &types.SimpleBlock[string]{
					BlockNumber: bn,
				})
				for index := range rawBlock.Transactions {
					txToRemove := &types.SimpleTransaction{
						BlockNumber:      bn,
						TransactionIndex: uint64(index),
					}
					locations = append(locations, txToRemove)
					locations = append(locations, &types.SimpleTraceGroup{
						BlockNumber:      bn,
						TransactionIndex: base.Txnum(index),
					})
				}
			}
		}
	}
	return locations, nil
}
