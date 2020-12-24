#!/bin/bash

kubectl apply -f hfbn-namespace-volume.yaml
sleep 10s
kubectl apply -f common-config.yaml
kubectl apply -f rca-deployment.yaml
sleep 10s
kubectl apply -f ica-deployment.yaml
sleep 30s
kubectl exec -n hfbn statefulset.apps/ica -c orderer-ica /data/scripts/create-orderers.sh
sleep 10s
kubectl exec -n hfbn statefulset.apps/ica -c peer-ica /data/scripts/create-peers.sh

sudo chown -R madmin:madmin artifacts/
mv artifacts/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts/cert.pem artifacts/peerOrganizations/org1.hf.m.io/users/admin/msp/admincerts/admin@org1.hf.m.io-cert.pem
mv artifacts/peerOrganizations/org1.hf.m.io/users/admin/msp/signcerts/cert.pem artifacts/peerOrganizations/org1.hf.m.io/users/admin/msp/signcerts/admin@org1.hf.m.io-cert.pem
mv artifacts/ordererOrganizations/hf.m.io/users/admin/msp/admincerts/cert.pem artifacts/ordererOrganizations/hf.m.io/users/admin/msp/admincerts/admin@hf.m.io-cert.pem
mv artifacts/ordererOrganizations/hf.m.io/users/admin/msp/signcerts/cert.pem artifacts/ordererOrganizations/hf.m.io/users/admin/msp/signcerts/admin@hf.m.io-cert.pem

export FABRIC_CFG_PATH=/home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/scripts/
/home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/scripts/bin/configtxgen -profile MyCreditChain -outputBlock /home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/orderer.genesis.block -channelID chain
/home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/scripts/bin/configtxgen -profile MyCreditChain -outputCreateChannelTx /home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/m.channel.tx -channelID m
/home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/scripts/bin/configtxgen -profile MyCreditChain -outputAnchorPeersUpdate /home/madmin/workspace/go/src/github.com/m/k8s-hf/artifacts/org1.m.anchors.tx -channelID m -asOrg Org1MyCreditChain

sudo chown -R madmin:madmin block-data/
kubectl apply -f hfbn-orderer-svcs.yaml
kubectl apply -f hfbn-orderers.yaml
kubectl apply -f hfbn-peer-svcs.yaml
kubectl apply -f hfbn-peers.yaml
sleep 60s
sudo chown -R madmin:madmin artifacts/
sudo chown -R madmin:madmin block-data/
#kubectl exec -n hfbn statefulset.apps/peer0 -c peer-cli /data/scripts/set-network.sh
