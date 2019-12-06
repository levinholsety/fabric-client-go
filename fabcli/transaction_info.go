package fabcli

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// TransactionInfo stores information of transaction.
type TransactionInfo struct {
	ID        string              `json:"id"`
	CreatedOn int64               `json:"createdOn"`
	ChannelID string              `json:"channelId"`
	Actions   []transactionAction `json:"actions"`
}

type transactionAction struct {
	Chaincode chaincode `json:"chaincode"`
	RWSets    []rwSet   `json:"rwSets"`
	InputArgs []string  `json:"inputArgs"`
}

type chaincode struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type rwSet struct {
	Namespace string  `json:"namespace"`
	Reads     []read  `json:"reads,omitempty"`
	Writes    []write `json:"writes,omitempty"`
}

type read struct {
	Key     string  `json:"key"`
	Version version `json:"version,omitempty"`
}

type version struct {
	BlockNum uint64 `json:"blockNum"`
	TxNum    uint64 `json:"txNum"`
}

type write struct {
	Key      string `json:"key"`
	IsDelete bool   `json:"isDelete"`
	Value    string `json:"value"`
}

// TransactionInfo queries and returns TransactionInfo.
func (p *ChannelContext) TransactionInfo(txID string) (txInfo TransactionInfo, err error) {
	ledgerClient, err := p.LedgerClient()
	if err != nil {
		return
	}
	tx, err := ledgerClient.QueryTransaction(fab.TransactionID(txID))
	if err != nil {
		return
	}
	var (
		payload       = new(common.Payload)
		channelHeader = new(common.ChannelHeader)
		transaction   = new(peer.Transaction)
	)
	if err = batchExecute(
		func() error { return payload.XXX_Unmarshal(tx.TransactionEnvelope.Payload) },
		func() error { return channelHeader.XXX_Unmarshal(payload.Header.ChannelHeader) },
		func() error { return transaction.XXX_Unmarshal(payload.Data) },
	); err != nil {
		return
	}
	txInfo = TransactionInfo{
		ID:        txID,
		CreatedOn: channelHeader.Timestamp.Seconds,
		ChannelID: channelHeader.ChannelId,
		Actions:   make([]transactionAction, len(transaction.Actions)),
	}
	for i, action := range transaction.Actions {
		if txInfo.Actions[i], err = getTransactionAction(action); err != nil {
			return
		}
	}
	return
}

func getTransactionAction(action *peer.TransactionAction) (act transactionAction, err error) {
	var (
		ccActPayload        = new(peer.ChaincodeActionPayload)
		ccProposalPayload   = new(peer.ChaincodeProposalPayload)
		proposalRespPayload = new(peer.ProposalResponsePayload)
		ccInvocationSpec    = new(peer.ChaincodeInvocationSpec)
		ccAction            = new(peer.ChaincodeAction)
		txrwset             = new(rwset.TxReadWriteSet)
	)
	if err = batchExecute(
		func() error { return ccActPayload.XXX_Unmarshal(action.Payload) },
		func() error { return ccProposalPayload.XXX_Unmarshal(ccActPayload.ChaincodeProposalPayload) },
		func() error { return proposalRespPayload.XXX_Unmarshal(ccActPayload.Action.ProposalResponsePayload) },
		func() error { return ccInvocationSpec.XXX_Unmarshal(ccProposalPayload.Input) },
		func() error { return ccAction.XXX_Unmarshal(proposalRespPayload.Extension) },
		func() error { return txrwset.XXX_Unmarshal(ccAction.Results) },
	); err != nil {
		return
	}
	rwsets, err := getRWSets(txrwset.NsRwset)
	if err != nil {
		return
	}
	act = transactionAction{
		RWSets: rwsets,
		Chaincode: chaincode{
			Name:    ccAction.ChaincodeId.Name,
			Version: ccAction.ChaincodeId.Version,
		},
		InputArgs: getInputArgs(ccInvocationSpec.ChaincodeSpec.Input.Args),
	}
	return
}

func getRWSets(nsrwsets []*rwset.NsReadWriteSet) (sets []rwSet, err error) {
	sets = make([]rwSet, len(nsrwsets))
	for i, nsrwset := range nsrwsets {
		kvrwset := new(kvrwset.KVRWSet)
		err = kvrwset.XXX_Unmarshal(nsrwset.Rwset)
		if err != nil {
			return
		}
		sets[i] = rwSet{
			Namespace: nsrwset.Namespace,
			Reads:     getReads(kvrwset),
			Writes:    getWrites(kvrwset),
		}
	}
	return
}

func getReads(kvrwset *kvrwset.KVRWSet) (reads []read) {
	reads = make([]read, len(kvrwset.Reads))
	for j, r := range kvrwset.Reads {
		reads[j] = read{
			Key:     r.Key,
			Version: getVersion(r.Version),
		}
	}
	return
}

func getVersion(v *kvrwset.Version) (ver version) {
	if v == nil {
		return
	}
	return version{
		BlockNum: v.BlockNum,
		TxNum:    v.TxNum,
	}
}

func getWrites(kvrwset *kvrwset.KVRWSet) (writes []write) {
	writes = make([]write, len(kvrwset.Writes))
	for j, w := range kvrwset.Writes {
		writes[j] = write{
			Key:      w.Key,
			IsDelete: w.IsDelete,
			Value:    string(w.Value),
		}
	}
	return
}

func getInputArgs(dataArray [][]byte) (args []string) {
	args = make([]string, len(dataArray))
	for i, data := range dataArray {
		args[i] = string(data)
	}
	return
}
