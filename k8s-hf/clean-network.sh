#!/bin/bash

kubectl delete -n hfbn statefulset.apps/peer0 statefulset.apps/peer1 statefulset.apps/peer2
kubectl delete -n hfbn service/peer0-svc service/peer1-svc service/peer2-svc
kubectl delete -n hfbn statefulset.apps/orderer0 statefulset.apps/orderer1 statefulset.apps/orderer2 statefulset.apps/orderer3 statefulset.apps/orderer4
kubectl delete -n hfbn service/orderer0-svc service/orderer1-svc service/orderer2-svc service/orderer3-svc service/orderer4-svc
kubectl delete -n hfbn statefulset.apps/ica service/ica-svc statefulset.apps/rca service/rca-svc

cd artifacts
sudo rm -rf ca-cert.pem m.channel.tx orderer.genesis.block ordererOrganizations org0-ca-chain.pem org1-ca-chain.pem org1.m.anchors.tx peerOrganizations m.block
cd ..
sudo rm -rf block-data/orderers/orderer0/orderer
sudo rm -rf block-data/orderers/orderer1/orderer
sudo rm -rf block-data/orderers/orderer2/orderer
sudo rm -rf block-data/orderers/orderer3/orderer
sudo rm -rf block-data/orderers/orderer4/orderer
sudo rm -rf block-data/peers/peer0/production
sudo rm -rf block-data/peers/peer1/production
sudo rm -rf block-data/peers/peer2/production
docker rm -f $(docker ps -qa --filter "name=dev-peer")
