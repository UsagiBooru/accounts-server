package impl

import (
	"context"

	. "github.com/UsagiBooru/accounts-server/gen"
)

type AccountsApiImplService struct {
	AccountsApiService
}

func NewAccountsApiImplService() AccountsApiServicer {
	return &AccountsApiImplService{}
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (ImplResponse, error) {
	// TODO - update GetAccount with the required logic for this service method.
	// Add api_accounts_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, AccountStruct{}) or use other options such as http.Ok ...
	//return Response(200, AccountStruct{}), nil

	//TODO: Uncomment the next line to return response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(404, GeneralMessageResponse{}), nil

	return Response(200, AccountStruct{}), nil
}
