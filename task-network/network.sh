../bin/cryptogen generate --config=./crypto-config.yaml

export FABRIC_CFG_PATH=$PWD

../bin/configtxgen -profile FourOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block

export CHANNEL_NAME=mychannel  && ../bin/configtxgen -profile FourOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/ManufacturerOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ManufacturerOrgMSP
../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/RetailOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg RetailOrgMSP
../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/DistributorOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg DistributorOrgMSP
../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/ShipOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ShipOrgMSP

docker-compose -f docker-compose-cli.yaml up -d

docker exec -it cli bash

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacturerorg.example.com/users/Admin@manufacturerorg.example.com/msp
CORE_PEER_ADDRESS=peer0.manufacturerorg.example.com:7051
CORE_PEER_LOCALMSPID="ManufacturerOrgMSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacturerorg.example.com/peers/peer0.manufacturerorg.example.com/tls/ca.crt
export CHANNEL_NAME=mychannel

peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx 
peer channel join -b mychannel.block

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/retailorg.example.com/users/Admin@retailorg.example.com/msp
CORE_PEER_ADDRESS=peer0.retailorg.example.com:9051
CORE_PEER_LOCALMSPID="RetailOrgMSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/retailorg.example.com/peers/peer0.retailorg.example.com/tls/ca.crt

peer channel join -b mychannel.block

peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/ManufacturerOrgMSPanchors.tx 

peer chaincode install -n task -v 1.0 -p github.com/chaincode/task/go/

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n task -v 1.0 -c '{"Args":["init"]}' -P "OR('ManufacturerOrgMSP.member')" --collections-config  $GOPATH/src/github.com/chaincode/task/task_config.json 

export PRODUCT=$(echo -n "{\"name\":\"TV\",\"quantity\":4,\"price\":35000,\"owner\":\"tom\"}" | base64 | tr -d \\n)
peer chaincode invoke -o orderer.example.com:7050 -C mychannel -n task -c '{"Args":["initProduct1"]}'  --transient "{\"product\":\"$PRODUCT\"}"

peer chaincode query -C mychannel -n task -c '{"Args":["readProduct1","TV"]}'