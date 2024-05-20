package chunksPkg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/manifest"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/pinning"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/usage"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/walk"
)

func (opts *ChunksOptions) HandlePin(blockNums []uint64) error {
	chain := opts.Globals.Chain
	if opts.Globals.TestMode {
		logger.Warn("Pinning option not tested.")
		return nil
	}

	if !opts.Globals.IsApiMode() && usage.QueryUser(pinWarning, "Check skipped") {
		if err := opts.doCheck(blockNums); err != nil {
			return err
		}
	}

	firstBlock := mustParseUint(os.Getenv("TB_CHUNKS_PINFIRSTBLOCK"))
	lastBlock := mustParseUint(os.Getenv("TB_CHUNKS_PINLASTBLOCK"))
	if lastBlock == 0 {
		lastBlock = utils.NOPOS
	}

	outPath := filepath.Join(config.PathToCache(chain), "tmp", "manifest.json")
	if opts.Rewrite {
		outPath = config.PathToManifest(chain)
	}

	man, err := manifest.ReadManifest(chain, opts.PublisherAddr, manifest.LocalCache)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	fetchData := func(modelChan chan types.Modeler[types.RawModeler], errorChan chan error) {
		hash := base.BytesToHash(config.HeaderHash(config.ExpectedVersion()))
		report := simpleChunkPinReport{
			Version:  config.VersionTags[hash.Hex()],
			Chain:    chain,
			SpecHash: base.IpfsHash(manifest.Specification()),
		}

		fileList := make([]string, 0, len(man.Chunks))
		listFiles := func(walker *walk.CacheWalker, path string, first bool) (bool, error) {
			rng, err := base.RangeFromFilenameE(path)
			if err != nil {
				return false, err
			}
			if rng.Last < firstBlock || rng.First > lastBlock {
				logger.Info("Skipping", path)
				return true, nil
			}
			if path != index.ToBloomPath(path) {
				return false, fmt.Errorf("should not happen in pinChunk")
			}
			if opts.Deep || len(blockNums) > 0 || man.ChunkMap[rng.String()] == nil {
				fileList = append(fileList, path)
			}
			return true, nil
		}

		walker := walk.NewCacheWalker(
			chain,
			opts.Globals.TestMode,
			100, /* maxTests */
			listFiles,
		)
		if err := walker.WalkBloomFilters(blockNums); err != nil {
			errorChan <- err
			// TODO: cancel probably doesn't cancel anything here does it? The walker doesn't even see it.
			cancel()
			return
		}

		sort.Slice(fileList, func(i, j int) bool {
			rng1, _ := base.RangeFromFilenameE(fileList[i])
			rng2, _ := base.RangeFromFilenameE(fileList[j])
			return rng1.First < rng2.First
		})

		for _, path := range fileList {
			if opts.Globals.Verbose {
				logger.Info("pinning path:", path)
			}
			local, remote, err := pinning.PinOneChunk(chain, path, opts.Remote)
			if err != nil {
				errorChan <- err
				logger.Error("Pin failed:", path, err)
			}

			blMatches, idxMatches := matches(&local, &remote)
			opts.matchReport(blMatches, local.BloomHash, remote.BloomHash)
			opts.matchReport(idxMatches, local.IndexHash, remote.IndexHash)

			if opts.Remote {
				man.Chunks = append(man.Chunks, remote)
			} else {
				man.Chunks = append(man.Chunks, local)
			}
			_ = man.SaveManifest(chain, outPath)

			if opts.Globals.Verbose {
				if opts.Remote {
					fmt.Println("result.Remote:", remote.String())
				} else {
					fmt.Println("result.Local:", local.String())
				}
			}

			sleep := opts.Sleep
			if sleep > 0 {
				ms := time.Duration(sleep*1000) * time.Millisecond
				if !opts.Globals.TestMode {
					logger.Info(fmt.Sprintf("Sleeping for %g seconds", sleep))
				}
				time.Sleep(ms)
			}
		}

		if len(blockNums) == 0 && firstBlock == 0 && lastBlock == utils.NOPOS {
			tsPath := config.PathToTimestamps(chain)
			if localHash, remoteHash, err := pinning.PinOneFile(chain, "timestamps", tsPath, opts.Remote); err != nil {
				errorChan <- err
				logger.Error("Pin failed:", tsPath, err)
			} else {
				opts.matchReport(localHash == remoteHash, localHash, remoteHash)
				report.TimestampHash = localHash
			}

			manPath := config.PathToManifest(chain)
			if opts.Deep {
				manPath = outPath
			}
			if localHash, remoteHash, err := pinning.PinOneFile(chain, "manifest", manPath, opts.Remote); err != nil {
				errorChan <- err
				logger.Error("Pin failed:", manPath, err)
			} else {
				opts.matchReport(localHash == remoteHash, localHash, remoteHash)
				report.ManifestHash = localHash
			}
		}

		logger.Info("The new manifest was written to", colors.BrightGreen+outPath+colors.Off, len(man.Chunks), "chunks")

		modelChan <- &report
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}

// matches returns true if the Result has both local and remote hashes for both the index and the bloom and they match
func matches(local, remote *types.SimpleChunkRecord) (bool, bool) {
	return local.BloomHash == remote.BloomHash, local.IndexHash == remote.IndexHash
}

func (opts *ChunksOptions) matchReport(matches bool, localHash, remoteHash base.IpfsHash) {
	if !opts.Remote || !config.IpfsRunning() {
		return // if we're not pinning in two places, don't report on matches
	}

	if matches {
		logger.Info(colors.BrightGreen+"Matches: "+localHash.String(), " ", localHash, colors.Off)
	} else {
		logger.Warn("Pins mismatch:", localHash.String(), " ", localHash)
	}
}

func (opts *ChunksOptions) doCheck(blockNums []uint64) error {
	if err, ok := opts.check(blockNums, false /* silent */); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("checks failed")
	}
	return nil
}

var pinWarning = `Do you want to run --check first (Yn)? `
