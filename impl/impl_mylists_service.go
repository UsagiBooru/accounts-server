package impl

import (
	"context"
	"errors"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/server"
	"go.mongodb.org/mongo-driver/mongo"
)

type MylistApiImplService struct {
	gen.MylistApiService
	md *mongo.Client
	ah mongo_models.MongoAccountHelper
	mh mongo_models.MongoMuteHelper
}

func NewMylistApiImplService() gen.MylistApiServicer {
	conf := server.GetConfig()
	md := server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
	return &MylistApiImplService{
		MylistApiService: gen.MylistApiService{},
		md:               md,
		ah:               mongo_models.NewMongoAccountHelper(md),
		mh:               mongo_models.NewMongoMuteHelper(md),
	}
}

// CreateMylist - Create user mylist
func (s *MylistApiImplService) CreateMylist(ctx context.Context, accountID int32, mylistStruct gen.MylistStruct) (gen.ImplResponse, error) {
	// TODO - update CreateMylist with the required logic for this service method.
	// Add api_mylist_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, MylistStruct{}) or use other options such as http.Ok ...
	//return Response(200, MylistStruct{}), nil

	//TODO: Uncomment the next line to return response Response(400, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(400, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return response Response(403, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(403, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return response Response(409, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(409, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return response Response(429, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(429, GeneralMessageResponse{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("CreateMylist method not implemented")
}

// GetUserMylists - Get user mylists
func (s *MylistApiImplService) GetUserMylists(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// TODO - update GetUserMylists with the required logic for this service method.
	// Add api_mylist_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, GetMylistListResponse{}) or use other options such as http.Ok ...
	//return Response(200, GetMylistListResponse{}), nil

	//TODO: Uncomment the next line to return response Response(403, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(403, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(404, GeneralMessageResponse{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetUserMylists method not implemented")
}
