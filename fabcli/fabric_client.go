package fabcli

import (
	"errors"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// FabricClient provides functions to operate fabric.
type FabricClient struct {
	ConfigPath     string
	channelContext func(*fabsdk.FabricSDK, *msp.Client, string) (context.ChannelProvider, error)
}

// Execute executes operations on fabric.
func (c *FabricClient) Execute(channelID string, executor func(ctx *ChannelContext) error) error {
	sdk, err := fabsdk.New(config.FromFile(c.ConfigPath))
	if err != nil {
		return err
	}
	defer sdk.Close()
	mspClient, err := msp.New(sdk.Context())
	if err != nil {
		return err
	}
	ctx, err := c.channelContext(sdk, mspClient, channelID)
	if err != nil {
		return err
	}
	return executor(&ChannelContext{ctx: ctx})
}

// ChannelIDs returns all channel ID of current fabric.
func (c *FabricClient) ChannelIDs() (channelIDs []string, err error) {
	sdk, err := fabsdk.New(config.FromFile(c.ConfigPath))
	if err != nil {
		return
	}
	defer sdk.Close()
	cb, err := sdk.Config()
	if err != nil {
		return
	}
	v, ok := cb.Lookup("channels")
	if !ok {
		err = errors.New("cannot find channels")
		return
	}
	for channelID := range v.(map[string]interface{}) {
		channelIDs = append(channelIDs, channelID)
	}
	return
}
