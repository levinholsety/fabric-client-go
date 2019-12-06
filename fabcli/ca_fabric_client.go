package fabcli

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// NewCAFabricClient creates fabric client with ca server.
func NewCAFabricClient(configPath, user, secret string) *FabricClient {
	return &FabricClient{
		ConfigPath: configPath,
		channelContext: func(sdk *fabsdk.FabricSDK, mspClient *msp.Client, channelID string) (ctx context.ChannelProvider, err error) {
			_, err = mspClient.GetSigningIdentity(user)
			if err == msp.ErrUserNotFound {
				err = mspClient.Enroll(user, msp.WithSecret(secret))
				if err != nil {
					return
				}
			} else if err != nil {
				return
			}
			ctx = sdk.ChannelContext(channelID, fabsdk.WithUser(user))
			return
		},
	}
}
