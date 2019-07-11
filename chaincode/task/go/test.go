package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"strings"
	
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type product struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Quantity   int    `json:"quantity"`
	Price      int    `json:"price"`
	Owner      string `json:"owner"`
}

// type supplyChain2 struct {
// 	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
// 	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
// 	Quantity   int    `json:"quantity"`
// 	Price      int    `json:"price"`
// 	Owner      string `json:"owner"`
// }

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "initProduct1":
		//create a new product
		return t.initProduct1(stub, args)
	case "initProduct2":
		//create a new product
		return t.initProduct1(stub, args)		
	case "readProduct1":
		//read a product
		return t.readProduct1(stub, args)
	case "readProduct2":
		//read a product
		return t.readProduct2(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

// ============================================================
// initproduct - create a new product for supplyChain1, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initProduct1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type productTransientInput struct {
		Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
		Quantity   int    `json:"quantity"`
		Price      int    `json:"price"`
		Owner      string `json:"owner"`
	}

	// ==== Input sanitation ====
	fmt.Println("- start init product1")

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private product data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}
    fmt.Println("the transient is: ", transMap);
	if _, ok := transMap["product"]; !ok {
		return shim.Error("product must be a key in the transient map")
	}
	
	fmt.Println("the transmap is: ", transMap["product"]);
	fmt.Println("the transmap is: ", len(transMap["product"]));
	
	if len(transMap["product"]) == 0 {
		return shim.Error("product value in the transient map must be a non-empty JSON string")
	}

	var productInput productTransientInput
	err = json.Unmarshal(transMap["product"], &productInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["product"]))
	}

	if len(productInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if productInput.Quantity <= 0 {
		return shim.Error("quantity field must be a positive integer")
	}
	if len(productInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if productInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}

	//Get the client ID object
	id, err := cid.New(stub)
	if err != nil {
		return shim.Error("unable to create a client ID object")
	}
	mspid, err := id.GetMSPID()
	if err != nil {
		return shim.Error("unable to get msp id")
	}
	 fmt.Println("mspid is: " + mspid)

	if mspid == "ManufacturerOrgMSP" {
	// ==== Check if product already exists ====
	productAsBytes, err := stub.GetPrivateData("supplyChain1", productInput.Name)
	if err != nil {
		return shim.Error("Failed to get product: " + err.Error())
	} else if productAsBytes != nil {
		fmt.Println("This product already exists: " + productInput.Name)
		return shim.Error("This product already exists: " + productInput.Name)
	}
	fmt.Println("Inside to add to database" )
	// ==== Create product object, marshal to JSON, and save to state ====
	product := &product{
		ObjectType: "product",
		Name:       productInput.Name,
		Quantity:   productInput.Quantity,
		Price:      productInput.Price,
		Owner:      productInput.Owner,
	}
	productJSONasBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutPrivateData("supplyChain1", productInput.Name, productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	 }
     }else {shim.Error("Only Manufacturer can write")}
	// ==== product saved and indexed. Return success ====
	fmt.Printf("- end init product")
	return shim.Success(nil)
}

// ============================================================
// initproduct - create a new product for supplyChain2, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initProduct2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type productTransientInput struct {
		Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
		Quantity   int    `json:"quantity"`
		Price      int    `json:"price"`
		Owner      string `json:"owner"`
	}

	// ==== Input sanitation ====
	fmt.Println("- start init product")

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private product data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	if _, ok := transMap["product"]; !ok {
		return shim.Error("product must be a key in the transient map")
	}

	if len(transMap["product"]) == 0 {
		return shim.Error("product value in the transient map must be a non-empty JSON string")
	}

	var productInput productTransientInput
	err = json.Unmarshal(transMap["product"], &productInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["product"]))
	}

	if len(productInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if productInput.Quantity <= 0 {
		return shim.Error("quantity field must be a positive integer")
	}
	if len(productInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if productInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}
	//Get the client ID object
	id, err := cid.New(stub)
	if err != nil {
		return shim.Error("unable to create a client ID object")
	}
	mspid, err := id.GetMSPID()
	if err != nil {
		return shim.Error("unable to get msp id")
	}

	if mspid == "org1MSP" {
	// ==== Check if product already exists ====
	productAsBytes, err := stub.GetPrivateData("supplyChain2", productInput.Name)
	if err != nil {
		return shim.Error("Failed to get product: " + err.Error())
	} else if productAsBytes != nil {
		fmt.Println("This product already exists: " + productInput.Name)
		return shim.Error("This product already exists: " + productInput.Name)
	}

	// ==== Create product object, marshal to JSON, and save to state ====
	product := &product{
		ObjectType: "product",
		Name:       productInput.Name,
		Quantity:   productInput.Quantity,
		Price:      productInput.Price,
		Owner:      productInput.Owner,
	}
	productJSONasBytes, err := json.Marshal(product)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	err = stub.PutPrivateData("supplyChain2", productInput.Name, productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
    }else {shim.Error("Only Manufacturer can write")}

	// ==== product saved and indexed. Return success ====
	fmt.Printf("- end init product")
	return shim.Success(nil)
}

// ===============================================
// readProduct1 - read a product from chaincode state
// ===============================================
func (t *SimpleChaincode) readProduct1(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error
    
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the product to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("supplyChain1", name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Product does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

// ===============================================
// readProduct2 - read a product from chaincode state
// ===============================================
func (t *SimpleChaincode) readProduct2(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the product to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("supplyChain2", name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Product does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}