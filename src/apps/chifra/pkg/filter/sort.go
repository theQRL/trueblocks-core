package filter

import (
	"sort"

	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/index"
)

type AppearanceSort int

const (
	NotSorted AppearanceSort = iota
	Sorted
	Reversed
)

func (f *AppearanceFilter) Sort(fromDisc []index.AppearanceRecord) {
	if f.sortBy == Sorted || f.sortBy == Reversed {
		sort.Slice(fromDisc, func(i, j int) bool {
			if f.sortBy == Reversed {
				i, j = j, i
			}
			si := (uint64(fromDisc[i].BlockNumber) << 32) + uint64(fromDisc[i].TransactionIndex)
			sj := (uint64(fromDisc[j].BlockNumber) << 32) + uint64(fromDisc[j].TransactionIndex)
			return si < sj
		})
	}
}
