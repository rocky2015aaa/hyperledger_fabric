#!/bin/bash

cd fixtures
docker-compose -f docker-compose-cas.yaml up -d
sleep 15s
docker exec ica.hf.m.io scripts/create-orderers.sh
docker exec ica.org1.hf.m.io scripts/create-peers.sh
sleep 15s
sudo chown -R madmin:madmin data
scripts/prepare-channel.sh
mv ./data/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts/cert.pem ./data/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts/admin@org1.hf.m.io-cert.pem
mv ./data/peerOrganizations/org1.hf.m.io/users/admin/msp/signcerts/cert.pem ./data/peerOrganizations/org1.hf.m.io/users/admin/msp/signcerts/admin@org1.hf.m.io-cert.pem
mv ./data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts/cert.pem ./data/ordererOrganizations/hf.m.io/users/admin/msp/admincerts/admin@hf.m.io-cert.pem
mv ./data/ordererOrganizations/hf.m.io/users/admin/msp/signcerts/cert.pem ./data/ordererOrganizations/hf.m.io/users/admin/msp/signcerts/admin@hf.m.io-cert.pem
