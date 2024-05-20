package pinning

import (
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/theQRL/trueblocks-core/src/apps/chifra/pkg/config"
)

// pinFileLocally pins the file locally to the local IPFS node
func (s *Service) pinFileLocally(chain, filepath string) (base.IpfsHash, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sh := shell.NewShell(config.GetPinning().LocalPinUrl)
	cid, err := sh.Add(file, shell.Pin(true))
	if err != nil {
		return "", err
	}
	return base.IpfsHash(cid), nil
}
