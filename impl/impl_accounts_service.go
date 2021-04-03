package impl

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils"
	"github.com/UsagiBooru/accounts-server/utils/mongo_models"
	"github.com/elastic/go-elasticsearch/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var account gen.AccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		utils.Debug(err.Error())
		return utils.NewNotFoundError(), nil
	}
	// Find inviter account
	col = s.md.Database("accounts").Collection("users")
	filter = bson.M{"accountID": account.Inviter.AccountID}
	var inviter gen.AccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&inviter); err != nil {
		utils.Debug(err.Error())
		return utils.NewInternalError(), nil
	}
	// Create response
	accountResp := gen.AccountStruct{
		AccountID:   account.AccountID,
		DisplayID:   account.DisplayID,
		Permission:  account.Permission,
		Name:        account.Name,
		Description: account.Description,
		Favorite:    account.Favorite,
		Access: gen.AccountStructAccess{
			CanInvite:      account.Access.CanInvite,
			CanLike:        account.Access.CanLike,
			CanComment:     account.Access.CanComment,
			CanCreatePost:  account.Access.CanCreatePost,
			CanEditPost:    account.Access.CanEditPost,
			CanApprovePost: account.Access.CanApprovePost,
		},
		Inviter: gen.LightAccountStruct{
			AccountID: account.Inviter.AccountID,
			Name:      inviter.Name,
		},
	}
	return gen.Response(200, accountResp), nil
}

// CreateAccount - Create account
func (s *AccountsApiImplService) CreateAccount(ctx context.Context, accountStruct gen.AccountStruct) (gen.ImplResponse, error) {
	// Timeout of this method is 3 seconds
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Validate request fields
	if resp := utils.ValidateRequiredFields(
		accountStruct,
		[]string{"name", "displayID", "password", "mail"},
	); resp.Code != http.StatusOK {
		return resp, nil
	}
	if resp := utils.ValidateRequiredFields(
		accountStruct.Invite,
		[]string{"code"},
	); resp.Code != http.StatusOK {
		return resp, nil
	}

	var user mongo_models.MongoAccount
	// Use transaction to prevent duplicate request
	err := s.md.UseSession(ctx, func(sc mongo.SessionContext) error {
		// Start transaction
		err := sc.StartTransaction()
		if err != nil {
			return err
		}
		// Find invite code
		col := s.md.Database("accounts").Collection("invites")
		filter := bson.M{
			"code":    accountStruct.Invite.Code,
			"invitee": 0,
		}
		var invite mongo_models.MongoInvite
		if err := col.FindOne(context.Background(), filter).Decode(&invite); err != nil {
			return errors.New("invite code was not found")
		}
		// Find inviter account
		col = s.md.Database("accounts").Collection("users")
		filter = bson.M{"accountID": invite.Inviter}
		var inviter mongo_models.MongoAccount
		if err := col.FindOne(context.Background(), filter).Decode(&inviter); err != nil {
			return errors.New("inviter account was not found")
		}
		// Get latest-1 accountID
		col = s.md.Database("accounts").Collection("sequence")
		filter = bson.M{"key": "accountID"}
		var seq mongo_models.MongoSequence
		if err := col.FindOne(context.Background(), filter).Decode(&seq); err != nil {
			return errors.New("accountID sequence was not found")
		}
		// Update invite invitee
		col = s.md.Database("accounts").Collection("invites")
		filter = bson.M{
			"_id": invite.ID,
		}
		set := bson.M{"$set": bson.M{"invitee": seq.Value + 1}}
		if _, err = col.UpdateOne(ctx, filter, set); err != nil {
			return errors.New("update invite invitee failed")
		}
		// Get password hash
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(accountStruct.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			return errors.New("password hash create failed")
		}
		// Create new invite for new account
		newInviteCodeForNew := utils.GetShortUUID(8)
		newInviteForNew := mongo_models.MongoInvite{
			ID:      primitive.NewObjectID(),
			Code:    newInviteCodeForNew,
			Inviter: seq.Value + 1,
			Invitee: 0,
		}
		if _, err = col.InsertOne(ctx, utils.ConvertStructToBson(newInviteForNew)); err != nil {
			return errors.New("insert new invite for new account failed")
		}
		// Create new invite for old account
		newInviteCodeForOld := utils.GetShortUUID(8)
		newInviteForOld := mongo_models.MongoInvite{
			ID:      primitive.NewObjectID(),
			Code:    newInviteCodeForOld,
			Inviter: invite.Inviter,
		}
		if _, err = col.InsertOne(ctx, utils.ConvertStructToBson(newInviteForOld)); err != nil {
			return errors.New("insert new invite for old account failed")
		}
		// Create new mongo user model
		col = s.md.Database("accounts").Collection("users")
		user = mongo_models.MongoAccount{
			ID: primitive.NewObjectID(),
			AccountStruct: gen.AccountStruct{
				AccountID:   seq.Value + 1,
				DisplayID:   accountStruct.DisplayID,
				ApiSeq:      0,
				Permission:  0,
				Password:    string(hashedPassword),
				Mail:        accountStruct.Mail,
				TotpEnabled: false,
				Name:        accountStruct.Name,
				Description: "",
				Favorite:    0,
				Access: gen.AccountStructAccess{
					CanInvite:      true,
					CanLike:        true,
					CanComment:     true,
					CanCreatePost:  true,
					CanEditPost:    false,
					CanApprovePost: false,
				},
				Inviter: gen.LightAccountStruct{
					AccountID: invite.Inviter,
				},
				Invite: gen.AccountStructInvite{
					InvitedCount: -1,
					Code:         newInviteCodeForNew,
				},
				Notify: gen.AccountStructNotify{
					HasLineNotify: false,
					HasWebNotify:  false,
				},
				Ipfs: gen.AccountStructIpfs{
					GatewayUrl:     "https://cloudflare-ipfs.com",
					NodeUrl:        "",
					GatewayEnabled: false,
					NodeEnabled:    false,
					PinEnabled:     false,
				},
			},
		}
		// Insert new user
		if _, err = col.InsertOne(ctx, utils.ConvertStructToBson(user)); err != nil {
			return errors.New("insert new user failed")
		}
		// Update sequence
		col = s.md.Database("accounts").Collection("sequence")
		filter = bson.M{"key": "accountID"}
		set = bson.M{"$set": bson.M{"value": seq.Value + 1}}
		if _, err = col.UpdateOne(ctx, filter, set); err != nil {
			return errors.New("update accountID sequence failed")
		}
		// Update inviter's invite count
		col = s.md.Database("accounts").Collection("users")
		filter = bson.M{"accountID": invite.Inviter}
		if inviter.Invite.InvitedCount == -1 {
			inviter.Invite.InvitedCount = 0
		}
		set = bson.M{"$set": bson.M{
			"invite.invitedCount": inviter.Invite.InvitedCount + 1,
			"invite.code":         newInviteCodeForOld,
		}}
		if _, err := col.UpdateOne(ctx, filter, set); err != nil {
			return errors.New("update inviter's invite count failed")
		}
		// Commit insert user / update sequence / update invite code
		return sc.CommitTransaction(sc)

	})
	if err != nil {
		utils.Debug(err.Error())
		return utils.NewInternalError(), nil
	}
	user.AccountStruct.Password = ""
	return gen.Response(200, user.AccountStruct), nil
}

// EditAccount - Edit account info
func (s *AccountsApiImplService) EditAccount(ctx context.Context, accountID int32, accountStruct gen.AccountStruct) (gen.ImplResponse, error) {
	issuerID, err := utils.GetUserID(ctx)
	issuerPermission, err2 := utils.GetUserPermission(ctx)
	if err != nil || err2 != nil {
		if err != nil {
			utils.Debug(err.Error())
		} else {
			utils.Debug(err2.Error())
		}
		return gen.Response(500, gen.GeneralMessageResponse{Message: utils.MessageInternalError}), nil
	}
	if accountID != issuerID && issuerPermission != utils.PermissionAdmin {
		return gen.Response(403, gen.GeneralMessageResponse{Message: utils.MessagePermissionError}), nil
	}

	// TODO - update EditAccount with the required logic for this service method.
	// Add api_accounts_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, AccountStruct{}) or use other options such as http.Ok ...
	//return Response(200, AccountStruct{}), nil

	//TODO: Uncomment the next line to return response Response(400, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(400, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return response Response(403, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(403, GeneralMessageResponse{}), nil

	//TODO: Uncomment the next line to return response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(404, GeneralMessageResponse{}), nil

	return gen.Response(http.StatusOK, gen.AccountStruct{}), nil
}
