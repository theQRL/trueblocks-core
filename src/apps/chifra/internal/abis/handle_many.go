package abisPkg

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/abi"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/articulate"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
)

func (opts *AbisOptions) HandleMany() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawModeler], errorChan chan error) {
		for _, addr := range opts.Addrs {
			abiCache := articulate.NewAbiCache(opts.Conn, opts.Known)
			address := base.HexToAddress(addr)
			if len(opts.ProxyFor) > 0 {
				address = base.HexToAddress(opts.ProxyFor)
			}
			err = abi.LoadAbi(opts.Conn, address, &abiCache.AbiMap)
			if err != nil {
				if errors.Is(err, rpc.ErrNotAContract) {
					msg := fmt.Errorf("address %s is not a smart contract", address.Hex())
					errorChan <- msg
					// Report but don't quit processing
				} else {
					// Cancel on all other errors
					errorChan <- err
					cancel()
				}
				// } else if len(opts.ProxyFor) > 0 {
				// TODO: We need to copy the proxied-to ABI to the proxy (replacing)
			}

			abi := simpleAbi{}
			abi.Address = address
			names := abiCache.AbiMap.Keys()
			sort.Strings(names)
			for _, name := range names {
				function := abiCache.AbiMap.GetValue(name)
				if function != nil {
					abi.Functions = append(abi.Functions, *function)
				}
			}
			modelChan <- &abi
		}
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
