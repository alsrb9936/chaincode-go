package chaincode_test

import (
    "encoding/json"
    "testing"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/hyperledger/fabric-contract-api-go/contractapi/mocks"
    "github.com/stretchr/testify/require"
    "smartcontract/chaincode"
)

func configureStub() (*chaincode.SmartContract, *mocks.TransactionContext) {
    sc := new(chaincode.SmartContract)
    ctx := new(mocks.TransactionContext)
    stub := new(mocks.ChaincodeStub)
    ctx.On("GetStub").Return(stub)
    return sc, ctx
}

func TestInitLedger(t *testing.T) {
    sc, ctx := configureStub()
    ctx.GetStub().On("PutState", mock.Anything, mock.Anything).Return(nil)

    err := sc.InitLedger(ctx)
    require.NoError(t, err)
}

func TestCreatePhoto(t *testing.T) {
    sc, ctx := configureStub()
    ctx.GetStub().On("GetState", "1").Return(nil, nil) // Assume "1" does not exist
    ctx.GetStub().On("PutState", "1", mock.Anything).Return(nil)

    err := sc.CreatePhoto(ctx, "1", "Alice", "hashvalue1")
    require.NoError(t, err)

    // Test for existing photo
    ctx.GetStub().On("GetState", "1").Return([]byte("exists"), nil)
    err = sc.CreatePhoto(ctx, "1", "Alice", "hashvalue1")
    require.Error(t, err)
}

func TestQueryPhoto(t *testing.T) {
    sc, ctx := configureStub()
    photo := chaincode.Photo{PhotoID: "1", UserName: "Alice", PhotoHash: "hashvalue1"}
    photoBytes, _ := json.Marshal(photo)
    ctx.GetStub().On("GetState", "1").Return(photoBytes, nil)

    result, err := sc.QueryPhoto(ctx, "1")
    require.NoError(t, err)
    require.Equal(t, &photo, result)

    // Test for non-existing photo
    ctx.GetStub().On("GetState", "2").Return(nil, nil)
    _, err = sc.QueryPhoto(ctx, "2")
    require.Error(t, err)
}

func TestDeletePhoto(t *testing.T) {
    sc, ctx := configureStub()
    ctx.GetStub().On("GetState", "1").Return([]byte("exists"), nil)
    ctx.GetStub().On("DelState", "1").Return(nil)

    err := sc.DeletePhoto(ctx, "1")
    require.NoError(t, err)

    // Test for non-existing photo
    ctx.GetStub().On("GetState", "2").Return(nil, nil)
    err = sc.DeletePhoto(ctx, "2")
    require.Error(t, err)
}

func TestPhotoExists(t *testing.T) {
    sc, ctx := configureStub()
    ctx.GetStub().On("GetState", "1").Return([]byte("exists"), nil)

    exists, err := sc.PhotoExists(ctx, "1")
    require.NoError(t, err)
    require.True(t, exists)

    // Test for non-existing photo
    ctx.GetStub().On("GetState", "2").Return(nil, nil)
    exists, err = sc.PhotoExists(ctx, "2")
    require.NoError(t, err)
    require.False(t, exists)
}

func TestGetAllPhoto(t *testing.T) {
    sc, ctx := configureStub()
    photos := []chaincode.Photo{
        {PhotoID: "1", UserName: "Alice", PhotoHash: "hashvalue1"},
        {PhotoID: "2", UserName: "Bob", PhotoHash: "hashvalue2"},
    }
    photosBytes, _ := json.Marshal(photos)
    ctx.GetStub().On("GetStateByRange", "", "").Return(mock.NewQueryIterator(photosBytes), nil)

    result, err := sc.GetAllPhoto(ctx)
    require.NoError(t, err)
    require.Len(t, result, len(photos))
}