package fabcli

import (
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
)

// BlockInfo stores information of block.
type BlockInfo struct {
	CreatedOn         int64    `json:"createdOn"`
	BlockNumber       uint64   `json:"blockNumber"`
	BlockHash         string   `json:"blockHash"`
	PreviousBlockHash string   `json:"previousBlockHash"`
	DataHash          string   `json:"dataHash"`
	TransactionIDs    []string `json:"transactionIds"`
}

// BlockInfoList queries and returns BlockInfo list from specified block number in count.
func (p *ChannelContext) BlockInfoList(blockNumber uint64, count int) (infoList []BlockInfo, err error) {
	client, err := p.LedgerClient()
	if err != nil {
		return
	}
	for i := 0; i < count; i++ {
		var info BlockInfo
		info, err = blockInfo(client, blockNumber)
		if err != nil {
			return
		}
		infoList = append(infoList, info)
		if blockNumber == 0 {
			break
		}
		blockNumber--
	}
	return
}

func blockInfo(client *ledger.Client, blockNumber uint64) (info BlockInfo, err error) {
	block, err := client.QueryBlock(blockNumber)
	if err != nil {
		return
	}
	info = BlockInfo{
		BlockNumber:       block.Header.Number,
		BlockHash:         hex.EncodeToString(blockHash(block.Header.Number, block.Header.PreviousHash, block.Header.DataHash)),
		PreviousBlockHash: hex.EncodeToString(block.Header.PreviousHash),
		DataHash:          hex.EncodeToString(block.Header.DataHash),
		TransactionIDs:    make([]string, len(block.Data.Data)),
	}
	for i, data := range block.Data.Data {
		var (
			envelope      = new(common.Envelope)
			payload       = new(common.Payload)
			channelHeader = new(common.ChannelHeader)
		)
		if err = batchExecute(
			func() error { return envelope.XXX_Unmarshal(data) },
			func() error { return payload.XXX_Unmarshal(envelope.Payload) },
			func() error { return channelHeader.XXX_Unmarshal(payload.Header.ChannelHeader) },
		); err != nil {
			return
		}
		timestamp := channelHeader.Timestamp.Seconds
		if info.CreatedOn == 0 || (timestamp > 0 && timestamp < info.CreatedOn) {
			info.CreatedOn = timestamp
		}
		info.TransactionIDs[i] = channelHeader.TxId
	}
	return
}

func blockHash(number uint64, previousHash, dataHash []byte) []byte {
	data, _ := asn1.Marshal(asn1Header{
		Number:       int64(number),
		PreviousHash: previousHash,
		DataHash:     dataHash,
	})
	hash := sha256.Sum256(data)
	return hash[:]
}

type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}
