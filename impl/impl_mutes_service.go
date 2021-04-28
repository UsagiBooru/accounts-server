package impl

import (
	"context"
	"errors"
	"net/http"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/mongomodels"
	"github.com/UsagiBooru/accounts-server/utils/request"
	"github.com/UsagiBooru/accounts-server/utils/response"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

type MutesApiImplService struct {
	gen.MutesApiService
	md       *mongo.Client
	ah       mongomodels.MongoAccountHelper
	mh       mongomodels.MongoMuteHelper
	validate *validator.Validate
}

// NewMutesApiImplService creates mutes api service
func NewMutesApiImplService(md *mongo.Client) gen.MutesApiServicer {
	return &MutesApiImplService{
		MutesApiService: gen.MutesApiService{},
		md:              md,
		ah:              mongomodels.NewMongoAccountHelper(md),
		mh:              mongomodels.NewMongoMuteHelper(md),
		validate:        validator.New(),
	}
}

// AddMute - Add mute
func (s *MutesApiImplService) AddMute(ctx context.Context, accountID int32, muteStruct gen.MuteStruct) (gen.ImplResponse, error) {
	// Validate struct
	err := s.validate.Struct(s.mh.ToMongo(muteStruct))
	if err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), nil
	}
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
	_, err = s.ah.FindAccount(mongomodels.AccountID(issuerID))
	if err != nil {
		return response.NewNotFoundErrorWithMessage("specified account was not found"), nil
	}
	// Find mute does already exists
	err = s.mh.FindDuplicatedMute(muteStruct.TargetType, muteStruct.TargetID, mongomodels.AccountID(issuerID))
	if err != nil {
		return response.NewConflictedError(), nil
	}
	// Use transaction to prevent duplicate request
	var newMute *mongomodels.MongoMuteStruct
	err = s.md.UseSession(ctx, func(sc mongo.SessionContext) error {
		err := sc.StartTransaction()
		if err != nil {
			return err
		}
		// Get muteIDSeq
		muteSequenceHelper := mongomodels.NewMongoSequenceHelper(s.md, "accounts", "muteID")
		seq, err := muteSequenceHelper.GetSeq()
		if err != nil {
			return err
		}
		// Create new mute
		newMute, err = s.mh.CreateMute(
			seq+1,
			muteStruct.TargetType,
			muteStruct.TargetID,
		)
		if err != nil {
			return err
		}
		// Update seq
		if err := muteSequenceHelper.UpdateSeq(); err != nil {
			return err
		}
		return sc.CommitTransaction(sc)
	})
	if err != nil {
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
	// Delete mute
	err = s.mh.DeleteMute(muteID)
	if err != nil {
		return response.NewNotFoundError(), nil
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

	//TODO: Uncomment the next line to return gen.Response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return gen.Response(404, GeneralMessageResponse{}), nil

	return gen.Response(http.StatusNotImplemented, nil), errors.New("GetMutes method not implemented")
}
