**Hyperledger Fabric Implementation of Private Data in Supply Chain Management :**

This project showcases the implementation of private data in the scenario of supply chain management where distributor 1, shipper 1 and retail 1 cannot see the data of   distributor 2, shipper 2 and retail 2 and vice versa.
Other than that only the manufacturer is allowed to upload data for a product while all the other participants have the read rights only, which is implemented using cid library(access control).

**Structure adapted for the project:**

    1. Bin: It consist of the platform-specific binaries you will need to set up your network.
    2. Chaincode: It consist of the smart contract written for our network.(Golang is used in this case)
    3. Task-network: It consist of files responsible for bringing up the network, building channel artifacts and crypto-config for the participants. The ordering here is solo.
    4. Task-network-kafka: It consist of the implementation of network using kafka as the ordering service.
    5. Taskapp- This folder consist of application level implementation. API created in nodejs SDK responsible for invoking the chaincodes to update the ledger.

**Pre-requisites to build this project:**

    • cURL
    • Docker and Docker Compose
    • Go Programming Language
    • Node.js and NPM
    • Python
    • Binaries and docker images for fabric
    • fabric version:1.4

**Step by step build the project on your PC:**

    • Clone the github repository
        ◦ git clone https://github.com/Dishant1997/Task_Solulab_Dishant.git

    • Create crypto config for our network
        ◦ ../bin/cryptogen generate –config=./crypto-config.yaml

    • Create channel artifacts for genesis block
        ◦ ../bin/configtxgen -profile FourOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block

    • Create channel.tx for channel mychannel
        ◦ export CHANNEL_NAME=mychannel  && ../bin/configtxgen -profile FourOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

    • Channel artifacts for each organizations
        ◦ ../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/ManufacturerOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ManufacturerOrgMSP
        ◦ ../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/RetailOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg RetailOrgMSP
        ◦ ../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/DistributorOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg DistributorOrgMSP
        ◦ ../bin/configtxgen -profile FourOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/ShipOrgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg ShipOrgMSP

    • Up all the containers using docker compose
        ◦ docker-compose -f docker-compose-cli.yaml up -d

    • Bash into the Cli container
        ◦ docker exec -it cli bash

    • Set the environment variable to Manufacturer Org
        ◦ CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacturerorg.example.com/users/Admin@manufacturerorg.example.com/msp
        ◦ CORE_PEER_ADDRESS=peer0.manufacturerorg.example.com:7051
        ◦ CORE_PEER_LOCALMSPID="ManufacturerOrgMSP"
        ◦ CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacturerorg.example.com/peers/peer0.manufacturerorg.example.com/tls/ca.crt
        ◦ export CHANNEL_NAME=mychannel

    • Connect this peer to the channel
        ◦ peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx 

    • Connect the channel to the mychannel.block
        ◦ peer channel join -b mychannel.block

    • Follow the above three steps for each organizations

    • Update the anchor peers
        ◦ peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/ManufacturerOrgMSPanchors.tx 

    • Install the chaincodeon the peers
        ◦ peer chaincode install -n task -v 1.0 -p github.com/chaincode/task/go/

    • Now go the taskapp/javascript
        ◦ npm install
        ◦ node enrollAdmin.js
        ◦ node registerUser.js
        ◦ node invoke.js
        ◦ node query.js
