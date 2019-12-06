package fabcli

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
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
