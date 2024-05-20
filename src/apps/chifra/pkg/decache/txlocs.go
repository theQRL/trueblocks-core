package decache

import (
	"errors"

	"github.com/theQRL/go-zond"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/identifiers"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
)

func LocationsFromTransactionIds(conn *rpc.Connection, ids []identifiers.Identifier) ([]cache.Locator, error) {
	locations := make([]cache.Locator, 0)
	for _, rng := range ids {
		txIds, err := rng.ResolveTxs(conn.Chain)
		if err != nil && !errors.Is(err, zond.NotFound) {
			continue
		}
		for _, app := range txIds {
			tx := &types.SimpleTransaction{
				BlockNumber:      uint64(app.BlockNumber),
				TransactionIndex: uint64(app.TransactionIndex),
			}
			locations = append(locations, tx)
			locations = append(locations, &types.SimpleTraceGroup{
				BlockNumber:      tx.BlockNumber,
				TransactionIndex: tx.TransactionIndex,
			})
		}
	}
	return locations, nil
}
