package impl

import (
	"context"
	"errors"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/request"
	"github.com/UsagiBooru/accounts-server/utils/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MutesApiImplService struct {
	gen.MutesApiService
	md *mongo.Client
	ah mongo_models.MongoAccountHelper
	mh mongo_models.MongoMuteHelper
}

func NewMutesApiImplService(md *mongo.Client) gen.MutesApiServicer {
	return &MutesApiImplService{
		MutesApiService: gen.MutesApiService{},
		md:              md,
		ah:              mongo_models.NewMongoAccountHelper(md),
		mh:              mongo_models.NewMongoMuteHelper(md),
	}
}

// AddMute - Add mute
func (s *MutesApiImplService) AddMute(ctx context.Context, accountID int32, muteStruct gen.MuteStruct) (gen.ImplResponse, error) {
	// Get issuerId/ issuerPermission
	issuerID, issuerPermission, err := request.GetHeaders(ctx)
	if err != nil {
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
	_, err = s.ah.FindAccount(mongo_models.AccountID(issuerID))
	if err != nil {
		return response.NewNotFoundError(), nil
	}
	// Find mute does already exists
	filter := bson.M{
		"targetType": muteStruct.TargetType,
		"targetID":   muteStruct.TargetID,
	}
	_, err = s.mh.FindMuteUsingFilter(filter)
	if err == nil {
		return response.NewConflictedError(), nil
	}
	// Get muteIDSeq
	muteSequenceHelper := mongo_models.NewMongoSequenceHelper(s.md, "accounts", "muteID")
	seq, err := muteSequenceHelper.GetSeq()
	if err != nil {
		return response.NewInternalError(), err
	}
	// Create new mute
	newMute, err := s.mh.CreateMute(
		seq+1,
		muteStruct.TargetType,
		muteStruct.TargetID,
	)
	if err != nil {
		return response.NewInternalError(), err
	}
	// Update seq
	if err := muteSequenceHelper.UpdateSeq(); err != nil {
		return response.NewInternalError(), err
	}
	return gen.Response(200, newMute.ToOpenApi()), nil
}

// DeleteMute - Delete mute
func (s *MutesApiImplService) DeleteMute(ctx context.Context, accountID int32, muteID int32) (gen.ImplResponse, error) {
	// Get issuerId/ issuerPermission
	issuerID, issuerPermission, err := request.GetHeaders(ctx)
	if err != nil {
		return response.NewInternalError(), err
	}
	if err := request.ValidatePermission(issuerPermission, issuerID, accountID); err != nil {
		return response.NewPermissionErrorWithMessage(err.Error()), err
	}
	// Find mute
	_, err = s.mh.FindMute(muteID)
	if err != nil {
		return response.NewNotFoundError(), nil
	}
	// Delete mute
	err = s.mh.DeleteMute(muteID)
	if err != nil {
		return response.NewInternalError(), nil
	}
	return gen.Response(204, nil), nil
}

// GetMute - Get mute
func (s *MutesApiImplService) GetMute(ctx context.Context, accountID int32, muteID int32) (gen.ImplResponse, error) {
	// Get issuerId/ issuerPermission
	issuerID, issuerPermission, err := request.GetHeaders(ctx)
	if err != nil {
		return response.NewInternalError(), err
	}
	if err := request.ValidatePermission(issuerPermission, issuerID, accountID); err != nil {
		return response.NewPermissionErrorWithMessage(err.Error()), err
	}
	// Find mute
	mute, err := s.mh.FindMute(muteID)
	if err != nil {
		return response.NewNotFoundError(), nil
	}
	return gen.Response(200, mute.ToOpenApi()), nil
}

// GetMutes - Get mute list
func (s *MutesApiImplService) GetMutes(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// Get issuerId/ issuerPermission
	issuerID, issuerPermission, err := request.GetHeaders(ctx)
	if err != nil {
		return response.NewInternalError(), err
	}
	if err := request.ValidatePermission(issuerPermission, issuerID, accountID); err != nil {
		return response.NewPermissionErrorWithMessage(err.Error()), err
	}
	//TODO: Uncomment the next line to return gen.Response Response(200, GetMutesResponse{}) or use other options such as http.Ok ...
	//return gen.Response(200, GetMutesResponse{}), nil

	//TODO: Uncomment the next line to return gen.Response Response(403, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return gen.Response(403, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return gen.Response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return gen.Response(404, GeneralMessageResponse{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetMutes method not implemented")
}
