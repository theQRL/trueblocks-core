package types

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/ethereum/go-ethereum/common"
)

type IpfsHash string

func (h IpfsHash) String() string {
	return string(h)
}

type SimpleTimestamp struct {
	BlockNumber uint64 `json:"blockNumber"`
	TimeStamp   uint64 `json:"timestamp"`
	Diff        uint64 `json:"diff"`
}

type SimpleNamedBlock struct {
	BlockNumber uint64 `json:"blockNumber"`
	TimeStamp   uint64 `json:"timestamp"`
	Date        string `json:"date"`
	Name        string `json:"name,omitempty"`
}

type SimpleAppearance struct {
	Address          string `json:"address"`
	BlockNumber      uint32 `json:"blockNumber"`
	TransactionIndex uint32 `json:"transactionIndex"`
}

type SimpleFunction struct {
	Encoding  string `json:"encoding,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type SimpleMonitor struct {
	Address     string `json:"address"`
	NRecords    int    `json:"nRecords"`
	FileSize    int64  `json:"fileSize"`
	LastScanned uint32 `json:"lastScanned"`
}

type SimpleChunkRecord struct {
	FileName  string `json:"fileName"`
	BloomHash string `json:"bloomHash"`
	IndexHash string `json:"indexHash"`
}

type SimpleBloom struct {
	Range     cache.FileRange `json:"range"`
	Count     uint32          `json:"nBlooms"`
	NInserted uint64          `json:"nInserted"`
	Size      int64           `json:"size"`
	Width     uint64          `json:"byteWidth"`
}

type SimpleIndex struct {
	Range           cache.FileRange `json:"range"`
	Magic           uint32          `json:"magic"`
	Hash            common.Hash     `json:"hash"`
	AddressCount    uint32          `json:"nAddresses"`
	AppearanceCount uint32          `json:"nAppearances"`
	Size            int64           `json:"fileSize"`
}

type SimpleIndexAppearance struct {
	BlockNumber      uint32 `json:"blockNumber"`
	TransactionIndex uint32 `json:"transactionIndex"`
}

type SimpleIndexAddress struct {
	Address string `json:"address"`
	Offset  uint32 `json:"offset"`
	Count   uint32 `json:"count"`
}

type SimpleIndexAddressBelongs struct {
	Address string                  `json:"address"`
	Offset  uint32                  `json:"offset"`
	Count   uint32                  `json:"count"`
	Apps    []SimpleIndexAppearance `json:"apps"`
}
