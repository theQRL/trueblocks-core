package articulate

import (
	"errors"

	goEthAbi "github.com/theQRL/go-zond/accounts/abi"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/abi"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
)

func (abiCache *AbiCache) ArticulateTrace(trace *types.SimpleTrace) (err error) {
	found, err := articulateTrace(trace, &abiCache.AbiMap)
	if err != nil {
		return err

	} else if found != nil {
		trace.ArticulatedTrace = found
		return nil

	} else {
		address := trace.Action.To
		if !abiCache.loadedMap.GetValue(address) && !abiCache.skipMap.GetValue(address) {
			if err = abi.LoadAbi(abiCache.Conn, address, &abiCache.AbiMap); err != nil {
				abiCache.skipMap.SetValue(address, true)
				if !errors.Is(err, rpc.ErrNotAContract) {
					// Not being a contract is not an error because we want to articulate the input in case it's a message
					return err
				}
			} else {
				abiCache.loadedMap.SetValue(address, true)
			}
		}

		if !abiCache.skipMap.GetValue(address) {
			if trace.ArticulatedTrace, err = articulateTrace(trace, &abiCache.AbiMap); err != nil {
				return err
			}
		}

		return nil
	}
}

func articulateTrace(trace *types.SimpleTrace, abiMap *abi.SelectorSyncMap) (articulated *types.SimpleFunction, err error) {
	input := trace.Action.Input
	if len(input) < 10 {
		return
	}

	encoding := input[:10]
	articulated = abiMap.GetValue(encoding)

	if trace.Result == nil || articulated == nil {
		return
	}

	var abiMethod *goEthAbi.Method

	if len(trace.Action.Input) > 10 {
		abiMethod, err = articulated.GetAbiMethod()
		if err != nil {
			return nil, err
		}
		err = articulateArguments(
			abiMethod.Inputs,
			trace.Action.Input[10:],
			nil,
			articulated.Inputs,
		)
		if err != nil {
			return
		}
	}

	abiMethod, err = articulated.GetAbiMethod()
	if err != nil {
		return nil, err
	}
	err = articulateArguments(
		abiMethod.Outputs,
		trace.Result.Output[2:],
		nil,
		articulated.Outputs,
	)

	return
}
