package eventclient

import (
	"math"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/orderer"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils" //"github.com/hyperledger/fabric/common/tools/protolator"
)

var (
	ledgerClient     *ledger.Client
	Seek             int
	transactionCount int
	roundCount       int
	Producer         sarama.SyncProducer

	oldest  = &orderer.SeekPosition{Type: &orderer.SeekPosition_Oldest{Oldest: &orderer.SeekOldest{}}}
	newest  = &orderer.SeekPosition{Type: &orderer.SeekPosition_Newest{Newest: &orderer.SeekNewest{}}}
	maxStop = &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: math.MaxUint64}}}
	Logger  = flogging.MustGetLogger("eventsclient")
)

const (
	OLDEST      = -2
	NEWEST      = -1
	MSG_TOPIC   = "batch-resp"
	CONFIG_FILE = "/go/src/github.com/m/echo-rest-server/config/"
	CHANNEL_ID  = "m"
	USER_NAME   = "admin"
	PATH_PREFIX = ""
)

// deliverClient abstracts common interface
// for deliver and deliverfiltered services
type DeliverClient interface {
	Send(*common.Envelope) error
	Recv() (*peer.DeliverResponse, error)
}

// eventsClient client to get connected to the
// events peer delivery system
type EventsClient struct {
	Client      DeliverClient
	Signer      crypto.LocalSigner
	TlsCertHash []byte
}

func InitLedgerClient(fileName string) error {
	sdk, err := fabsdk.New(config.FromFile(os.Getenv("HOME") + PATH_PREFIX + CONFIG_FILE + fileName))
	if err != nil {
		return err
	}
	clientProvider := sdk.ChannelContext(CHANNEL_ID, fabsdk.WithUser(USER_NAME))
	ledgerClient, err = ledger.New(clientProvider)
	if err != nil {
		return err
	}
	return nil
}

func (r *EventsClient) seekOldest() error {
	return r.Client.Send(r.seekHelper(oldest, maxStop))
}

func (r *EventsClient) seekNewest() error {
	return r.Client.Send(r.seekHelper(newest, maxStop))
}

func (r *EventsClient) seekSingle(blockNumber uint64) error {
	specific := &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: blockNumber}}}
	return r.Client.Send(r.seekHelper(specific, specific))
}

func (r *EventsClient) seekHelper(start *orderer.SeekPosition, stop *orderer.SeekPosition) *common.Envelope {
	env, err := utils.CreateSignedEnvelopeWithTLSBinding(common.HeaderType_DELIVER_SEEK_INFO, "m", r.Signer, &orderer.SeekInfo{
		Start:    start,
		Stop:     stop,
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}, 0, 0, r.TlsCertHash)
	if err != nil {
		panic(err)
	}
	return env
}

func (r *EventsClient) Seek(s int) error {
	var err error
	switch Seek {
	case OLDEST:
		err = r.seekOldest()
	case NEWEST:
		err = r.seekNewest()
	default:
		err = r.seekSingle(uint64(Seek))
	}
	return err
}

func (r *EventsClient) ReadEventsStream() {
	_, _ = r.Client.Recv()
	for {
		msg, err := r.Client.Recv()
		if err != nil {
			Logger.Info("Error receiving:", err)
			return
		}

		block := msg.GetBlock()
		//_ = protolator.DeepMarshalJSON(os.Stdout, block)
		for _, data := range block.GetData().GetData() {
			env, _ := utils.GetEnvelopeFromBlock(data)
			cAction, _ := utils.GetActionFromEnvelopeMsg(env)
			channelHeader, err := utils.ChannelHeader(env)
			if err != nil {
				Logger.Error("error:", err)
			}
			Logger.Info("block number:", block.GetHeader().GetNumber(), "/ chaincode id information -", cAction.GetChaincodeId(), "/ transaction id:", channelHeader.GetTxId())
			cEvent, _ := utils.GetChaincodeEvents(cAction.GetEvents())
			transaction, err := ledgerClient.QueryTransaction(fab.TransactionID(cEvent.GetTxId()))
			if err != nil {
				Logger.Error("error:", err)
			}
			Logger.Info("validation:", peer.TxValidationCode_VALID)
			if transaction.GetValidationCode() == int32(peer.TxValidationCode_VALID) {
				handleTransaction(cEvent)
			} else {
				cEventName := cEvent.GetEventName()
				cEventPayloadAdditionalData := strings.SplitN(cEventName, "_", 2)
				Logger.Error("transaction error:", peer.TxValidationCode_name[int32(transaction.ValidationCode)], "/ id:", cEventPayloadAdditionalData[1])
			}
		}
	}
}

func handleTransaction(cEvent *peer.ChaincodeEvent) {
	cEventName := cEvent.GetEventName()
	cEventPayloadAdditionalData := strings.SplitN(cEventName, "_", 2)
	if len(cEventPayloadAdditionalData) == 2 {
		cEventPayloadAddtiionalDataJsonFields := ",\"svc\":\"" + cEventPayloadAdditionalData[0] + "\",\"id\":\"" + cEventPayloadAdditionalData[1] + "\",\"code\":\"RS_RES_N01\"}"
		cEventPayload := string(cEvent.GetPayload())
		cEventPayloadJsonStr := cEventPayload[:len(cEventPayload)-1] + cEventPayloadAddtiionalDataJsonFields
		Logger.Info("chaincode id:", cEvent.GetChaincodeId(), "/ transaction id:", cEvent.GetTxId(), "/ chaincode event name:", cEventName, "/ chaincode event payload:", cEventPayload, "/ chaincode event payload for kafka resp:", cEventPayloadJsonStr)
		partition, offset, err := Producer.SendMessage(
			&sarama.ProducerMessage{
				Topic: MSG_TOPIC,
				Value: sarama.StringEncoder(cEventPayloadJsonStr),
			})
		if err != nil {
			Logger.Error(err, "chaincode id:", cEvent.GetChaincodeId(), "/ transaction id:", cEvent.GetTxId(), "/ chaincode event name:", cEventName, "/ chaincode event payload:", cEventPayload, "/ chaincode event payload for kafka resp:", cEventPayloadJsonStr)
			return
		}

		Logger.Info("Message is stored in topic(", MSG_TOPIC, ")/partition(", partition, ")/offset(", offset, ")")
	}

}
