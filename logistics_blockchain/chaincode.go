// Chaincode for Tokenized Logistics and Transportation Blockchain

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Smart Contract for Logistics
type LogisticsContract struct {
	contractapi.Contract
}

// Shipper and Consignee structure
type Shipper struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Consignee struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Freight Quotation structure
type FreightQuotation struct {
	ID       string `json:"id"`
	Carrier  string `json:"carrier"`
	Amount   string `json:"amount"`
	Confirmed bool   `json:"confirmed"`
}

// Shipment structure
type Shipment struct {
	ID          string `json:"id"`
	ShipperID   string `json:"shipperId"`
	ConsigneeID string `json:"consigneeId"`
	Status      string `json:"status"`
}

// Function to create a new shipper
func (c *LogisticsContract) CreateShipper(ctx contractapi.TransactionContextInterface, id string, name string) error {
	shipper := Shipper{
		ID:   id,
		Name: name,
	}
	shipperJSON, err := json.Marshal(shipper)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, shipperJSON)
}

// Function to create a new consignee
func (c *LogisticsContract) CreateConsignee(ctx contractapi.TransactionContextInterface, id string, name string) error {
	consignee := Consignee{
		ID:   id,
		Name: name,
	}
	consigneeJSON, err := json.Marshal(consignee)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, consigneeJSON)
}

// Function to create a freight quotation
func (c *LogisticsContract) CreateFreightQuotation(ctx contractapi.TransactionContextInterface, id string, carrier string, amount string) error {
	quotation := FreightQuotation{
		ID:       id,
		Carrier:  carrier,
		Amount:   amount,
		Confirmed: false,
	}
	quotationJSON, err := json.Marshal(quotation)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, quotationJSON)
}

// Function to book a shipment
func (c *LogisticsContract) BookShipment(ctx contractapi.TransactionContextInterface, id string, shipperID string, consigneeID string) error {
	shipment := Shipment{
		ID:          id,
		ShipperID:   shipperID,
		ConsigneeID: consigneeID,
		Status:      "Booked",
	}
	shipmentJSON, err := json.Marshal(shipment)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, shipmentJSON)
}

func (c *LogisticsContract) ConfirmFreightQuotation(ctx contractapi.TransactionContextInterface, id string) error {
	quotationJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("could not retrieve quotation: %v", err)
	}
	if quotationJSON == nil {
		return fmt.Errorf("quotation not found: %s", id)
	}

	var quotation FreightQuotation
	err = json.Unmarshal(quotationJSON, &quotation)
	if err != nil {
		return err
	}

	quotation.Confirmed = true
	updatedQuotationJSON, err := json.Marshal(quotation)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, updatedQuotationJSON)
}

func (c *LogisticsContract) TrackShipment(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	shipmentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("could not retrieve shipment: %v", err)
	}
	if shipmentJSON == nil {
		return "", fmt.Errorf("shipment not found: %s", id)
	}

	return string(shipmentJSON), nil
}

func (c *LogisticsContract) UpdateShipmentStatus(ctx contractapi.TransactionContextInterface, id string, status string) error {
	shipmentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("could not retrieve shipment: %v", err)
	}
	if shipmentJSON == nil {
		return fmt.Errorf("shipment not found: %s", id)
	}

	var shipment Shipment
	err = json.Unmarshal(shipmentJSON, &shipment)
	if err != nil {
		return err
	}

	shipment.Status = status
	updatedShipmentJSON, err := json.Marshal(shipment)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, updatedShipmentJSON)
}

func (c *LogisticsContract) SetRoutingInfo(ctx contractapi.TransactionContextInterface, id string, routingInfo string) error {
	shipmentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("could not retrieve shipment: %v", err)
	}
	if shipmentJSON == nil {
		return fmt.Errorf("shipment not found: %s", id)
	}

	var shipment Shipment
	err = json.Unmarshal(shipmentJSON, &shipment)
	if err != nil {
		return err
	}

	// Assuming routingInfo is a string that contains routing details
	shipment.Status = routingInfo
	updatedShipmentJSON, err := json.Marshal(shipment)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, updatedShipmentJSON)
}

func (c *LogisticsContract) AddThirdPartyLogistics(ctx contractapi.TransactionContextInterface, id string, providerName string) error {
	// Logic to add third-party logistics provider
	// This could involve creating a new record in the ledger
	return nil
}

func (c *LogisticsContract) HandleCustomsOperations(ctx contractapi.TransactionContextInterface, shipmentID string, customsInfo string) error {
	// Logic to handle customs operations for a shipment
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(LogisticsContract))
	if err != nil {
		fmt.Printf("Error creating logistics chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting logistics chaincode: %s", err.Error())
	}
}
