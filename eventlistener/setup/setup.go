package setup

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric/peer/common"
	"github.com/spf13/viper"

	"github.com/m/eventlistener/eventclient"
	"github.com/m/eventlistener/util"
)

var (
	config           string
	ServerAddrs      []string
	TlsEnabled       bool
	MTlsEnabled      bool
	ClientKeyPath    string
	ClientCertPath   string
	ServerRootCAPath string

	channelID = "m"
	// server.key
	clientKey = os.Getenv("HOME") + eventclient.PATH_PREFIX + "/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls/keystore/server.key"
	// server.crt
	clientCert = os.Getenv("HOME") + eventclient.PATH_PREFIX + "/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls/signcerts/server.crt"
	// ca.crt
	rootCert      = os.Getenv("HOME") + eventclient.PATH_PREFIX + "/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/org1-ca-chain.pem"
	defaultConfig = "fingerCloud"
)

const (
	ROOT       = "configtx"
	KAFKA_PORT = "9092"
)

func ReadCLInputs() {
	ServerAddrs = []string{"peer0.org1.hf.m.io:7051", "peer1.org1.hf.m.io:7051", "peer2.org1.hf.m.io:7051"}
	flag.StringVar(&channelID, "channelID", channelID, "The channel ID to deliver from.")
	flag.BoolVar(&TlsEnabled, "tls", true, "TLS enabled/disabled")
	flag.BoolVar(&MTlsEnabled, "mTls", true, "Mutual TLS enabled/disabled (whenever server side validates clients TLS certificate)")
	flag.StringVar(&ClientKeyPath, "clientKey", clientKey, "Specify path to the client TLS key")
	flag.StringVar(&ClientCertPath, "clientCert", clientCert, "Specify path to the client TLS certificate")
	flag.StringVar(&ServerRootCAPath, "rootCert", rootCert, "Specify path to the server root CA certificate")
	flag.IntVar(&eventclient.Seek, "seek", eventclient.NEWEST, "Specify the range of requested blocks."+
		"Acceptable values:"+
		"-2 (or -1) to start from oldest (or newest) and keep at it indefinitely."+
		"N >= 0 to fetch block N only.")
	flag.StringVar(&config, "config", defaultConfig, "The config to deliver from.")
	flag.Parse()
}

func InitMSP() {
	// Init the MSP
	var mspMgrConfigDir = os.Getenv("HOME") + eventclient.PATH_PREFIX + "/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/peerOrganizations/org1.hf.m.io/users/admin/msp"
	var mspType = "bccsp"
	var mspID = "org1.hf.m.io"

	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}
	err := common.InitCrypto(mspMgrConfigDir, mspID, mspType)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Sprintf("Cannot run client because %s", err.Error()))
	}
}

func InitConfig() {
	// For environment variables.
	viper.SetEnvPrefix(ROOT)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	err := common.InitConfig(ROOT)
	if err != nil { // Handle errors reading the config file
		eventclient.Logger.Error(fmt.Sprintf("fatal error when initializing %s config : %s", ROOT, err))
		return
	}
}

func InitKafkaProducer() {
	var err error

	//var err error
	pConfig := sarama.NewConfig()
	pConfig.Producer.RequiredAcks = sarama.WaitForLocal // Wait for all in-sync replicas to ack the message
	pConfig.Producer.Retry.Max = 10                     // Retry up to 10 times to produce the message
	pConfig.Producer.Return.Successes = true

	configuration := getConfiguration(config)
	if configuration == nil {
		return
	}
	kafkaURLs := reflect.ValueOf(util.ExtractMap(configuration, "server", config)["kafka"]).Interface().([]interface{})
	kafkaURLStrs := []string{}
	for _, kafkaURL := range kafkaURLs {
		kafkaURLStrs = append(kafkaURLStrs, kafkaURL.(string)+":"+KAFKA_PORT)
	}
	fmt.Println(kafkaURLStrs)
	eventclient.Producer, err = sarama.NewSyncProducer(kafkaURLStrs, pConfig)
	if err != nil {
		eventclient.Logger.Error("Failed to start Sarama producer:", err)
		return
	}

	eventclient.Logger.Info(config, "producer is ready")
}

func InitLedgerClientSetup() {
	configuration := getConfiguration(config)
	if configuration == nil {
		eventclient.Logger.Error("Failed to load configuration")
		return
	}
	clientConfigFile := reflect.ValueOf(util.ExtractMap(configuration, "server", config)["ledgerClient"]).Interface().(interface{})
	err := eventclient.InitLedgerClient(clientConfigFile.(string))
	if err != nil {
		eventclient.Logger.Error("Failed to create ledger client:", err)
		return
	}
}

func getConfiguration(config string) map[string]interface{} {
	if config != "production" && config != "staging" && config != "development" && config != "local" && config != "fingerCloud" {
		eventclient.Logger.Error("the config flag is wrong")
		return nil
	}
	filename, _ := filepath.Abs(os.Getenv("HOME") + eventclient.PATH_PREFIX + "/go/src/github.com/m/eventlistener/config.json")

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		eventclient.Logger.Error("Failed to read config file", err)
		return nil
	}

	systemConfig := make(map[string]interface{})
	err = json.Unmarshal(yamlFile, &systemConfig)
	if err != nil {
		eventclient.Logger.Error("Failed to unmarshal config file", err)
		return nil
	}

	return util.ExtractMap(systemConfig)
}
