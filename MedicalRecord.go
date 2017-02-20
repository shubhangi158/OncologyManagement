package main

import(
	//"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"strconv"
)

var logger = shim.NewLogger("PatientRecord")

type Chaincode struct{
}

type Contact struct{
	ContactId string
	Name string
	Age int
	Gender string
	Race string
}

type Case struct{
	T_Stage string
	N_Stage string
	Grade string
	Condition string
	Survival_Time int												//Survival Time is in percentage.
	Cancer_Diagnosis string
	Metastasis_Location string
	Cancer_Stage string
	NCCN_Distress_Score	int
} 

type BackgroundInformation struct{
	Affected_Breast string
	ER_Status_LB string
	ER_Status_RB string
	HER2_Status_LB string
	HER2_Status_RB string
	PR_Status_LB string
	PR_Status_RB string
}

type MedicalRecord struct{
	Contact_Rec Contact
	Case_Rec Case
	BackgroundInfo BackgroundInformation
}

func main(){
	err := shim.Start(new(Chaincode))
	if err != nil{
		errors.New("MAIN: Error in starting chaincode.")
	}
}

func (c *Chaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	return nil, nil
}

func (c *Chaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	
	if len(args) != 5 {
		errors.New("INVOKE: Incorrect number of arguments passed.")
	}
	
	age, err := strconv.Atoi(args[2])
	if err != nil{
		errors.New("INVOKE: String to int conversion failed.")
	}
	
	if function == "createMedicalRecord" {
		return c.createMedicalRecord(stub, args[0], args[1], age, args[3], args[4])
	} else if function == "updateMedicalRecord" {
		return c.updateMedicalRecord(stub, args[0], args[1], age, args[3], args[4])
	}
	
	return nil, nil
}

func (c *Chaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	return nil, nil
}

func (c *Chaincode) createMedicalRecord(stub shim.ChaincodeStubInterface, contactId string, name string, age int, gender string, race string) ([]byte, error){

	if contactId == ""{
		return nil, errors.New("CREATE_MEDICAL_RECORD: Invalid Contact ID provided.")
	}
	
	var cont Contact
	contactJson  := []byte(`{"ContactId":"` + contactId + `","Name":"` + name + `","Age":` + strconv.Itoa(age) + `,"Gender":"` + gender + `","Race":"` + race + `"}`)
	err := json.Unmarshal(contactJson, &cont)
	if err != nil{
		errors.New("CREATE_MEDICAL_RECORD: Invalid JSON object.")
	}
	
	c.saveRecord(stub, cont)
	
	return nil, nil
}

/*func (c *Chaincode) updateMedicalRecord(stub shim.ChaincodeStubInterface, contactId string, name string, age int, gender string, race string) ([]byte, error){
	return nil, nil
}*/

func (c *Chaincode) saveRecord(stub shim.ChaincodeStubInterface, cont Contact) (bool, error){
	
	contact, err := json.Marshal(cont)
	if err != nil{
		errors.New("SAVE_RECORD: Error encoding Medical Record.")
	}
	
	err = stub.PutState(contact.ContactId, contact)
	if err != nil{
		errors.New("SAVE_RECORD: Error saving Medical Record.")
	}
	
	return true, nil
}
