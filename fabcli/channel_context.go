package fabcli

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// ChannelContext wraps fabric channel context.
type ChannelContext struct {
	ctx context.ChannelProvider
}

// LedgerClient returns *ledger.Client.
func (p *ChannelContext) LedgerClient() (client *ledger.Client, err error) {
	client, err = ledger.New(p.ctx)
	return
}

// ChannelClient returns *channel.Client.
func (p *ChannelContext) ChannelClient() (client *channel.Client, err error) {
	client, err = channel.New(p.ctx)
	return
}

// Invoke invokes chaincode to save data.
func (p *ChannelContext) Invoke(chaincodeID, fcn string, args [][]byte) (transactionID fab.TransactionID, payload []byte, err error) {
	cli, err := p.ChannelClient()
	if err != nil {
		return
	}
	resp, err := cli.Execute(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        args,
	})
	if err != nil {
		return
	}
	transactionID = resp.TransactionID
	payload = resp.Payload
	return
}

// Query invokes chaincode to query data.
func (p *ChannelContext) Query(chaincodeID, fcn string, args [][]byte) (transactionID fab.TransactionID, payload []byte, err error) {
	cli, err := p.ChannelClient()
	if err != nil {
		return
	}
	resp, err := cli.Query(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        args,
	})
	if err != nil {
		return
	}
	transactionID = resp.TransactionID
	payload = resp.Payload
	return
}
