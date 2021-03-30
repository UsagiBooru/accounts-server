package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils"
	"github.com/UsagiBooru/accounts-server/utils/mongo_models"
	"github.com/elastic/go-elasticsearch/v7"

	"go.mongodb.org/mongo-driver/mongo"
)

type AccountsApiImplService struct {
	gen.AccountsApiService
	es *elasticsearch.Client
	md *mongo.Client
}

func NewAccountsApiImplService() gen.AccountsApiServicer {
	conf := utils.GetConfig()
	return &AccountsApiImplService{
		AccountsApiService: gen.AccountsApiService{},
		es:                 utils.NewElasticSearchClient(conf.ElasticHost, conf.ElasticUser, conf.ElasticPass),
		md:                 utils.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass),
	}
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// TODO - update GetAccount with the required logic for this service method.
	// Add api_accounts_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, AccountStruct{}) or use other options such as http.Ok ...
	//return Response(200, AccountStruct{}), nil

	//TODO: Uncomment the next line to return response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(404, GeneralMessageResponse{}), nil

	// s.es.hogehoge で ElasticSearchが呼べる?
	// s.md.hogehoge で MongoDBが呼べる?
	user := mongo_models.MongoAccount{
		AccountStruct: gen.AccountStruct{
			AccountID: 1,
			DisplayID: "domao",
		},
		TotpKey:      "NewTotpKey",
		PasswordSalt: "",
	}
	user_doc := utils.ConvertOpenApiStructToBson(user)
	col := s.md.Database("accounts").Collection("users")
	_, err := col.InsertOne(ctx, user_doc)
	if err != nil {
		fmt.Println(err)
		return gen.Response(500, gen.GeneralMessageResponse{Message: "Failed"}), nil
	}
	return gen.Response(200, gen.AccountStruct{}), nil
}
