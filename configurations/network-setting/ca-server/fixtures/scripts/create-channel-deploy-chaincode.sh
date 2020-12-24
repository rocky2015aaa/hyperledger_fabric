#!/bin/bash

# Set the msp path of the admin user in org1. Only admin can do following tasks with current setting.
export CORE_PEER_MSPCONFIGPATH=/data/peerOrganizations/org1.hf.m.io/users/admin/msp/

# Create the channel based on channel file created in previous step
peer channel create -c m -f /data/m.channel.tx -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE --outputBlock /data/m.block

# Join the channle
peer channel join -b /data/m.block

# Set anchor peer
peer channel update -c m -f /data/peerOrganizations/org1.hf.m.io/org1.m.anchors.tx -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE

# Set golang information. To deploy chain code, golang is needed.
#export PATH=$PATH:/usr/local/go/bin 
#export GOPATH=/opt/gopath

cd $PEER_HOME

# Insatll chain code. If the file format in the path for -p option is golang file, it is fine. 
peer chaincode install -n m -v 1.0 -p github.com/chaincode/

sleep 10s

# Instantiate chain code. chain code compile works in this step.
peer chaincode instantiate -C m -n m -v 1.0 -c '{"Args":["init"]}' -o orderer0.hf.m.io:7050 --tls --cafile /data/org0-ca-chain.pem --clientauth --keyfile $CORE_PEER_TLS_KEY_FILE --certfile $CORE_PEER_TLS_CERT_FILE
