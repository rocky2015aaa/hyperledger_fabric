# create genesis blocka and channel and set anchor peer

export FABRIC_CFG_PATH=/home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures
/home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/bin/configtxgen -profile M -outputBlock /home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/orderer.genesis.block -channelID chain

/home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/bin/configtxgen -profile M -outputCreateChannelTx /home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/m.channel.tx -channelID m

/home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/bin/configtxgen -profile M -outputAnchorPeersUpdate /home/m/workspace/go/src/github.com/m/configurations/network-setting/ca-server/fixtures/data/org1.m.anchors.tx -channelID m -asOrg Org1M

