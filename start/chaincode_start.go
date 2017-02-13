package main

import (
	"errors"
	"fmt"

	"encoding/json"



	"github.com/hyperledger/fabric/core/chaincode/shim"
)


type SimpleChaincode struct {
}


var openOrdersStr = "_openorders"


type MilkContainer struct{

        ContainerID string `json:"containerid"`
        User string        `json:"user"`

        Litres string        `json:"litres"`

}

type SupplyCoin struct{

        CoinID string `json:"coinid"`
        User string        `json:"user"`
}

type Order struct{
        OrderID string `json:"orderid"`
       User string `json:"user"`
       Status string `json:"status"`
       Litres string    `json:"litres"`
}



type AllOrders struct{
	OpenOrders []Order `json:"open_orders"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
func(t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
                                                          
}

       err = stub.PutState("hello world",[]byte("welcome ti supply chain management"))  //Just to check the network 
       if err != nil {
		return nil, err
}
/*
      var empty []string
	jsonAsBytes, _ := json.Marshal(empty)                          //marshal an emtpy array of strings to clear the index
	err = stub.PutState(milkIndexStr, jsonAsBytes)                 //Making milk container list as empty - resetting
	if err != nil {
		return nil, err
} 
        err = stub.PutState(coinIndexStr, jsonAsBytes)                 //Making coin list as empty
        if err != nil {
                return nil, err
}
*/
	
	
	var orders AllOrders
	jsonAsBytes, _ = json.Marshal(orders)								//clear the open trade struct
	err = stub.PutState(openOrdersStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
        return nil, nil

}



func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
        }else if function == "Create_milkcontainer" {		//creates a milk container-invoked by supplier   
		return t.Create_milkcontainer(stub, args)      

	}else if function == "Create_coinmarket" {		//creates a coin - invoked by market
		return t.Create_coinmarket(stub, args)	

        }else if function == "Create_coinlogistics"{              //creates a coin - invoked by logistics 
                return t.Create_coinmarket(stub, args)
	} else if function == "Order_milk"{
		return t.Order_milk(stub,args)
	}


     return nil,nil

}


func (t *SimpleChaincode) Create_milkcontainer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
var err error

// "1x22" "supplier" 20 
// args[0] args[1] args[2] 

id := args[0]
user := args[1]
litres :=args[2] 
milkAsBytes, err := stub.GetState(id) 
if err != nil {
		return nil, errors.New("Failed to get details og given id") 
}

res := MilkContainer{} 
json.Unmarshal(milkAsBytes, &res)

if res.ContainerID == id{

        fmt.Println("Container already exixts")
        fmt.Println(res)
        return nil,errors.New("This cpontainer alreadt exists")
}

res.ContainerID = id
res.User = user
res.Litres = litres
milkAsBytes, _ =json.Marshal(res)

stub.PutState(id,milkAsBytes)

return nil,nil

}


func (t *SimpleChaincode) Create_coinmarket(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

//"1x245" "Market"
id := args[0]
user:= args[1]

coinAsBytes , err := stub.GetState(id)
if err != nil{
              return nil, errors.New("Failed to get details of given id")
} 

res :=SupplyCoin{}

json.Unmarshal(coinAsBytes, &res)

if res.CoinID == id{

          fmt.Println("Coin already exists")
          fmt.Println(res)
          return nil,errors.New("This coin already exists")
}

res.CoinID = id
res.User = user

coinAsBytes, _ = json.Marshal(res)
stub.PutState(id,coinAsBytes)
return nil,nil
}
 

func(t *SimpleChaincode) Create_coinlogistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


// "1x226" "Logistics"

id := args[0]
user:= args[1]

coinAsBytes , err := stub.GetState(id)
if err != nil{
              return nil, errors.New("Failed to get details of given id")
}

res :=SupplyCoin{}

json.Unmarshal(coinAsBytes, &res)

if res.CoinID == id{

          fmt.Println("Coin already exists")
          fmt.Println(res)
          return nil,errors.New("This coin already exists")
}

res.CoinID = id
res.User = user

coinAsBytes, _ = json.Marshal(res)
stub.PutState(id,coinAsBytes)
return nil,nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

if function == "read" {						//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}


func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the variable to query")
	}


	name = args[0]
	valAsbytes, err := stub.GetState(name)				//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for  \"}"
		return nil, errors.New(jsonResp)
	}

	



return valAsbytes, nil										       //send it onward
}

func (t *SimpleChaincode) Order_milk(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
Openorder := Order{}
Openorder.User = "Market"
Openorder.Status = "pending"
Openorder.OrderID = "abcd"
Openorder.Litres = args[0]
orderasbytes,_ := json.Marshal(Openorder)
	var err error
err = stub.PutState("abcd",orderasbytes)
	
	if err != nil {
		return nil, err
        }

	
	
	
	                                                      // writing the order to orders list
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)										//un stringify it aka JSON.parse()
	
	orders.OpenOrders = append(trades.OpenOrders, Openorder);						//append to open trades
	fmt.Println("! appended the order to orders list so that supplier can view")
	jsonAsBytes, _ = json.Marshal(orders)
	err = stub.PutState(openOrdersStr, jsonAsBytes)								//rewrite open orders
	if err != nil {
		return nil, err
	}
	
	t.read(stub,Openorder.OrderID)             //printing the Order

return nil,nil
}
func (t *SimpleChaincode) view_orders(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	
/*orderasbytes, _ := stub.GetState(args[0])
Performorder := Order{}
json.Unmarshal(orderasbytes, &Performorder)
containerasbytes, _ := stub.GetState("1x23")
Container := MilkContainer{}
json.Unmarshal(containerasbytes, &Container)
if(Container.Litres == Performorder.Litres){
fmt.Println("Hurray, we got want u want")
Performorder.Status="received"
orderasbytes,_=json.Marshal(Performorder)
stub.PutState(args[0],orderasbytes)
var a []string
a[0] = args[0]
a[1] ="1x23" 
t.init_logistics(stub,a)
return nil,nil
} else{
fmt.Println("Sorry")
return nil,nil



}

*/
	ordersAsBytes, err := stub.GetState(openOrdersStr)
	if err != nil {
		return nil, errors.New("Failed to get opentrades")
	}
	var orders AllOrders
	json.Unmarshal(ordersAsBytes, &orders)										//un stringify it aka JSON.parse()
	
for i := range orders.OpenOrders{
	
	milkAsBytes, err := stub.GetState(args[0])
			if err != nil {
				return nil, errors.New("Failed to get thing")
			}
			saleMilk := MilkContainer{}
			json.Unmarshal(milkAsBytes, &saleMilk)	
	
	if orders.OpenOrders[i].Litres == saleMilk.Litres{
		
		fmt.Println("Shipping will be done soon, u are inside view_order function")
		orders.OpenOrders[i].Status = "Received order"
	}
}
	
	return nil,nil
	
	
	
}
func (t *SimpleChaincode) init_logistics(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
orderid := args[0]
fmt.Println("Recevied order for shipping is %s ",orderid)
orderasbytes, _ := stub.GetState(orderid)
Shiporder := Order{}
json.Unmarshal(orderasbytes, &Shiporder)
Shiporder.Status = "Shipped and in transit"
orderasbytes,_ = json.Marshal(Shiporder)
stub.PutState(orderid, orderasbytes)
var a []string
a[0] = args[0]
a[1] ="1x23"
t.completedelivery(stub,a)
return nil,nil
	

	
	
	
}
func (t *SimpleChaincode) completedelivery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
orderid := args[0]
orderasbytes, _ := stub.GetState(orderid)
Shiporder := Order{}
json.Unmarshal(orderasbytes, &Shiporder)
milkasbytes, _ := stub.GetState(args[1])
milkcont := MilkContainer{}
json.Unmarshal(milkasbytes, &milkcont)
milkcont.User="Market"
Shiporder.Status = "Delivered"
orderasbytes,_ = json.Marshal(Shiporder)
stub.PutState(orderid, orderasbytes)
milkasbytes,_ = json.Marshal(milkcont)
stub.PutState(args[1],milkasbytes)
var a []string
a[0] = args[0]
a[1] =args[1]
t.checkproduct(stub,a)
return nil,nil
}
func (t *SimpleChaincode) checkproduct(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
milkcontainerid := args[1]
milkasbytes, _ := stub.GetState(milkcontainerid)
milkcontainer := MilkContainer{}
json.Unmarshal(milkasbytes,&milkcontainer)
if(milkcontainer.User == "Market"){
var a []string
a[0] = "1x245"
       t.init_cointransfer(stub,a)
       return nil,nil
}else{
      return nil,errors.New("Couldn't transfer,please try again")
}
return nil,nil
}
func (t *SimpleChaincode) init_cointransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
coinid := args[0]
coinasbytes, _ := stub.GetState(coinid)
Finalcoin:= SupplyCoin{}
json.Unmarshal(coinasbytes,&Finalcoin)
Finalcoin.User = "Supplier"
coinasbytes,_ = json.Marshal(Finalcoin)
stub.PutState(coinid, coinasbytes)
return nil,nil
}
