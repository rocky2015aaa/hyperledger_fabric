#!/bin/bash

# Initialize the root CA
fabric-ca-server init -b $BOOTSTRAP_USER_PASS

# Copy the root CA's signing certificate to the data directory to be used by others
cp $FABRIC_CA_SERVER_HOME/ca-cert.pem $TARGET_CERTFILE

# Add the custom orgs
sed -i "/   org1:/,+2 d" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/   org2:/   org0: []\n   org1:/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/      - C: US/      - C: KR/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/        ST: \"North Carolina\"/        ST: \"Seoul\"/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml
sed -i "s/        O: Hyperledger/        O: m.io/" $FABRIC_CA_SERVER_HOME/fabric-ca-server-config.yaml


# Start the root CA
fabric-ca-server start

