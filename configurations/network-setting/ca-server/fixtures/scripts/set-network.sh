#!/bin/bash

peer channel create -c m -f /data/m.channel.tx -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE --outputBlock /data/m.block
sleep 5s
peer channel join -b /data/m.block
sleep 5s
peer channel update -c m -f /data/org1.m.anchors.tx -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n confcc -v 1.0 -p github.com/m/chaincode/confcc
sleep 10s
peer chaincode instantiate -C m -n confcc -v 1.0 -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n dadlog -v 1.0 -p github.com/m/chaincode/dadlog
sleep 5s
peer chaincode instantiate -C m -n dadlog -v 1.0 -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n m -v 1.0 -p github.com/m/chaincode/m
sleep 5s
peer chaincode instantiate -C m -n m -v 1.0 -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
export CORE_PEER_TLS_KEY_FILE=/data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls/keystore/server.key
export CORE_PEER_TLS_CERT_FILE=/data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls/signcerts/server.crt
export CORE_PEER_ADDRESS=peer1.org1.hf.m.io:7051
sleep 5s
peer channel join -b /data/m.block
sleep 5s
peer chaincode install -n confcc -v 1.0 -p github.com/m/chaincode/confcc
sleep 10s
peer chaincode invoke -C m -n confcc -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n dadlog -v 1.0 -p github.com/m/chaincode/dadlog
sleep 5s
peer chaincode invoke -C m -n dadlog -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n m -v 1.0 -p github.com/m/chaincode/m
sleep 5s
peer chaincode invoke -C m -n m -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
export CORE_PEER_TLS_KEY_FILE=/data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls/keystore/server.key
export CORE_PEER_TLS_CERT_FILE=/data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls/signcerts/server.crt
export CORE_PEER_ADDRESS=peer2.org1.hf.m.io:7051
sleep 5s
peer channel join -b /data/m.block
sleep 5s
peer chaincode install -n confcc -v 1.0 -p github.com/m/chaincode/confcc
sleep 10s
peer chaincode invoke -C m -n confcc -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n dadlog -v 1.0 -p github.com/m/chaincode/dadlog
sleep 5s
peer chaincode invoke -C m -n dadlog -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
sleep 5s
peer chaincode install -n m -v 1.0 -p github.com/m/chaincode/m
sleep 5s
peer chaincode invoke -C m -n m -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE

