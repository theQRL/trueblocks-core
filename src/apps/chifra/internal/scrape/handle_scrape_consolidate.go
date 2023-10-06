package scrapePkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

const asciiAppearanceSize = 59

// HandleScrapeConsolidate calls into the block scraper to (a) call Blaze and (b) consolidate if applicable
func (opts *ScrapeOptions) HandleScrapeConsolidate(progressThen *rpc.MetaData, blazeOpts *BlazeOptions) (bool, error) {
	chain := blazeOpts.Chain

	// Get a sorted list of files in the ripe folder
	ripeFolder := filepath.Join(config.PathToIndex(chain), "ripe")
	ripeFileList, err := os.ReadDir(ripeFolder)
	if err != nil {
		return true, err
	}

	stageFolder := filepath.Join(config.PathToIndex(chain), "staging")
	if len(ripeFileList) == 0 {
		// On active chains, this most likely never happens, but on some less used or private chains, this is a frequent occurrence.
		// return a message, but don't do anything about it.
		msg := fmt.Sprintf("No new blocks at block %d (%d away from head)%s", progressThen.Latest, (progressThen.Latest - progressThen.Ripe), spaces)
		logger.Info(msg)

		// we need to move the file to the end of the scraped range so we show progress
		stageFn, _ := file.LatestFileInFolder(stageFolder) // it may not exist...
		stageRange := base.RangeFromFilename(stageFn)
		newRangeLast := utils.Min(blazeOpts.RipeBlock, blazeOpts.StartBlock+opts.BlockCnt-1)
		if stageRange.Last < newRangeLast {
			newRange := base.FileRange{First: stageRange.First, Last: newRangeLast}
			newFilename := filepath.Join(stageFolder, newRange.String()+".txt")
			_ = os.Rename(stageFn, newFilename)
			_ = os.Remove(stageFn) // seems redundant, but may not be on some operating systems
		}
		return true, nil
	}

	// Check to see if we got as many ripe files as we were expecting. In the case when AllowMissing is true, we
	// can't really know, but if AllowMissing is false, then the number of files should be the same as the range width
	ripeCnt := len(ripeFileList)
	if uint64(ripeCnt) < (blazeOpts.BlockCount - blazeOpts.UnripeDist) {
		// Then, if they are not at least sequential, clean up and try again...
		allowMissing := config.GetScrape(chain).AllowMissing
		if err := isListSequential(chain, ripeFileList, allowMissing); err != nil {
			_ = index.CleanTemporaryFolders(config.PathToCache(chain), false)
			return true, err
		}
	}

	stageFn, _ := file.LatestFileInFolder(stageFolder) // it may not exist...
	nAppsThen := int(file.FileSize(stageFn) / asciiAppearanceSize)

	// ripeRange := rangeFromFileList(ripeFileList)
	stageRange := base.RangeFromFilename(stageFn)

	curRange := base.FileRange{First: blazeOpts.StartBlock, Last: blazeOpts.StartBlock + opts.BlockCnt - 1}
	if file.FileExists(stageFn) {
		curRange = stageRange
	}

	// Note, this file may be empty or non-existant
	tmpPath := filepath.Join(config.PathToCache(chain) + "tmp")
	backupFn, err := file.MakeBackup(tmpPath, stageFn)
	if err != nil {
		return true, errors.New("Could not create backup file: " + err.Error())
	}

	// logger.Info("Created backup file for stage")
	defer func() {
		if backupFn != "" && file.FileExists(backupFn) {
			// If the backup file exists, something failed, so we replace the original file.
			// logger.Info("Replacing backed up staging file")
			_ = os.Rename(backupFn, stageFn)
			_ = os.Remove(backupFn) // seems redundant, but may not be on some operating systems
		}
	}()

	appearances := file.AsciiFileToLines(stageFn)
	os.Remove(stageFn) // we have a backup copy, so it's not so bad to delete it here
	for _, ripeFile := range ripeFileList {
		ripePath := filepath.Join(ripeFolder, ripeFile.Name())
		appearances = append(appearances, file.AsciiFileToLines(ripePath)...)
		os.Remove(ripePath) // if we fail halfway through, this will get noticed next time around and cleaned up
		curCount := uint64(len(appearances))

		ripeRange := base.RangeFromFilename(ripePath)
		curRange.Last = ripeRange.Last

		isSnap := (curRange.Last >= opts.Settings.FirstSnap && (curRange.Last%opts.Settings.SnapToGrid) == 0)
		isOvertop := (curCount >= uint64(opts.Settings.AppsPerChunk))

		if isSnap || isOvertop {
			appMap := make(index.AddressAppearanceMap, len(appearances))
			for _, line := range appearances {
				parts := strings.Split(line, "\t")
				if len(parts) == 3 {
					addr := strings.ToLower(parts[0])
					bn, _ := strconv.ParseUint(parts[1], 10, 32)
					txid, _ := strconv.ParseUint(parts[2], 10, 32)
					appMap[addr] = append(appMap[addr], index.AppearanceRecord{
						BlockNumber:   uint32(bn),
						TransactionId: uint32(txid),
					})
				}
			}

			indexPath := config.PathToIndex(chain) + "finalized/" + curRange.String() + ".bin"
			if report, err := index.WriteChunk(chain, opts.PublisherAddr, indexPath, appMap, len(appearances)); err != nil {
				return false, err
			} else if report == nil {
				logger.Fatal("Should not happen, write chunk returned empty report")
			} else {
				report.Snapped = isSnap
				report.Report()
			}

			curRange.First = curRange.Last + 1
			appearances = []string{}
		}
	}

	if len(appearances) > 0 {
		lineLast := appearances[len(appearances)-1]
		parts := strings.Split(lineLast, "\t")
		Last := uint64(0)
		if len(parts) > 1 {
			Last, _ = strconv.ParseUint(parts[1], 10, 32)
			Last = utils.Max(utils.Min(blazeOpts.RipeBlock, blazeOpts.StartBlock+opts.BlockCnt-1), Last)
		} else {
			return true, errors.New("Cannot find last block number at lineLast in consolidate: " + lineLast)
		}

		conn := rpc.TempConnection(chain)
		m, _ := conn.GetMetaData(opts.Globals.TestMode)
		rng := base.FileRange{First: m.Finalized + 1, Last: Last}
		f := fmt.Sprintf("%s.txt", rng)
		fileName := filepath.Join(config.PathToIndex(chain), "staging", f)
		err = file.LinesToAsciiFile(fileName, appearances)
		if err != nil {
			os.Remove(fileName) // cleans up by replacing the previous stage
			return true, err
		}
		// logger.Info(colors.Red, "fileName:", fileName, colors.Off)
		// logger.Info(colors.Red, "curRange:", curRange, colors.Off)
	}

	stageFn, _ = file.LatestFileInFolder(stageFolder) // it may not exist...
	nAppsNow := int(file.FileSize(stageFn) / asciiAppearanceSize)
	opts.Report(nAppsThen, nAppsNow)

	// logger.Info("Removing backup file as it's not needed.")
	os.Remove(backupFn) // commits the change

	return true, err
}

func (opts *ScrapeOptions) Report(nAppsThen, nAppsNow int) {
	msg := "Block={%d} have {%d} appearances of {%d} ({%0.1f%%}). Need {%d} more. Added {%d} records ({%0.2f} apps/blk)."
	need := opts.Settings.AppsPerChunk - utils.Min(opts.Settings.AppsPerChunk, uint64(nAppsNow))
	seen := nAppsNow
	if nAppsThen < nAppsNow {
		seen = nAppsNow - nAppsThen
	}
	pct := float64(nAppsNow) / float64(opts.Settings.AppsPerChunk)
	pBlk := float64(seen) / float64(opts.BlockCnt)
	height := opts.StartBlock + opts.BlockCnt - 1
	msg = strings.Replace(msg, "{", colors.Green, -1)
	msg = strings.Replace(msg, "}", colors.Off, -1)
	logger.Info(fmt.Sprintf(msg, height, nAppsNow, opts.Settings.AppsPerChunk, pct*100, need, seen, pBlk))
}

func isListSequential(chain string, ripeFileList []os.DirEntry, allowMissing bool) error {
	prev := base.NotARange
	for _, file := range ripeFileList {
		fileRange := base.RangeFromFilename(file.Name())
		if prev != base.NotARange && prev != fileRange {
			if !prev.Preceeds(fileRange, !allowMissing) {
				msg := fmt.Sprintf("Ripe files are not sequential (%s ==> %s)", prev, fileRange)
				return errors.New(msg)
			}
		}
		prev = fileRange
	}
	return nil
}

var spaces = strings.Repeat(" ", 40)

// TODO: chifra scrape misreports appearances per block when consolidating #2291 (closed, but copied here as a TODO)
