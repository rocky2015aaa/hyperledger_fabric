#!/bin/bash

# Generating ica.hf.m.io-admin's msp and storing to the location FABRIC_CA_CLIENT_HOME
fabric-ca-client enroll -d -u https://ica.hf.m.io-admin:adminpw@ica.hf.m.io:7054 --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Registering orderers and admin.hf.m.io(admin user in org0)
fabric-ca-client register -d --id.name orderer0.hf.m.io --id.secret org0pw --id.type orderer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name orderer1.hf.m.io --id.secret org0pw --id.type orderer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name orderer2.hf.m.io --id.secret org0pw --id.type orderer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name orderer3.hf.m.io --id.secret org0pw --id.type orderer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name orderer4.hf.m.io --id.secret org0pw --id.type orderer --csr.names "C=KR,ST=Seoul,O=hf.m.io"

fabric-ca-client register -d --id.name admin.hf.m.io --id.secret org0pw --id.attrs 'admin=true:ecert' --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Getting CA certs for organization org0 and storing in /data/ordererOrganizations/hf.m.io/msp. And making tlscacerts and tlsintermediatecerts
fabric-ca-client getcacert -d -u https://ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/msp
cp -r /data/ordererOrganizations/hf.m.io/msp/cacerts /data/ordererOrganizations/hf.m.io/msp/tlscacerts
cp -r /data/ordererOrganizations/hf.m.io/msp/intermediatecerts /data/ordererOrganizations/hf.m.io/msp/tlsintermediatecerts

# Setting the location to store the msp files including certificates of the admin user in org0. In this case, storing in /data, the docker container volume
export FABRIC_CA_CLIENT_HOME=/data/ordererOrganizations/hf.m.io/users/admin

# Generating the msp of the admin user in org0 and storing to the location FABRIC_CA_CLIENT_HOME. And making admincerts and org0's admincerts
fabric-ca-client enroll -d -u https://admin.hf.m.io:org0pw@ica.hf.m.io:7054 --csr.names "C=KR,ST=Seoul,O=hf.m.io"
cp -r $FABRIC_CA_CLIENT_HOME/msp/signcerts $FABRIC_CA_CLIENT_HOME/msp/admincerts
cp -r $FABRIC_CA_CLIENT_HOME/msp/signcerts /data/ordererOrganizations/hf.m.io/msp/admincerts

# Enroll to get orderer's TLS cert (using the "tls" profile) : orderer0
fabric-ca-client enroll -d --enrollment.profile tls -u https://orderer0.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/tls --csr.hosts orderer0.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/tls/keystore/*_sk /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/tls/keystore/server.key
mv /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/tls/signcerts/cert.pem /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/tls/signcerts/server.crt

# Enroll again to get the orderer's enrollment certificate (default profile)
fabric-ca-client enroll -d -u https://orderer0.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Finish setting up the local MSP for the orderer
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/msp/cacerts /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/msp/tlscacerts
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/msp/intermediatecerts /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/msp/tlsintermediatecerts
cp -r /data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts /data/ordererOrganizations/hf.m.io/orderers/orderer0.hf.m.io/msp/admincerts



# Enroll to get orderer's TLS cert (using the "tls" profile) : orderer1
fabric-ca-client enroll -d --enrollment.profile tls -u https://orderer1.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/tls --csr.hosts orderer1.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/tls/keystore/*_sk /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/tls/keystore/server.key
mv /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/tls/signcerts/cert.pem /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/tls/signcerts/server.crt

# Enroll again to get the orderer's enrollment certificate (default profile)
fabric-ca-client enroll -d -u https://orderer1.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Finish setting up the local MSP for the orderer
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/msp/cacerts /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/msp/tlscacerts
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/msp/intermediatecerts /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/msp/tlsintermediatecerts
cp -r /data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts /data/ordererOrganizations/hf.m.io/orderers/orderer1.hf.m.io/msp/admincerts




# Enroll to get orderer's TLS cert (using the "tls" profile) : orderer2
fabric-ca-client enroll -d --enrollment.profile tls -u https://orderer2.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/tls --csr.hosts orderer2.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/tls/keystore/*_sk /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/tls/keystore/server.key
mv /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/tls/signcerts/cert.pem /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/tls/signcerts/server.crt

# Enroll again to get the orderer's enrollment certificate (default profile)
fabric-ca-client enroll -d -u https://orderer2.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Finish setting up the local MSP for the orderer
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/msp/cacerts /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/msp/tlscacerts
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/msp/intermediatecerts /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/msp/tlsintermediatecerts
cp -r /data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts /data/ordererOrganizations/hf.m.io/orderers/orderer2.hf.m.io/msp/admincerts






# Enroll to get orderer's TLS cert (using the "tls" profile) : orderer3
fabric-ca-client enroll -d --enrollment.profile tls -u https://orderer3.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/tls --csr.hosts orderer3.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/tls/keystore/*_sk /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/tls/keystore/server.key
mv /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/tls/signcerts/cert.pem /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/tls/signcerts/server.crt

# Enroll again to get the orderer's enrollment certificate (default profile)
fabric-ca-client enroll -d -u https://orderer3.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Finish setting up the local MSP for the orderer
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/msp/cacerts /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/msp/tlscacerts
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/msp/intermediatecerts /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/msp/tlsintermediatecerts
cp -r /data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts /data/ordererOrganizations/hf.m.io/orderers/orderer3.hf.m.io/msp/admincerts







# Enroll to get orderer's TLS cert (using the "tls" profile) : orderer4
fabric-ca-client enroll -d --enrollment.profile tls -u https://orderer4.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/tls --csr.hosts orderer4.hf.m.io --csr.names "C=KR,ST=Seoul,O=hf.m.io"

mv /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/tls/keystore/*_sk /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/tls/keystore/server.key
mv /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/tls/signcerts/cert.pem /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/tls/signcerts/server.crt

# Enroll again to get the orderer's enrollment certificate (default profile)
fabric-ca-client enroll -d -u https://orderer4.hf.m.io:org0pw@ica.hf.m.io:7054 -M /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/msp --csr.names "C=KR,ST=Seoul,O=hf.m.io"

# Finish setting up the local MSP for the orderer
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/msp/cacerts /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/msp/tlscacerts
cp -r /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/msp/intermediatecerts /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/msp/tlsintermediatecerts
cp -r /data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts /data/ordererOrganizations/hf.m.io/orderers/orderer4.hf.m.io/msp/admincerts

