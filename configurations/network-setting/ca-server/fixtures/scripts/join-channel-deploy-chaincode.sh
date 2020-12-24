#!/bin/bash

# Set the msp path of the admin user in org1. Only admin can do following tasks with current setting.
export CORE_PEER_MSPCONFIGPATH=/data/peerOrganizations/org1.hf.m.io/users/admin/msp/

# Create the channel based on channel file created in previous step
peer channel join -b /data/m.block

# Set golang information. To deploy chain code, golang is needed.
#export PATH=$PATH:/usr/local/go/bin 
#export GOPATH=/opt/gopath

cd $PEER_HOME

# Insatll chain code. If the file format in the path for -p option is golang file, it is fine.
peer chaincode install -n m -v 1.0 -p github.com/chaincode/

sleep 10s

# Instantiate chain codei by querying. chain code compile works in this step.
peer chaincode query -C m -n m -c '{"Args":["invoke","query","hello"]}'
