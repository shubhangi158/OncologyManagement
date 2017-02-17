//Go programs run in packages. main package is the starting point to run the program.
package main

// Package shim provides APIs for the chaincode to access its state variables, transaction context and call other chaincodes.
import(
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"strconv"
)

//NewLogger allows a Go language chaincode to create one or more logging objects 
//whose logs will be formatted consistently with, 
//and temporally interleaved with the logs created by the shim interface. 
//The logs created by this object can be distinguished from shim logs 
//by the name provided, which will appear in the logs.
var logger = shim.NewLogger("PatientChaincode");

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct{
}

//==============================================================================================================================
//	Patient - Defines the structure for a Patient object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON age -> Struct age.
//==============================================================================================================================
type Patient struct{
	PatientId string `json:"patientId"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Illness   string `json:"illness"`
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main(){
	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Printf("Error in starting chaincode : %s", err)
	}
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Patient - Creates the initial JSON for the patient and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode)create_patient(stub shim.ChaincodeStubInterface, patientId string, age int, gender string, illness string) ([]byte, error) {
	
	var p Patient
	
	patient_json := []byte(`{"patientId": "` + patientId + `", "age": ` + strconv.Itoa(age) + `, "gender": "` + gender +  `", "illness": "` + illness + `"}`)						//Build the Patient JSON object
	
	if patientId == "" {
		fmt.Printf("CREATE_PATIENT: Invalid Patient ID provided.")
		return nil,errors.New("Invalid Patient ID provided.")
	}
	
	err := json.Unmarshal(patient_json, &p)													//Converts the JSON defined above into a Patient object 
	if err != nil{
		fmt.Printf("Invalid JSON Object")
		return nil,errors.New("Invalid JSON Object")
	}

	/*record, err := stub.GetState(p.PatientId)														
	if record != nil{																				//If not an error then a record exists, so cant create a new  with this V5cID 
		return nil, errors.New("Patient already exists.")
	}*/
	
	_, err = t.save_changes(stub, p)
	if err != nil { 
		fmt.Printf("CREATE_PATIENT: Error saving changes: %s", err); 
		return nil, errors.New("Error saving changes") 
	}
	
	return nil, nil
	
}

//==============================================================================================================================
// save_changes - Writes to the ledger the Patient struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, p Patient) (bool, error){

	patientJson, err := json.Marshal(p)
	if err != nil{
		fmt.Printf("SAVE_CHANGES: Error converting Patient record: %s", err)
		return false, errors.New("Error converting Patient record")
	}
	
	err = stub.PutState(p.PatientId, patientJson)
	if err != nil{
		fmt.Printf("SAVE_CHANGES: Error storing Patient record: %s", err)
		return false, errors.New("Error storing Patient record")
	}
	
	return true, nil
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	fmt.Printf("Inside Init ...")
	return nil, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){

	if len(args) != 4 {
		return nil, errors.New("INVOKE: Incorrect number of parameters passed. Expecting 4.")
	}

	age, err := strconv.Atoi(args[1])
	if err != nil{
		return nil, errors.New("INVOKE: String to int conversion failed")
	}
	
	if function == "create_patient"{
		return t.create_patient(stub, args[0], age, args[2], args[3])
	}                                                                                                                                            
	
	return nil, errors.New("Function " + function + "doesn't exist")
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){

	logger.Debug("Function: ", function)
	
	if len(args) != 1{
		fmt.Printf("Incorrect number of arguments passed")
		return nil, errors.New("QUERY: Incorrect number of arguments passed")
	}
	
	if function == "getPatientRecord" {
		p, err := t.getPatient(stub, args[0])
		if err != nil{
			fmt.Printf("QUERY: Error retrieving Patient : %s", err)
			return nil, errors.New("QUERY: Error retrieving Patient " + err.Error())
		}
		return t.getPatientRecord(stub, p)
	}
	
	return nil, errors.New("Received unknown function invocation " + function)
	
}

//==============================================================================================================================
//	 getPatient - Gets the state of the data for patientId in the ledger then converts it from the stored
//				  JSON into the Patient struct for use in the contract. Returns the Patient struct.
//				  Returns empty p if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) getPatient(stub shim.ChaincodeStubInterface, patientId string) (Patient, error){

	var p Patient
	bytes, err := stub.GetState(patientId)
	if err != nil {	
		fmt.Printf("GET_PATIENT: Failed to get Patient record: %s", err) 
		return p, errors.New("GET_PATIENT: Error retrieving patient record with patientId = " + patientId) 
	}

	err = json.Unmarshal(bytes, &p)
	if err != nil {	
		fmt.Printf("GET_PATIENT: Corrupt patient record " + string(bytes)+": %s", err) 
		return p, errors.New("GET_PATIENT: Corrupt patient record"+string(bytes))	
	}

	return p, nil
}

func (t *SimpleChaincode) getPatientRecord(stub shim.ChaincodeStubInterface, p Patient) ([]byte, error){
	bytes, err := json.Marshal(p)
	if err != nil {
		return bytes, errors.New("GET_PATIENT_RECORD: Invalid Patient object.")
	}
	return bytes, nil
}



//To-Do List

//1. Implement struct to store all patient IDs.