package impl

import (
	"context"
	"errors"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/request"
	"github.com/UsagiBooru/accounts-server/utils/response"
	"github.com/UsagiBooru/accounts-server/utils/server"
	jwt "github.com/form3tech-oss/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AccountsApiImplService struct {
	gen.AccountsApiService
	// es *elasticsearch.Client
	md        *mongo.Client
	jwtSecret string
}

func NewAccountsApiImplService() gen.AccountsApiServicer {
	conf := server.GetConfig()
	return &AccountsApiImplService{
		AccountsApiService: gen.AccountsApiService{},
		// es:                 server.NewElasticSearchClient(conf.ElasticHost, conf.ElasticUser, conf.ElasticPass),
		md:        server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass),
		jwtSecret: conf.JwtSecret,
	}
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var account mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		server.Debug(err.Error())
		return response.NewNotFoundError(), nil
	}
	return gen.Response(200, account.ToOpenApi(s.md)), nil
}

// CreateAccount - Create account
func (s *AccountsApiImplService) CreateAccount(ctx context.Context, accountStruct gen.AccountStruct) (gen.ImplResponse, error) {
	// Timeout of this method is 3 seconds
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Validate request fields
	if err := request.ValidateRequiredFields(
		accountStruct,
		[]string{"name", "displayID", "password", "mail"},
	); err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), nil
	}
	if err := request.ValidateRequiredFields(
		accountStruct.Invite,
		[]string{"code"},
	); err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), nil
	}

	var account mongo_models.MongoAccountStruct
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
			server.Debug(err.Error())
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
		newInviteCodeForNew := server.GetShortUUID(8)
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
		newInviteCodeForOld := server.GetShortUUID(8)
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
		account = mongo_models.MongoAccountStruct{
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
		if _, err = col.InsertOne(ctx, account); err != nil {
			return errors.New("insert new account failed")
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
		server.Debug(err.Error())
		return response.NewInternalError(), nil
	}
	return gen.Response(200, account.ToOpenApi(s.md)), nil
}

// EditAccount - Edit account info
func (s *AccountsApiImplService) EditAccount(ctx context.Context, accountID int32, accountChange gen.AccountStruct) (gen.ImplResponse, error) {
	// Get issuer id/permission
	issuerID, err := request.GetUserID(ctx)
	issuerPermission, err2 := request.GetUserPermission(ctx)
	if err != nil || err2 != nil {
		if err != nil {
			server.Debug(err.Error())
		} else {
			server.Debug(err2.Error())
		}
		return response.NewInternalError(), nil
	}
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var accountCurrent mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&accountCurrent); err != nil {
		server.Debug(err.Error())
		return response.NewNotFoundError(), nil
	}

	/* Validate Permission */
	notAdmin := issuerPermission != request.PermissionAdmin
	notMod := issuerPermission < request.PermissionModerator
	notSelf := accountID != issuerID
	notSelfOrAdmin := notAdmin && notSelf
	// Deny changing invite / inviter / notify
	if (accountChange.Invite != gen.AccountStructInvite{}) ||
		(accountChange.Inviter != gen.LightAccountStruct{}) ||
		(accountChange.Notify != gen.AccountStructNotify{}) {
		server.Debug("Denied since tried to change invite / inviter / notify")
		return response.NewRequestError(), nil
	}
	// Deny changing different account if not greater than moderator
	if notSelf && notMod {
		return response.NewPermissionError(), nil
	}
	// Deny changing if target permission is greater than moderator except target is ownself
	if issuerPermission == request.PermissionModerator &&
		accountCurrent.Permission >= request.PermissionModerator &&
		notSelf {
		server.Debug("Denied since changing permission with moderator, and target was not normal user.")
		return response.NewPermissionError(), nil
	}
	// Deny changing permission if not admin
	if accountChange.Permission != accountCurrent.Permission && notAdmin {
		server.Debug("Denied since changing permission with not admin")
		return response.NewPermissionError(), nil
	}
	// Deny changing access if not greater than moderator
	if (accountChange.Access != gen.AccountStructAccess{}) && notMod {
		server.Debug("Denied since changing access with not greater than moderator")
		return response.NewPermissionError(), nil
	}
	// Deny changing password if not admin except target is ownself
	if accountChange.Password != "" && notSelfOrAdmin {
		server.Debug("Denied since changing password with not admin and target wasn't ownself")
		return response.NewPermissionError(), nil
	}
	// Deny changing totp if not admin except target is ownself
	if accountChange.TotpEnabled != accountCurrent.TotpEnabled && notSelfOrAdmin {
		server.Debug("Denied since changing totp with not admin and target wasn't ownself")
		return response.NewPermissionError(), nil
	}
	// Deny changing mail if not admin except target is ownself
	if accountChange.Mail != "" && notSelfOrAdmin {
		server.Debug("Denied since changing mail with not admin and target wasn't ownself")
		return response.NewPermissionError(), nil
	}

	/* Update using input */
	if err := accountCurrent.UpdateDisplayID(col, accountChange.DisplayID); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateName(col, accountChange.Name); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateApiSeq(accountChange.ApiSeq); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdatePermission(accountChange.Permission); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdatePassword(accountChange.OldPassword, accountChange.Password); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateDescription(accountChange.Description); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateMail(accountChange.Mail); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateFavorite(accountChange.Favorite); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateAccess(accountChange.Access); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateIpfs(accountChange.Ipfs); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}

	// Update account
	filter = bson.M{"accountID": accountCurrent.AccountID}
	set := bson.M{"$set": accountCurrent}
	if _, err = col.UpdateOne(ctx, filter, set); err != nil {
		server.Debug(err.Error())
		return response.NewInternalError(), nil
	}
	return gen.Response(200, accountCurrent.ToOpenApi(s.md)), nil
}

// DeleteAccount - Delete account info
func (s *AccountsApiImplService) DeleteAccount(ctx context.Context, accountID int32, password string) (gen.ImplResponse, error) {
	// Get issuer id/permission
	issuerID, err := request.GetUserID(ctx)
	issuerPermission, err2 := request.GetUserPermission(ctx)
	if err != nil || err2 != nil {
		if err != nil {
			server.Debug(err.Error())
		} else {
			server.Debug(err2.Error())
		}
		return response.NewInternalError(), nil
	}
	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": accountID}
	var account mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		return response.NewNotFoundError(), nil
	}

	// Validate permission
	notMod := issuerPermission < request.PermissionModerator
	if accountID != issuerID && notMod {
		return response.NewPermissionError(), nil
	}
	// Validate old password hash
	if notMod {
		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
			return response.NewPermissionErrorWithMessage("password mismatched"), nil
		}
	}

	// Update account
	if issuerPermission == request.PermissionUser {
		account.AccountStatus = mongo_models.STATUS_DELETED_SELF
	} else {
		account.AccountStatus = mongo_models.STATUS_DELETED_MOD
	}
	filter = bson.M{"accountID": account.AccountID}
	set := bson.M{"$set": account}
	if _, err = col.UpdateOne(ctx, filter, set); err != nil {
		server.Debug(err.Error())
		return response.NewInternalError(), nil
	}
	return gen.Response(204, nil), nil
}

// LoginWithForm - Login with form
func (s *AccountsApiImplService) LoginWithForm(ctx context.Context, req gen.PostLoginWithFormRequest) (gen.ImplResponse, error) {
	accountIdOrMail := req.Id
	accountPassword := req.Password

	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"displayID": accountIdOrMail}
	var account mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		return response.NewNotFoundError(), nil
	}
	if err := account.ValidatePassword(accountPassword); err != nil {
		return response.NewUnauthorizedError(), nil
	}

	// Deny if account deleted
	if account.AccountStatus != mongo_models.STATUS_NORMAL {
		return response.NewLockedErrorWithMessage("the account was deleted"), nil
	}

	// Generate jwt token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	TWO_MONTH := time.Hour * 24 * 60
	claims["sub"] = account.AccountID
	claims["name"] = account.Name
	claims["permission"] = account.Permission
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(TWO_MONTH).Unix()
	signed_token, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return response.NewInternalError(), nil
	}

	return gen.Response(200, gen.PostLoginWithFormResponse{ApiKey: signed_token}), nil
}

// GetUploadHistory - Get upload history
func (s *AccountsApiImplService) GetUploadHistory(ctx context.Context, accountID int32, page int32, sort string, order string, perPage int32) (gen.ImplResponse, error) {
	// TODO - update GetUploadHistory with the required logic for this service method.

	//TODO: Uncomment the next line to return response Response(404, GeneralMessageResponse{}) or use other options such as http.Ok ...
	//return Response(404, GeneralMessageResponse{}), nil
	return gen.Response(200, gen.GetUploadHistoryResponse{}), nil
}

// GetAccountMe - Get user info (self)
func (s *AccountsApiImplService) GetAccountMe(ctx context.Context) (gen.ImplResponse, error) {
	// Get issuer id/permission
	issuerID, err := request.GetUserID(ctx)
	if err != nil {
		server.Debug(err.Error())
		return response.NewInternalError(), nil
	}

	// Find target account
	col := s.md.Database("accounts").Collection("users")
	filter := bson.M{"accountID": issuerID}
	var account mongo_models.MongoAccountStruct
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		server.Debug(err.Error())
		return response.NewNotFoundError(), nil
	}
	return gen.Response(200, account.ToOpenApi(s.md)), nil
}
