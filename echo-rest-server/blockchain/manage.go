package blockchain

import (
	"fmt"
	//"strconv"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	//	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"encoding/json"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
	"os"
	//"github.com/m/echo-rest-server/message"
)

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) SetupChannel(s interface{}) error {

	jsonBytes2, _ := json.Marshal(s)

	var ccst ChainCodeInfo
	err2 := json.Unmarshal(jsonBytes2, &ccst)
	if err2 != nil {
		panic(err2)
	}
	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(setup.sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}
	fmt.Println("mspClient pass")
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "failed to get admin signing identity")
	}
	fmt.Println("adminIdentity pass")
	fmt.Println(adminIdentity)

	req := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	fmt.Println(req)

	txID, err := setup.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	fmt.Println("Channel created")

	// Make admin user join the previously created channel
	if err = setup.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
		return errors.WithMessage(err, "failed to make admin join channel")
	}
	fmt.Println("Channel joined")

	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	// Creation of the client which will enables access to our channel events
	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")

	return nil
}

func (setup *FabricSetup) SetupCC(s interface{}) error {

	jsonBytes2, _ := json.Marshal(s)

	var ccst ChainCodeInfo
	err2 := json.Unmarshal(jsonBytes2, &ccst)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(ccst)

	// Create the chaincode package that will be sent to the peers
	ccPkg, err := packager.NewCCPackage(ccst.CcPath, os.Getenv("GOPATH"))
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("ccPkg created")

	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: ccst.CcId, Path: ccst.CcPath, Version: ccst.CcVs, Package: ccPkg}
	installRes, err := setup.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	fmt.Println("-------------Chaincode instal  result----------------")
	//fmt.Println(installRes, err)
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}

	if len(installRes) > 0 {
		for _, response := range installRes {

			fmt.Println(response)
		}
	}

	fmt.Println("---start---------Chaincode instantiated------------------")

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{ccst.CcPolicy})

	resp, err := setup.admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: ccst.CcId, Path: ccst.CcPath, Version: ccst.CcVs, Args: [][]byte{[]byte(ccst.CcInit)}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	return nil
}
