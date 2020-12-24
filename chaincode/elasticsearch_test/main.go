package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/olivere/elastic"
)

var (
	esClient *elastic.Client
	ctext    context.Context
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
	Id      int
	Title   string
	Year    string
	Created time.Time
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data, so be careful to avoid a scenario where you
// inadvertently clobber your ledger's data!
func (t *SimpleAsset) Init(shim.ChaincodeStubInterface) peer.Response {
	var err error
	esClient, err = elastic.NewClient(elastic.SetURL("http://172.17.8.11:9200", "http://172.17.8.12:9200", "http://172.17.8.13:9200"))
	if err != nil {
		// Handle error
		return shim.Error(err.Error())
	}
	ctext = context.Background()
	exists, err := esClient.IndexExists("test_chaincode").Do(ctext)
	if err != nil {
		// Handle error
		return shim.Error(err.Error())
	}
	if !exists {
		mapping := `
{
    "settings":{
        "number_of_shards":3,
        "number_of_replicas":3
    },
    "mappings":{
            "properties":{
				"id":{
                    "type":"integer"
                },
                "title":{
                    "type":"text"
                },
                "year":{
                    "type":"text"
                }
          }
    }
}
`
		createIndex, err := esClient.CreateIndex("test_chaincode").Body(mapping).Do(ctext)
		if err != nil {
			// Handle error
			return shim.Error(err.Error())
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			return shim.Error("create Index not acknowledgerd")
		}
	}

	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else {
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	arg0Int, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("arg type is wrong")
	}
	arg3Time, err := time.Parse(time.RFC3339, args[3])
	if err != nil {
		return "", fmt.Errorf("arg type is wrong")
	}

	var asset = SimpleAsset{
		Id:      arg0Int,
		Title:   args[1],
		Year:    args[2],
		Created: arg3Time,
	}

	bulkRequest := esClient.Bulk()
	req := elastic.NewBulkIndexRequest().Index("test_chaincode").Id(args[0]).Doc(asset)
	bulkRequest = bulkRequest.Add(req)
	_, err = bulkRequest.Do(ctext)
	if err != nil {
		return "", err
	}

	assetBytes, _ := json.Marshal(asset)

	err = stub.PutState("Test_"+args[0], assetBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
