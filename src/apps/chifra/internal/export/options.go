// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * This file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */

package exportPkg

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/theQRL/trueblocks-core/src/apps/chifra/internal/globals"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/caps"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/validate"
)

// ExportOptions provides all command options for the chifra export command.
type ExportOptions struct {
	Addrs       []string              `json:"addrs,omitempty"`       // One or more addresses (0x...) to export
	Topics      []string              `json:"topics,omitempty"`      // Filter by one or more log topics (only for --logs option)
	Fourbytes   []string              `json:"fourbytes,omitempty"`   // Filter by one or more fourbytes (only for transactions and trace options)
	Appearances bool                  `json:"appearances,omitempty"` // Export a list of appearances
	Receipts    bool                  `json:"receipts,omitempty"`    // Export receipts instead of transactional data
	Logs        bool                  `json:"logs,omitempty"`        // Export logs instead of transactional data
	Traces      bool                  `json:"traces,omitempty"`      // Export traces instead of transactional data
	Neighbors   bool                  `json:"neighbors,omitempty"`   // Export the neighbors of the given address
	Accounting  bool                  `json:"accounting,omitempty"`  // Attach accounting records to the exported data (applies to transactions export only)
	Statements  bool                  `json:"statements,omitempty"`  // For the accounting options only, export only statements
	Balances    bool                  `json:"balances,omitempty"`    // Traverse the transaction history and show each change in ETH balances
	Withdrawals bool                  `json:"withdrawals,omitempty"` // Export withdrawals for the given address
	Articulate  bool                  `json:"articulate,omitempty"`  // Articulate transactions, traces, logs, and outputs
	CacheTraces bool                  `json:"cacheTraces,omitempty"` // Force the transaction's traces into the cache
	Count       bool                  `json:"count,omitempty"`       // For --appearances mode only, display only the count of records
	FirstRecord uint64                `json:"firstRecord,omitempty"` // The first record to process
	MaxRecords  uint64                `json:"maxRecords,omitempty"`  // The maximum number of records to process
	Relevant    bool                  `json:"relevant,omitempty"`    // For log and accounting export only, export only logs relevant to one of the given export addresses
	Emitter     []string              `json:"emitter,omitempty"`     // For the --logs option only, filter logs to show only those logs emitted by the given address(es)
	Topic       []string              `json:"topic,omitempty"`       // For the --logs option only, filter logs to show only those with this topic(s)
	Reverted    bool                  `json:"reverted,omitempty"`    // Export only transactions that were reverted
	Asset       []string              `json:"asset,omitempty"`       // For the accounting options only, export statements only for this asset
	Flow        string                `json:"flow,omitempty"`        // For the accounting options only, export statements with incoming, outgoing, or zero value
	Factory     bool                  `json:"factory,omitempty"`     // For --traces only, report addresses created by (or self-destructed by) the given address(es)
	Unripe      bool                  `json:"unripe,omitempty"`      // Export transactions labeled upripe (i.e. less than 28 blocks old)
	Load        string                `json:"load,omitempty"`        // A comma separated list of dynamic traversers to load
	Reversed    bool                  `json:"reversed,omitempty"`    // Produce results in reverse chronological order
	NoZero      bool                  `json:"noZero,omitempty"`      // For the --count option only, suppress the display of zero appearance accounts
	FirstBlock  uint64                `json:"firstBlock,omitempty"`  // First block to process (inclusive)
	LastBlock   uint64                `json:"lastBlock,omitempty"`   // Last block to process (inclusive)
	Globals     globals.GlobalOptions `json:"globals,omitempty"`     // The global options
	Conn        *rpc.Connection       `json:"conn,omitempty"`        // The connection to the RPC server
	BadFlag     error                 `json:"badFlag,omitempty"`     // An error flag if needed
	// EXISTING_CODE
	// EXISTING_CODE
}

var defaultExportOptions = ExportOptions{
	MaxRecords: 250,
	LastBlock:  utils.NOPOS,
}

// testLog is used only during testing to export the options for this test case.
func (opts *ExportOptions) testLog() {
	logger.TestLog(len(opts.Addrs) > 0, "Addrs: ", opts.Addrs)
	logger.TestLog(len(opts.Topics) > 0, "Topics: ", opts.Topics)
	logger.TestLog(len(opts.Fourbytes) > 0, "Fourbytes: ", opts.Fourbytes)
	logger.TestLog(opts.Appearances, "Appearances: ", opts.Appearances)
	logger.TestLog(opts.Receipts, "Receipts: ", opts.Receipts)
	logger.TestLog(opts.Logs, "Logs: ", opts.Logs)
	logger.TestLog(opts.Traces, "Traces: ", opts.Traces)
	logger.TestLog(opts.Neighbors, "Neighbors: ", opts.Neighbors)
	logger.TestLog(opts.Accounting, "Accounting: ", opts.Accounting)
	logger.TestLog(opts.Statements, "Statements: ", opts.Statements)
	logger.TestLog(opts.Balances, "Balances: ", opts.Balances)
	logger.TestLog(opts.Withdrawals, "Withdrawals: ", opts.Withdrawals)
	logger.TestLog(opts.Articulate, "Articulate: ", opts.Articulate)
	logger.TestLog(opts.CacheTraces, "CacheTraces: ", opts.CacheTraces)
	logger.TestLog(opts.Count, "Count: ", opts.Count)
	logger.TestLog(opts.FirstRecord != 0, "FirstRecord: ", opts.FirstRecord)
	logger.TestLog(opts.MaxRecords != 250, "MaxRecords: ", opts.MaxRecords)
	logger.TestLog(opts.Relevant, "Relevant: ", opts.Relevant)
	logger.TestLog(len(opts.Emitter) > 0, "Emitter: ", opts.Emitter)
	logger.TestLog(len(opts.Topic) > 0, "Topic: ", opts.Topic)
	logger.TestLog(opts.Reverted, "Reverted: ", opts.Reverted)
	logger.TestLog(len(opts.Asset) > 0, "Asset: ", opts.Asset)
	logger.TestLog(len(opts.Flow) > 0, "Flow: ", opts.Flow)
	logger.TestLog(opts.Factory, "Factory: ", opts.Factory)
	logger.TestLog(opts.Unripe, "Unripe: ", opts.Unripe)
	logger.TestLog(len(opts.Load) > 0, "Load: ", opts.Load)
	logger.TestLog(opts.Reversed, "Reversed: ", opts.Reversed)
	logger.TestLog(opts.NoZero, "NoZero: ", opts.NoZero)
	logger.TestLog(opts.FirstBlock != 0, "FirstBlock: ", opts.FirstBlock)
	logger.TestLog(opts.LastBlock != 0 && opts.LastBlock != utils.NOPOS, "LastBlock: ", opts.LastBlock)
	opts.Conn.TestLog(opts.getCaches())
	opts.Globals.TestLog()
}

// String implements the Stringer interface
func (opts *ExportOptions) String() string {
	b, _ := json.MarshalIndent(opts, "", "  ")
	return string(b)
}

// exportFinishParseApi finishes the parsing for server invocations. Returns a new ExportOptions.
func exportFinishParseApi(w http.ResponseWriter, r *http.Request) *ExportOptions {
	copy := defaultExportOptions
	opts := &copy
	opts.FirstRecord = 0
	opts.MaxRecords = 250
	opts.FirstBlock = 0
	opts.LastBlock = utils.NOPOS
	for key, value := range r.URL.Query() {
		switch key {
		case "addrs":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Addrs = append(opts.Addrs, s...)
			}
		case "topics":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Topics = append(opts.Topics, s...)
			}
		case "fourbytes":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Fourbytes = append(opts.Fourbytes, s...)
			}
		case "appearances":
			opts.Appearances = true
		case "receipts":
			opts.Receipts = true
		case "logs":
			opts.Logs = true
		case "traces":
			opts.Traces = true
		case "neighbors":
			opts.Neighbors = true
		case "accounting":
			opts.Accounting = true
		case "statements":
			opts.Statements = true
		case "balances":
			opts.Balances = true
		case "withdrawals":
			opts.Withdrawals = true
		case "articulate":
			opts.Articulate = true
		case "cacheTraces":
			opts.CacheTraces = true
		case "count":
			opts.Count = true
		case "firstRecord":
			opts.FirstRecord = globals.ToUint64(value[0])
		case "maxRecords":
			opts.MaxRecords = globals.ToUint64(value[0])
		case "relevant":
			opts.Relevant = true
		case "emitter":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Emitter = append(opts.Emitter, s...)
			}
		case "topic":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Topic = append(opts.Topic, s...)
			}
		case "reverted":
			opts.Reverted = true
		case "asset":
			for _, val := range value {
				s := strings.Split(val, " ") // may contain space separated items
				opts.Asset = append(opts.Asset, s...)
			}
		case "flow":
			opts.Flow = value[0]
		case "factory":
			opts.Factory = true
		case "unripe":
			opts.Unripe = true
		case "load":
			opts.Load = value[0]
		case "reversed":
			opts.Reversed = true
		case "noZero":
			opts.NoZero = true
		case "firstBlock":
			opts.FirstBlock = globals.ToUint64(value[0])
		case "lastBlock":
			opts.LastBlock = globals.ToUint64(value[0])
		default:
			if !copy.Globals.Caps.HasKey(key) {
				opts.BadFlag = validate.Usage("Invalid key ({0}) in {1} route.", key, "export")
			}
		}
	}
	opts.Conn = opts.Globals.FinishParseApi(w, r, opts.getCaches())

	// EXISTING_CODE
	if len(opts.Addrs) > 0 {
		addrs := []string{}
		for _, addr := range opts.Addrs {
			if validate.IsValidTopic(addr) {
				opts.Topic = append(opts.Topic, addr)
				opts.Topics = append(opts.Topics, addr)
			} else if validate.IsValidFourByte(addr) {
				opts.Fourbytes = append(opts.Fourbytes, addr)
			} else {
				addrs = append(addrs, addr)
			}
		}
		opts.Addrs = addrs
	}
	// EXISTING_CODE
	opts.Addrs, _ = opts.Conn.GetEnsAddresses(opts.Addrs)
	opts.Emitter, _ = opts.Conn.GetEnsAddresses(opts.Emitter)
	opts.Asset, _ = opts.Conn.GetEnsAddresses(opts.Asset)

	return opts
}

// exportFinishParse finishes the parsing for command line invocations. Returns a new ExportOptions.
func exportFinishParse(args []string) *ExportOptions {
	// remove duplicates from args if any (not needed in api mode because the server does it).
	dedup := map[string]int{}
	if len(args) > 0 {
		tmp := []string{}
		for _, arg := range args {
			if value := dedup[arg]; value == 0 {
				tmp = append(tmp, arg)
			}
			dedup[arg]++
		}
		args = tmp
	}

	defFmt := "txt"
	opts := GetOptions()
	opts.Conn = opts.Globals.FinishParse(args, opts.getCaches())

	// EXISTING_CODE
	for _, arg := range args {
		if validate.IsValidTopic(arg) {
			opts.Topic = append(opts.Topic, arg)
			opts.Topics = append(opts.Topics, arg)
		} else if validate.IsValidFourByte(arg) {
			opts.Fourbytes = append(opts.Fourbytes, arg)
		} else {
			opts.Addrs = append(opts.Addrs, arg)
		}
	}
	// EXISTING_CODE
	opts.Addrs, _ = opts.Conn.GetEnsAddresses(opts.Addrs)
	opts.Emitter, _ = opts.Conn.GetEnsAddresses(opts.Emitter)
	opts.Asset, _ = opts.Conn.GetEnsAddresses(opts.Asset)
	if len(opts.Globals.Format) == 0 || opts.Globals.Format == "none" {
		opts.Globals.Format = defFmt
	}

	return opts
}

func GetOptions() *ExportOptions {
	// EXISTING_CODE
	// EXISTING_CODE
	return &defaultExportOptions
}

func ResetOptions(testMode bool) {
	// We want to keep writer between command file calls
	w := GetOptions().Globals.Writer
	defaultExportOptions = ExportOptions{}
	globals.SetDefaults(&defaultExportOptions.Globals)
	defaultExportOptions.Globals.TestMode = testMode
	defaultExportOptions.Globals.Writer = w
	capabilities := caps.Default // Additional global caps for chifra export
	// EXISTING_CODE
	capabilities = capabilities.Add(caps.Caching)
	capabilities = capabilities.Add(caps.Ether)
	// EXISTING_CODE
	defaultExportOptions.Globals.Caps = capabilities
}

func (opts *ExportOptions) getCaches() (m map[string]bool) {
	// EXISTING_CODE
	m = map[string]bool{
		// TODO: Enabled neighbors cache
		"transactions": true,
		"statements":   opts.Accounting,
		"traces":       opts.CacheTraces || (opts.Globals.Cache && (opts.Traces || opts.Neighbors)),
	}
	// EXISTING_CODE
	return
}

// EXISTING_CODE
// EXISTING_CODE

