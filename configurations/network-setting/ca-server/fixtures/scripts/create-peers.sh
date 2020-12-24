#!/bin/bash

# Generating ica.org1.hf.m.io-admin's msp and storing to the location FABRIC_CA_CLIENT_HOME
fabric-ca-client enroll -d -u https://ica.org1.hf.m.io-admin:adminpw@ica.org1.hf.m.io:7054 --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Registering orderers, admin.org1.hf.m.io(admin user in org1) and user in org1
fabric-ca-client register -d --id.name peer0.org1.hf.m.io --id.secret org1pw --id.type peer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name peer1.org1.hf.m.io --id.secret org1pw --id.type peer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name peer2.org1.hf.m.io --id.secret org1pw --id.type peer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name admin.org1.hf.m.io --id.secret org1pw --id.attrs 'hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert' --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Getting CA certs for organization org1 and storing in /data/peerOrganizations/org1.hf.m.io/msp. And making tlscacerts and tlsintermediatecerts
fabric-ca-client getcacert -d -u https://ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/msp

cp -r /data/peerOrganizations/org1.hf.m.io/msp/cacerts /data/peerOrganizations/org1.hf.m.io/msp/tlscacerts
cp -r /data/peerOrganizations/org1.hf.m.io/msp/intermediatecerts /data/peerOrganizations/org1.hf.m.io/msp/tlsintermediatecerts

# Setting the location to store the msp files including certificates of the admin user in org1. In this case, storing in /data, the docker container volume
export FABRIC_CA_CLIENT_HOME=/data/peerOrganizations/org1.hf.m.io/users/admin

# Generating the msp of the admin user in org1 and storing to the location FABRIC_CA_CLIENT_HOME. And making admincerts and org0's admincerts
fabric-ca-client enroll -d -u https://admin.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 --csr.names "C=KR,ST=Seoul,O=hf.m.io"

cp -r $FABRIC_CA_CLIENT_HOME/msp/signcerts $FABRIC_CA_CLIENT_HOME/msp/admincerts
cp -r $FABRIC_CA_CLIENT_HOME/msp/signcerts /data/peerOrganizations/org1.hf.m.io/msp/admincerts

fabric-ca-client register -d --id.name user.org1.hf.m.io --id.secret org1pw --csr.names "C=KR,ST=Seoul,O=hf.m.io"

#export FABRIC_CA_CLIENT_HOME=/data/peerOrganizations/org1.hf.m.io/user.org1.hf.m.io

#fabric-ca-client enroll -d -u https://user.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054

#cp -r /data/org1/admin.org1.hf.m.io/msp/signcerts $FABRIC_CA_CLIENT_HOME/msp/admincerts

# Generate server TLS cert and key pair for the peer : peer 0
fabric-ca-client enroll -d --enrollment.profile tls -u https://peer0.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls --csr.hosts peer0.org1.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls/keystore/*_sk /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls/keystore/server.key
mv /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls/signcerts/cert.pem /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/tls/signcerts/server.crt

# Enroll the peer to get an enrollment certificate and set up the core's local MSP directory
fabric-ca-client enroll -d -u https://peer0.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

cp -r /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/msp/cacerts /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/msp/tlscacerts
cp -r /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/msp/intermediatecerts /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/msp/tlsintermediatecerts
cp -r /data/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts /data/peerOrganizations/org1.hf.m.io/peers/peer0.org1.hf.m.io/msp/admincerts




# Generate server TLS cert and key pair for the peer : peer 1
fabric-ca-client enroll -d --enrollment.profile tls -u https://peer1.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls --csr.hosts peer1.org1.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls/keystore/*_sk /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls/keystore/server.key
mv /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls/signcerts/cert.pem /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/tls/signcerts/server.crt

# Enroll the peer to get an enrollment certificate and set up the core's local MSP directory
fabric-ca-client enroll -d -u https://peer1.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

cp -r /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/msp/cacerts /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/msp/tlscacerts
cp -r /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/msp/intermediatecerts /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/msp/tlsintermediatecerts
cp -r /data/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts /data/peerOrganizations/org1.hf.m.io/peers/peer1.org1.hf.m.io/msp/admincerts





# Generate server TLS cert and key pair for the peer : peer 2
fabric-ca-client enroll -d --enrollment.profile tls -u https://peer2.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls --csr.hosts peer2.org1.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls/keystore/*_sk /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls/keystore/server.key
mv /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls/signcerts/cert.pem /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/tls/signcerts/server.crt

# Enroll the peer to get an enrollment certificate and set up the core's local MSP directory
fabric-ca-client enroll -d -u https://peer2.org1.hf.m.io:org1pw@ica.org1.hf.m.io:7054 -M /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

cp -r /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/msp/cacerts /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/msp/tlscacerts
cp -r /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/msp/intermediatecerts /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/msp/tlsintermediatecerts
cp -r /data/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts /data/peerOrganizations/org1.hf.m.io/peers/peer2.org1.hf.m.io/msp/admincerts

