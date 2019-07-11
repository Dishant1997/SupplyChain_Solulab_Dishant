/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const base64 = require('base-64');
const utf8 = require('utf8');

const ccpPath = path.resolve(__dirname, '..', '..', 'task-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const buffer = require('buffer').Buffer
async function main() {
    try {

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        console.log(walletPath);
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        console.log('hi2');
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
        console.log('hi1');
        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');
        console.log('hi3');
        // Get the contract from the network.
        const contract = network.getContract('task');

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR10', 'Dave')

        // Private data sent as transient data: { [key: string]: Buffer 

        var transientInput = {
            product:  buffer.from("{\"name\":\"Sofa\",\"quantity\":8,\"price\":22000,\"owner\":\"steve\"}")
            // quantity: '4',
            // owner:    'Tom',
            // price:    '35000'
        };
        let product = Buffer.from('\"name\":\"TV\",\"quantity\":4,\"price\":35000,\"owner\":\"tom\"','base64');
        console.log('Sending tx with transient', JSON.stringify(product));
        //transientData.productname=JSON.stringify(transientData.productname);
        //transientData=JSON.stringify(transientData)
        const result = await contract.createTransaction('initProduct1')
            .setTransient(transientInput)
            .submit();

        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.log("inside catch")
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
