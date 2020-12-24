package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
	"github.com/pkg/errors"
)

func (setup *FabricSetup) GetTransactionByTxId(txId string) (string, error) {

	clientProvider := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	client, err := ledger.New(clientProvider)
	if err != nil {
		return "", errors.WithMessage(err, "failed to create ledger client")
	}

	transaction, err := client.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		return "", errors.WithMessage(err, "failed to get Transaction data")
	}
	env := transaction.GetTransactionEnvelope()
	newEnv := common.Envelope(*env)
	cAction, err := utils.GetActionFromEnvelopeMsg(&newEnv)
	if err != nil {
		return "", errors.WithMessage(err, "failed to create chaincode action structure")
	}

	return string(cAction.GetResponse().GetPayload()), nil
}

// qeury
func (setup *FabricSetup) Query(ccst *ChainCodeInfo) (string, error) {

	response, err := setup.client.Query(channel.Request{ChaincodeID: ccst.CcId, Fcn: ccst.CcFnc, Args: ccst.CcArgs})

	return string(response.Payload), err
}

// invoke
func (setup *FabricSetup) Invoke(ccst *ChainCodeInfo) (string, error) {

	transientDataMap := make(map[string][]byte)
	response, err := setup.client.Execute(channel.Request{ChaincodeID: ccst.CcId, Fcn: ccst.CcFnc, Args: ccst.CcArgs, TransientMap: transientDataMap})

	return string(response.Payload), err
}
