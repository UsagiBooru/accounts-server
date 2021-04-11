package impl

import (
	"context"
	"errors"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/request"
	"github.com/UsagiBooru/accounts-server/utils/response"
	"github.com/UsagiBooru/accounts-server/utils/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MutesApiImplService struct {
	gen.MutesApiService
	md *mongo.Client
}

func NewMutesApiImplService() gen.MutesApiServicer {
	conf := server.GetConfig()
	return &MutesApiImplService{
		MutesApiService: gen.MutesApiService{},
		md:              server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass),
	}
}

// AddMute - Add mute
func (s *MutesApiImplService) AddMute(ctx context.Context, accountID int32, muteStruct gen.MuteStruct) (gen.ImplResponse, error) {
	// Get issuer id/permission
	var (
		issuerID         int32
		issuerPermission int32
		err              error
	)
	if issuerID, err = request.GetUserID(ctx); err != nil {
		return response.NewInternalError(), err
	}
	if issuerPermission, err = request.GetUserID(ctx); err != nil {
		return response.NewInternalError(), err
	}
	// Validate request
	if err := request.ValidateRequiredFields(muteStruct, []string{"targetType", "targetID"}); err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), err
	}
	// Validate permission
	if err := request.ValidatePermission(issuerPermission, issuerID, accountID); err != nil {
		return response.NewPermissionErrorWithMessage(err.Error()), err
	}
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var account mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		return response.NewNotFoundError(), nil
	}
	// Find mute does already exists
	col = s.md.Database("accounts").Collection("mutes")
	filter = bson.M{
		"targetType": muteStruct.TargetType,
		"targetID":   muteStruct.TargetID,
	}
	if err := col.FindOne(context.Background(), filter); err != nil {
		return response.NewConflictedError(), nil
	}
	// Get muteIDSeq
	var muteIDSeq int32 = 0
	if muteIDSeq, err = mongo_models.GetSeq(s.md, "accounts", "muteID"); err != nil {
		return response.NewInternalError(), err
	}
	// Create new mute
	newMute := mongo_models.MongoMuteStruct{
		ID:         primitive.NewObjectID(),
		MuteID:     muteIDSeq + 1,
		TargetType: muteStruct.TargetType,
		TargetID:   muteStruct.TargetID,
	}
	if _, err = col.InsertOne(ctx, newMute); err != nil {
		return response.NewInternalError(), errors.New("insert new mute failed")
	}
	// Update muteIDSeq
	if err = mongo_models.UpdateSeq(s.md, "accounts", "muteID", muteIDSeq); err != nil {
		return response.NewInternalError(), err
	}
	return gen.Response(200, newMute.ToOpenApi()), nil
}

// DeleteMute - Delete mute
func (s *MutesApiImplService) DeleteMute(ctx context.Context, accountID int32, muteID int32) (gen.ImplResponse, error) {
	// Validate permission
	var (
		issuerID         int32
		issuerPermission int32
		err              error
	)
	if issuerID, err = request.GetUserID(ctx); err != nil {
		return response.NewInternalError(), err
	}
	if issuerPermission, err = request.GetUserID(ctx); err != nil {
		return response.NewInternalError(), err
	}
	if err := request.ValidatePermission(issuerPermission, issuerID, accountID); err != nil {
		return response.NewPermissionErrorWithMessage(err.Error()), err
	}
	// Find mute
	col := s.md.Database("accounts").Collection("mutes")
	filter := bson.M{"muteID": muteID}
	if err := col.FindOne(context.Background(), filter); err != nil {
		return response.NewNotFoundError(), nil
	}
	// Delete mute
	filter = bson.M{"muteID": muteID}
	if _, err := col.DeleteOne(context.Background(), filter); err != nil {
		return response.NewInternalError(), err
	}
	return gen.Response(204, nil), nil
}

// GetMute - Get mute
func (s *MutesApiImplService) GetMute(ctx context.Context, accountID int32, muteID int32) (gen.ImplResponse, error) {
	// Validate permission
	var (
		issuerID         int32
		issuerPermission int32
		err              error
	)
	if issuerID, err = request.GetUserID(ctx); err != nil {
		return response.NewInternalError(), err
	}
	if issuerPermission, err = request.GetUserID(ctx); err != nil {
		return response.NewInternalError(), err
	}
	if err := request.ValidatePermission(issuerPermission, issuerID, accountID); err != nil {
		return response.NewPermissionErrorWithMessage(err.Error()), err
	}
	// Find mute
	col := s.md.Database("accounts").Collection("mutes")
	filter := bson.M{"muteID": muteID}
	var mute mongo_models.MongoMuteStruct
	if err := col.FindOne(context.Background(), filter).Decode(&mute); err != nil {
		return response.NewNotFoundError(), nil
	}
	return gen.Response(200, mute.ToOpenApi()), nil
}

// GetMutes - Get mute list
func (s *MutesApiImplService) GetMutes(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// TODO - update GetMutes with the required logic for this service method.
	// Add api_mutes_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return gen.Response Response(200, GetMutesResponse{}) or use other options such as http.Ok ...
	//return gen.Response(200, GetMutesResponse{}), nil

	//TODO: Uncomment the next line to return gen.Response Response(403, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return gen.Response(403, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return gen.Response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return gen.Response(404, GeneralMessageResponse{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetMutes method not implemented")
}
