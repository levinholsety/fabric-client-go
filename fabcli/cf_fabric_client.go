package fabcli

import (
	"io/ioutil"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// NewCFFabricClient creates fabric client with cert files.
func NewCFFabricClient(configPath, certPath, keyPath string) *FabricClient {
	return &FabricClient{
		ConfigPath: configPath,
		channelContext: func(sdk *fabsdk.FabricSDK, mspClient *msp.Client, channelID string) (ctx context.ChannelProvider, err error) {
			var certData, keyData []byte
			if certData, err = ioutil.ReadFile(certPath); err != nil {
				return
			}
			if keyData, err = ioutil.ReadFile(keyPath); err != nil {
				return
			}
			id, err := mspClient.CreateSigningIdentity(mspctx.WithCert(certData), mspctx.WithPrivateKey(keyData))
			if err != nil {
				return
			}
			ctx = sdk.ChannelContext(channelID, fabsdk.WithIdentity(id))
			return
		},
	}
}
