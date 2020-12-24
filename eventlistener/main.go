/*
FABRIC_CFG_PATH must point at configtx.yaml file

export FABRIC_CFG_PATH=/home/vagrant/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/

*/

package main

import (
	"context"
	"flag"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/localmsp"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/comm"
	"github.com/hyperledger/fabric/protos/orderer"
	"github.com/hyperledger/fabric/protos/peer"

	"github.com/m/eventlistener/eventclient"
	"github.com/m/eventlistener/setup"
)

var (
	channelID string
	filtered  bool

	oldest  = &orderer.SeekPosition{Type: &orderer.SeekPosition_Oldest{Oldest: &orderer.SeekOldest{}}}
	newest  = &orderer.SeekPosition{Type: &orderer.SeekPosition_Newest{Newest: &orderer.SeekNewest{}}}
	maxStop = &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: math.MaxUint64}}}

	mm     peer.FilteredTransaction_TransactionActions
	logger = flogging.MustGetLogger("eventsclient")
)

func init() {
	os.Setenv("FABRIC_CFG_PATH", os.Getenv("HOME")+eventclient.PATH_PREFIX+"/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/")
	setup.ReadCLInputs()
	setup.InitKafkaProducer()
	setup.InitConfig()
	setup.InitMSP()
	setup.InitLedgerClientSetup()
}

func main() {
	if eventclient.Seek < eventclient.OLDEST {
		logger.Info("Invalid seek value")
		flag.PrintDefaults()
		return
	}

	clientConfig := comm.ClientConfig{
		KaOpts:  comm.DefaultKeepaliveOptions,
		SecOpts: &comm.SecureOptions{},
		Timeout: 5 * time.Minute,
	}

	if setup.TlsEnabled {
		clientConfig.SecOpts.UseTLS = true
		rootCert, err := ioutil.ReadFile(setup.ServerRootCAPath)
		if err != nil {
			logger.Error("error loading TLS root certificate", err)
			return
		}
		clientConfig.SecOpts.ServerRootCAs = [][]byte{rootCert}
		if setup.MTlsEnabled {
			clientConfig.SecOpts.RequireClientCert = true
			clientKey, err := ioutil.ReadFile(setup.ClientKeyPath)
			if err != nil {
				logger.Error("error loading client TLS key from", setup.ClientKeyPath)
				return
			}
			clientConfig.SecOpts.Key = clientKey

			clientCert, err := ioutil.ReadFile(setup.ClientCertPath)
			if err != nil {
				logger.Error("error loading client TLS certificate from path", setup.ClientCertPath)
				return
			}
			clientConfig.SecOpts.Certificate = clientCert
		}
	}

	grpcClient, err := comm.NewGRPCClient(clientConfig)
	if err != nil {
		logger.Error("Error creating grpc client:", err)
		return
	}

	var client eventclient.DeliverClient
	for _, serverAddr := range setup.ServerAddrs {
		conn, err := grpcClient.NewConnection(serverAddr, "")
		if err != nil {
			logger.Warn("Error connecting:", err)
		} else {
			client, err = peer.NewDeliverClient(conn).Deliver(context.Background())
			if err != nil {
				logger.Warn("Error connecting:", err)
				continue
			}
			logger.Info("eventclient is ready with ", serverAddr)

			break
		}
	}
	if err != nil {
		logger.Error("Error connecting any peers")
		return
	}

	//client, err := peer.NewDeliverClient(conn).DeliverFiltered(context.Background())
	//client, err := peer.NewChaincodeSupportClient(conn).Register(context.Background())

	events := &eventclient.EventsClient{
		Client: client,
		Signer: localmsp.NewSigner(),
	}

	if setup.MTlsEnabled {
		events.TlsCertHash = util.ComputeSHA256(grpcClient.Certificate().Certificate[0])
	}
	events.Seek(eventclient.Seek)
	if err != nil {
		logger.Error("Received error:", err)
		return
	}
	events.ReadEventsStream()
}
