#!/bin/bash

# Initialize the intermediate CA
fabric-ca-server init -b $BOOTSTRAP_USER_PASS -u $PARENT_URL

# Copy the intermediate CA's certificate chain to the data directory to be used by others
cp $FABRIC_CA_SERVER_HOME/ca-chain.pem $TARGET_CHAINFILE

# Add the custom orgs
sed -i "/   org1:/,+2 d" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/   org2:/   hf.m.io:/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/      expiry: 8760h/      expiry: 43800h/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/         expiry: 8760h/         expiry: 43800h/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/      - department1/      - orderer0\n      - orderer1\n      - orderer2\n      - orderer3\n      - orderer4/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/      - C: US/      - C: KR/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/        ST: \"North Carolina\"/        ST: \"Seoul\"/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/        O: Hyperledger/        O: hf.m.io/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sleep 10s

# Start the intermediate CA
fabric-ca-server start

