package impl

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils"
	"github.com/UsagiBooru/accounts-server/utils/mongo_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AccountsApiImplService struct {
	gen.AccountsApiService
	// es *elasticsearch.Client
	md *mongo.Client
}

func NewAccountsApiImplService() gen.AccountsApiServicer {
	conf := utils.GetConfig()
	return &AccountsApiImplService{
		AccountsApiService: gen.AccountsApiService{},
		// es:                 utils.NewElasticSearchClient(conf.ElasticHost, conf.ElasticUser, conf.ElasticPass),
		md: utils.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass),
	}
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var account mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		utils.Debug(err.Error())
		return utils.NewNotFoundError(), nil
	}
	return gen.Response(200, account.ToOpenApi(s.md)), nil
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

	var user mongo_models.MongoAccountStruct
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
		var inviter mongo_models.MongoAccountStruct
		if err := col.FindOne(context.Background(), filter).Decode(&inviter); err != nil {
			utils.Debug(err.Error())
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
		if _, err = col.InsertOne(ctx, newInviteForNew); err != nil {
			return errors.New("insert new invite for new account failed")
		}
		// Create new invite for old account
		newInviteCodeForOld := utils.GetShortUUID(8)
		newInviteForOld := mongo_models.MongoInvite{
			ID:      primitive.NewObjectID(),
			Code:    newInviteCodeForOld,
			Inviter: invite.Inviter,
		}
		if _, err = col.InsertOne(ctx, newInviteForOld); err != nil {
			return errors.New("insert new invite for old account failed")
		}
		// Create new mongo user model
		col = s.md.Database("accounts").Collection("users")
		user = mongo_models.MongoAccountStruct{
			ID:            primitive.NewObjectID(),
			AccountStatus: 0,
			AccountID:     seq.Value + 1,
			DisplayID:     accountStruct.DisplayID,
			ApiKey:        "",
			ApiSeq:        0,
			Permission:    0,
			Password:      string(hashedPassword),
			Mail:          accountStruct.Mail,
			TotpCode:      "",
			TotpEnabled:   false,
			Name:          accountStruct.Name,
			Description:   "",
			Favorite:      0,
			Access: mongo_models.MongoAccountStructAccess{
				CanInvite:      true,
				CanLike:        true,
				CanComment:     true,
				CanCreatePost:  true,
				CanEditPost:    false,
				CanApprovePost: false,
			},
			Inviter: mongo_models.LightMongoAccountStruct{
				AccountID: invite.Inviter,
			},
			Invite: mongo_models.MongoAccountStructInvite{
				InvitedCount: -1,
				Code:         newInviteCodeForNew,
			},
			Notify: mongo_models.MongoAccountStructNotify{
				HasLineNotify: false,
				HasWebNotify:  false,
			},
			Ipfs: mongo_models.MongoAccountStructIpfs{
				GatewayUrl:     "https://cloudflare-ipfs.com",
				NodeUrl:        "",
				GatewayEnabled: false,
				NodeEnabled:    false,
				PinEnabled:     false,
			},
		}
		// Insert new user
		if _, err = col.InsertOne(ctx, user); err != nil {
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
	return gen.Response(200, user.ToOpenApi(s.md)), nil
}

// EditAccount - Edit account info
func (s *AccountsApiImplService) EditAccount(ctx context.Context, accountID int32, accountChange gen.AccountStruct) (gen.ImplResponse, error) {
	// Get issuer id/permission
	issuerID, err := utils.GetUserID(ctx)
	issuerPermission, err2 := utils.GetUserPermission(ctx)
	if err != nil || err2 != nil {
		if err != nil {
			utils.Debug(err.Error())
		} else {
			utils.Debug(err2.Error())
		}
		return utils.NewInternalError(), nil
	}
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var accountCurrent mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&accountCurrent); err != nil {
		utils.Debug(err.Error())
		return utils.NewNotFoundError(), nil
	}

	/* Validate Permission */
	notAdmin := issuerPermission != utils.PermissionAdmin
	notMod := issuerPermission < utils.PermissionModerator
	notSelf := accountID != issuerID
	notSelfOrAdmin := notAdmin && notSelf
	// Deny changing invite / inviter / notify
	if (accountChange.Invite != gen.AccountStructInvite{}) ||
		(accountChange.Inviter != gen.LightAccountStruct{}) ||
		(accountChange.Notify != gen.AccountStructNotify{}) {
		utils.Debug("Denied since tried to change invite / inviter / notify")
		return utils.NewRequestError(), nil
	}
	// Deny changing different account if not greater than moderator
	if notSelf && notMod {
		return utils.NewPermissionError(), nil
	}
	// Deny changing if target permission is greater than moderator except target is ownself
	if issuerPermission == utils.PermissionModerator &&
		accountCurrent.Permission >= utils.PermissionModerator &&
		notSelf {
		utils.Debug("Denied since changing permission with moderator, and target was not normal user.")
		return utils.NewPermissionError(), nil
	}
	// Deny changing permission if not admin
	if accountChange.Permission != accountCurrent.Permission && notAdmin {
		utils.Debug("Denied since changing permission with not admin")
		return utils.NewPermissionError(), nil
	}
	// Deny changing access if not greater than moderator
	if (accountChange.Access != gen.AccountStructAccess{}) && notMod {
		utils.Debug("Denied since changing access with not greater than moderator")
		return utils.NewPermissionError(), nil
	}
	// Deny changing password if not admin except target is ownself
	if accountChange.Password != "" && notSelfOrAdmin {
		utils.Debug("Denied since changing password with not admin and target wasn't ownself")
		return utils.NewPermissionError(), nil
	}
	// Deny changing totp if not admin except target is ownself
	if accountChange.TotpEnabled != accountCurrent.TotpEnabled && notSelfOrAdmin {
		utils.Debug("Denied since changing totp with not admin and target wasn't ownself")
		return utils.NewPermissionError(), nil
	}
	// Deny changing mail if not admin except target is ownself
	if accountChange.Mail != "" && notSelfOrAdmin {
		utils.Debug("Denied since changing mail with not admin and target wasn't ownself")
		return utils.NewPermissionError(), nil
	}

	/* Update using input */
	if err := accountCurrent.UpdateDisplayID(col, accountChange.DisplayID); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateName(col, accountChange.Name); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateApiSeq(accountChange.ApiSeq); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdatePermission(accountChange.Permission); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdatePassword(accountChange.OldPassword, accountChange.Password); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateDescription(accountChange.Description); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateMail(accountChange.Mail); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateFavorite(accountChange.Favorite); err != nil {
		return utils.NewLockedErrorWithMessage(err.Error()), nil
	}
	if (accountChange.Access != gen.AccountStructAccess{}) {
		accountCurrent.Access = mongo_models.MongoAccountStructAccess(accountChange.Access)
	}
	if (accountChange.Ipfs != gen.AccountStructIpfs{}) {
		accountCurrent.Ipfs = mongo_models.MongoAccountStructIpfs(accountChange.Ipfs)
	}

	// Update account
	filter = bson.M{"accountID": accountCurrent.AccountID}
	set := bson.M{"$set": accountCurrent}
	if _, err = col.UpdateOne(ctx, filter, set); err != nil {
		utils.Debug(err.Error())
		return utils.NewInternalError(), nil
	}
	return gen.Response(200, accountCurrent.ToOpenApi(s.md)), nil
}
