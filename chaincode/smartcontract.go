package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Photo struct {
	PhotoID    string `json:"photoID"`
	UserName   string `json:"userName"`
	PhotoHash  string `json:"photoHash"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	photos := []Photo{
		{PhotoID: "1", UserName: "Alice", PhotoHash: "hashvalue1"},
		{PhotoID: "2", UserName: "Bob", PhotoHash: "hashvalue2"},
		{PhotoID: "3", UserName: "Bob", PhotoHash: "hashvalue2"},

	}

	for _, photo := range photos {
		photoJSON, err := json.Marshal(photo)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(photo.PhotoID, photoJSON)
		if err != nil {
			return fmt.Errorf("failed to put photo to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreatePhoto(ctx contractapi.TransactionContextInterface, photoID string, userName string, photoHash string) error {
	exists, err := s.PhotoExists(ctx, photoID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", photoID)
	}

	photo := Photo{
		PhotoID:   photoID,
		UserName:  userName,
		PhotoHash: photoHash,
	}

	photoJSON, err := json.Marshal(photo)
	if err != nil {
		return err
	}
	
	return ctx.GetStub().PutState(photoID, photoJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) QueryPhoto(ctx contractapi.TransactionContextInterface, photoID string) (*Photo, error) {
	photoJSON, err := ctx.GetStub().GetState(photoID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if photoJSON == nil {
		return nil, fmt.Errorf("the photo %s does not exist", photoID)
	}

	var photo Photo
	err = json.Unmarshal(photoJSON, &photo)
	if err != nil {
		return nil, err
	}

	return &photo, nil
}


// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeletePhoto(ctx contractapi.TransactionContextInterface, photoID string) error {
	exists, err := s.PhotoExists(ctx, photoID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the photo %s does not exist", photoID)
	}

	return ctx.GetStub().DelState(photoID)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) PhotoExists(ctx contractapi.TransactionContextInterface, photoID string) (bool, error) {
	photoJSON, err := ctx.GetStub().GetState(photoID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return photoJSON != nil, nil
}


// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllPhoto(ctx contractapi.TransactionContextInterface) ([]*Photo, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Photo
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Photo
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

