package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	CHANNEL          = "m"
	TEST_TARGET      = "main_"
	CONFIG_CHAINCODE = "confcc"
)

var TLogger = shim.NewLogger(TEST_TARGET + "test")

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit(TEST_TARGET+"Init", args)
	if res.Status != shim.OK {
		TLogger.Info("Init failed", string(res.Message))
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) error {
	bytes := stub.State[name]
	if bytes == nil {
		return fmt.Errorf("failed to get value")
	}
	if string(bytes) != value {
		return fmt.Errorf("not as expected")
	}
	return nil
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string, value *SimpleAsset) error {
	res := stub.MockInvoke(TEST_TARGET+"Query", [][]byte{[]byte("query"), []byte(name)})
	if res.Status != shim.OK {
		return fmt.Errorf(string(res.Message))
	}
	if res.Payload == nil {
		return fmt.Errorf("failed to get value")
	}
	sa := SimpleAsset{}
	err := json.Unmarshal(res.Payload, &sa)
	if err != nil {
		return fmt.Errorf(string(res.Message))
	}
	if !reflect.DeepEqual(sa, *value) {
		return fmt.Errorf("not as expected")
	}
	return nil
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) error {
	res := stub.MockInvoke(TEST_TARGET+"Invoke", args)
	if res.Status != shim.OK {
		return fmt.Errorf(string(res.Message))
	}
	return nil
}

func TestInit(t *testing.T) {
	m := new(SimpleAsset)
	stub := shim.NewMockStub(TEST_TARGET, m)

	TLogger.Info("------ Init test ------")
	checkInit(t, stub, [][]byte{[]byte("init")})
}

func TestSetGetInvoke(t *testing.T) {
	m := new(SimpleAsset)
	stub := shim.NewMockStub(TEST_TARGET, m)

	var setTests = []struct {
		input [][]byte // input
	}{
		{[][]byte{[]byte("set"), []byte("10"), []byte("Test1"), []byte("2000"), []byte("2014-11-12T11:45:26.371Z")}},
		{[][]byte{[]byte("set"), []byte("1000"), []byte("Test2"), []byte("2001"), []byte("2015-11-12T11:45:26.371Z")}},
	}
	TLogger.Info("------ set test ------")
	for _, setTest := range setTests {
		err := checkInvoke(t, stub, setTest.input)
		if err != nil {
			TLogger.Error(err)
		}
	}

	time1, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26.371Z")
	time2, _ := time.Parse(time.RFC3339, "2015-11-12T11:45:26.371Z")

	var getTests = []struct {
		key   string
		value SimpleAsset // input
	}{
		{"Test_10", SimpleAsset{Id: 10, Title: "Test1", Year: "2000", Created: time1}},
		{"Test_1000", SimpleAsset{Id: 1000, Title: "Test2", Year: "2001", Created: time2}},
	}

	for _, getTest := range getTests {
		TLogger.Info("------ get test ------")
		checkQuery(t, stub, getTest.key, &getTest.value)
	}
}
