package util

import (
	"encoding/json"
	"fmt"
	//"os"
	"reflect"
)

func ExtractFilteredBlock(buf []byte) []interface{} {
	var data map[string]interface{}
	err := json.Unmarshal(buf, &data)
	if err != nil {
		fmt.Printf("Error pretty printing block: %s", err)
	}
	return reflect.ValueOf(data["filtered_transactions"]).Interface().([]interface{})
}

func ExtractChaincodeEvents(filteredTransactionList []interface{}) (chaincodeEvents []map[string]interface{}) {
	chaincodeEvents = []map[string]interface{}{}
	for _, filteredTransaction := range filteredTransactionList {
		chaincodeAction := reflect.ValueOf(ExtractMap(filteredTransaction, "transaction_actions")["chaincode_actions"]).Interface().([]interface{})
		chaincodeEvents = append(chaincodeEvents, ExtractMap(chaincodeAction[0], "chaincode_event"))
	}
	return chaincodeEvents
}

func ExtractMap(data interface{}, fields ...string) map[string]interface{} {
	mapData := data.(map[string]interface{})
	for _, field := range fields {
		mapData = mapData[field].(map[string]interface{})
	}
	return mapData
}
