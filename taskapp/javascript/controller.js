//SPDX-License-Identifier: Apache-2.0

/*
  This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/fabcar/query.js
  and https://github.com/hyperledger/fabric-samples/blob/release/fabcar/invoke.js
 */
//packages that we need
'use strict';
// var express = require('express');  
// var app = express();   
// app.use(express.json());
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const base64 = require('base-64');
const utf8 = require('utf8');

const ccpPath = path.resolve(__dirname, '..', '..', 'task-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
const buffer = require('buffer').Buffer

module.exports = (function() {
return{
	invoke: async function(req, res){
		try {
			
			//var array = req.params.split(",");
			console.log(req.body);

			var name1 = req.body.name;
			var quantity1 = req.body.quantity;
			var price1 = req.body.price;
			var owner1 = req.body.owner;
			//console.log(req.body.name);
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
			const myquery=console.log('\\"name\\":\\"', name1 ,'\\",\\"quantity\\":' ,quantity1, ',\\"price\\":' ,price1,',\\"owner\\":\\"' ,owner1,'\\"');
			console.log("hi", myquery);
			var transientInput = {
				product:  buffer.from("{'",myquery,"'}")
				// quantity: '4',
				// owner:    'Tom',
				// price:    '35000'
			};
			let product = Buffer.from("{'\"name\":\", name1 ,\",\"quantity\":4,\"price\":35000,\"owner\":\"tom\"','base64'}");
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
	},
	query: async function(req, res){
		try {
			 var key = req.params.product
			 console.log(key);
			// Create a new file system based wallet for managing identities.
			const walletPath = path.join(process.cwd(), 'wallet');
			const wallet = new FileSystemWallet(walletPath);
			console.log(`Wallet path: ${walletPath}`);
	
			// Check to see if we've already enrolled the user.
			const userExists = await wallet.exists('user1');
			if (!userExists) {
				console.log('An identity for the user "user1" does not exist in the wallet');
				console.log('Run the registerUser.js application before retrying');
				return;
			}
	
			// Create a new gateway for connecting to our peer node.
			const gateway = new Gateway();
			await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });
	
			// Get the network (channel) our contract is deployed to.
			const network = await gateway.getNetwork('mychannel');
	
			// Get the contract from the network.
			const contract = network.getContract('task');
	
			// Evaluate the specified transaction.
			// queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
			// queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
			const result = await contract.evaluateTransaction('readProduct1', key);
			console.log(`Transaction has been evaluated, result is:${result.toString()}`);
	
		} catch (error) {
			console.error(`Failed to evaluate transaction: ${error}`);
			process.exit(1);
		}
	}

}
})();