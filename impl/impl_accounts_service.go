package impl

import (
	"context"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/mongo_models"
	"github.com/UsagiBooru/accounts-server/utils/request"
	"github.com/UsagiBooru/accounts-server/utils/response"
	"github.com/UsagiBooru/accounts-server/utils/server"
	jwt "github.com/form3tech-oss/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type AccountsApiImplService struct {
	gen.AccountsApiService
	// es *elasticsearch.Client
	md        *mongo.Client
	ih        mongo_models.MongoInviteHelper
	ah        mongo_models.MongoAccountHelper
	validate  *validator.Validate
	jwtSecret string
}

func NewAccountsApiImplService() gen.AccountsApiServicer {
	conf := server.GetConfig()
	md := server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
	return &AccountsApiImplService{
		AccountsApiService: gen.AccountsApiService{},
		// es:                 server.NewElasticSearchClient(conf.ElasticHost, conf.ElasticUser, conf.ElasticPass),
		md:        md,
		ih:        mongo_models.NewMongoInviteHelper(md),
		ah:        mongo_models.NewMongoAccountHelper(md),
		validate:  validator.New(),
		jwtSecret: conf.JwtSecret,
	}
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// Find target account
	account, err := s.ah.FindAccount(mongo_models.AccountID(accountID))
	if err != nil {
		return response.NewNotFoundError(), nil
	}
	return gen.Response(200, account.ToOpenApi(s.md)), nil
}

// CreateAccount - Create account
func (s *AccountsApiImplService) CreateAccount(ctx context.Context, accountStruct gen.AccountStruct) (gen.ImplResponse, error) {
	// Timeout of this method is 3 seconds
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	// Validate struct
	err := s.validate.Struct(s.ah.ToMongo(accountStruct))
	if err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), nil
	}
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
	var account *mongo_models.MongoAccountStruct
	// Use transaction to prevent duplicate request
	err = s.md.UseSession(ctx, func(sc mongo.SessionContext) error {
		err := sc.StartTransaction()
		if err != nil {
			return err
		}
		// Create sequence helper
		accountSequenceHelper := mongo_models.NewMongoSequenceHelper(s.md, "accounts", "accountID")
		// Get latest -1 accountID
		seq, err := accountSequenceHelper.GetSeq()
		if err != nil {
			return err
		}
		if err := accountSequenceHelper.UpdateSeq(); err != nil {
			return err
		}
		// Get invite info
		invite, err := s.ih.FindInvite(accountStruct.Invite.Code)
		if err != nil {
			return err
		}
		// Find inviter account
		inviterAccountID := invite.Inviter
		newAccountID := mongo_models.AccountID(seq + 1)
		inviter, err := s.ah.FindAccount(inviterAccountID)
		if err != nil {
			return err
		}
		// Use invite code
		if err := s.ih.UseInvite(invite.ID, newAccountID); err != nil {
			return err
		}
		// Generate new invite for new account
		inviteCodeForNew := server.GetShortUUID(8)
		if err := s.ih.CreateInvite(inviteCodeForNew, newAccountID); err != nil {
			return err
		}
		// Generate new invite for inviter account
		inviteCodeForInviter := server.GetShortUUID(8)
		if err := s.ih.CreateInvite(inviteCodeForInviter, inviterAccountID); err != nil {
			return err
		}
		// Create new account
		account, err = s.ah.CreateAccount(
			newAccountID,
			accountStruct.DisplayID,
			accountStruct.Password,
			accountStruct.Mail,
			accountStruct.Name,
			inviterAccountID,
			inviteCodeForNew,
		)
		if err != nil {
			return err
		}
		// Update inviter's invite count
		if err := s.ah.UpdateInvite(
			inviterAccountID,
			inviteCodeForInviter,
			inviter.Invite.InvitedCount+1,
		); err != nil {
			return err
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
	issuerID, issuerPermission, err := request.GetHeaders(ctx)
	if err != nil {
		return response.NewInternalError(), err
	}
	// Validate struct
	err = s.validate.Struct(s.ah.ToMongo(accountChange))
	if err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), nil
	}
	// Find target account
	accountCurrent, err := s.ah.FindAccount(mongo_models.AccountID(accountID))
	if err != nil {
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
	// Update using input
	col := s.md.Database("accounts").Collection("users")
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
	if err := s.ah.UpdateAccount(mongo_models.AccountID(accountID), *accountCurrent); err != nil {
		return response.NewInternalError(), err
	}
	return gen.Response(200, accountCurrent.ToOpenApi(s.md)), nil
}

// DeleteAccount - Delete account info
func (s *AccountsApiImplService) DeleteAccount(ctx context.Context, accountID int32, password string) (gen.ImplResponse, error) {
	issuerID, issuerPermission, err := request.GetHeaders(ctx)
	if err != nil {
		return response.NewInternalError(), err
	}
	// Find target account
	account, err := s.ah.FindAccount(mongo_models.AccountID(accountID))
	if err != nil {
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
	if err := s.ah.UpdateAccount(mongo_models.AccountID(accountID), *account); err != nil {
		return response.NewInternalError(), err
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
	account, err := s.ah.FindAccount(mongo_models.AccountID(issuerID))
	if err != nil {
		return response.NewNotFoundError(), nil
	}
	return gen.Response(200, account.ToOpenApi(s.md)), nil
}
