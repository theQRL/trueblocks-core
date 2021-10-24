/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
	"github.com/spf13/cobra"
)

type ExploreType uint8

const (
	ExploreNone ExploreType = 1 << iota
	ExploreAddress
	ExploreName
	ExploreEnsName
	ExploreTx
	ExploreBlock
	ExploreFourByte
)

func (t ExploreType) String() string {
	switch t {
	case ExploreNone:
		return "ExploreNone"
	case ExploreAddress:
		return "ExploreAddress"
	case ExploreName:
		return "ExploreName"
	case ExploreEnsName:
		return "ExploreEnsName"
	case ExploreTx:
		return "ExploreTx"
	case ExploreBlock:
		return "ExploreBlock"
	case ExploreFourByte:
		return "ExploreFourByte"
	default:
		return fmt.Sprintf("%d", t)
	}
}

type ExploreUrl struct {
	term     string
	termType ExploreType
}

var urls []ExploreUrl

func validateExploreArgs(cmd *cobra.Command, args []string) error {
	TestLogExplore(args)

	if ExploreOpts.google && ExploreOpts.local {
		return validate.Usage("Choose either --google or --local, not both.")
	}

	for _, arg := range args {
		arg = strings.ToLower(arg)

		valid, _ := validate.IsValidAddress(arg)
		if valid {
			utils.TestLogBool("is_addr", true)
			if strings.Contains(arg, ".eth") {
				urls = append(urls, ExploreUrl{arg, ExploreEnsName})
			} else {
				urls = append(urls, ExploreUrl{arg, ExploreAddress})
			}
			continue
		}

		if ExploreOpts.google {
			return validate.Usage("Option --google allows only an address term.")
		}

		valid, _ = validate.IsValidTransId([]string{arg}, validate.ValidTransId)
		if valid {
			// TODO: Transactions are block_hash.tx_id, block_num.tx_id or tx_hash
			// TODO: We need to check to see if this argument is a valid on-chain
			// TODO: transaction and if not fail in the first two cases and pass
			// TODO: it on if it's a tx_hash (because it might be a block_hash)
			utils.TestLogBool("is__tx", true)
			urls = append(urls, ExploreUrl{arg, ExploreTx})
			continue
		}

		valid, _ = validate.IsValidBlockId([]string{arg}, validate.ValidBlockId)
		if valid {
			// TODO: The block number needs to be resolved (for example from a hash)
			// TODO: or a special block
			utils.TestLogBool("is_block", true)
			urls = append(urls, ExploreUrl{arg, ExploreBlock})
			continue
		}

		valid, _ = validate.IsValidFourByte(arg)
		if valid {
			utils.TestLogBool("is_fourbyte", true)
			urls = append(urls, ExploreUrl{arg, ExploreFourByte})
			continue
		}

		return validate.Usage("The argument ({0}) does not appear to be valid.", arg)
	}

	if len(urls) == 0 {
		urls = append(urls, ExploreUrl{"", ExploreNone})
	}
	err := validateGlobalFlags(cmd, args)
	if err != nil {
		return err
	}

	return nil
}

func (u *ExploreUrl) getUrl() string {
	if ExploreOpts.google {
		// TODO: How does one do an assertion?
		// assert(u.termType == ExploreAddress)
		return "https://www.google.com/search?q=" + u.term + "+-etherscan+-etherchain+-bloxy+-bitquery+-ethplorer+-tokenview+-anyblocks+-explorer"
	}

	if u.termType == ExploreFourByte {
		return "https://www.4byte.directory/signatures/?bytes4_signature=" + u.term
	}

	if u.termType == ExploreEnsName {
		return "https://etherscan.io/enslookup-search?search=" + u.term
	}

	url := "https://etherscan.io/"
	query := ""
	switch u.termType {
	case ExploreNone:
		return url
	case ExploreTx:
		query = "tx/" + u.term
	case ExploreBlock:
		query = "block/" + u.term
	case ExploreName:
		// TODO: we must resolve the name if possible or fail
		fallthrough
	case ExploreAddress:
		fallthrough
	default:
		query = "address/" + u.term
	}

	if ExploreOpts.local {
		url = "http://localhost:1234/"
		query = strings.Replace(query, "tx/", "explorer/transactions/", -1)
		query = strings.Replace(query, "block/", "explorer/blocks/", -1)
		query = strings.Replace(query, "address/", "dashboard/accounts?address=", -1)
	}

	return url + query
}

func runExplore(cmd *cobra.Command, args []string) {
	for _, url := range urls {
		fmt.Printf("Opening %s\n", url.getUrl())
		if !utils.IsTestMode() {
			utils.OpenBrowser(url.getUrl())
		}
	}
}

// TODO: If isHash determine if it's a tx or a block
// TODO: Turn off OPT_FMT OPT_VERBOSE
// TODO: Read base URLs from config file
