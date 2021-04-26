package impl

import (
	"context"
	"time"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/const_models/account_const"
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

func NewAccountsApiImplService(md *mongo.Client, jwtSecret string) gen.AccountsApiServicer {
	return &AccountsApiImplService{
		AccountsApiService: gen.AccountsApiService{},
		// es:                 server.NewElasticSearchClient(conf.ElasticHost, conf.ElasticUser, conf.ElasticPass),
		md:        md,
		ih:        mongo_models.NewMongoInviteHelper(md),
		ah:        mongo_models.NewMongoAccountHelper(md),
		validate:  validator.New(),
		jwtSecret: jwtSecret,
	}
}

// GetAccount - Get account info
func (s *AccountsApiImplService) GetAccount(ctx context.Context, accountID int32) (gen.ImplResponse, error) {
	// Find target account
	account, err := s.ah.FindAccount(mongo_models.AccountID(accountID))
	if err != nil {
		return response.NewNotFoundError(), nil
	}
	// Read permission for block getting deleted account
	issuerPermission, _ := request.GetUserPermission(ctx)
	if issuerPermission == account_const.PERMISSION_USER &&
		account.AccountStatus != account_const.STATUS_ACTIVE {
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
		if err == server.ErrInviteNotFound {
			return response.NewNotFoundErrorWithMessage(err.Error()), nil
		} else {
			return response.NewInternalError(), nil
		}
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
	notAdmin := issuerPermission != account_const.PERMISSION_ADMIN
	notMod := issuerPermission < account_const.PERMISSION_MOD
	notSelf := accountID != issuerID
	notSelfOrAdmin := notAdmin && notSelf
	// Deny changing invite / inviter / notify
	if (accountChange.Invite != gen.AccountStructInvite{}) ||
		(accountChange.Inviter != gen.LightAccountStruct{}) ||
		(accountChange.Notify != gen.AccountStructNotify{}) {
		return response.NewRequestErrorWithMessage("invite/inviter/notify are not editable"), nil
	}
	// Deny changing different account if not greater than moderator
	if notSelf && notMod {
		return response.NewPermissionErrorWithMessage("you can't edit different account!"), nil
	}
	// Deny changing if target permission is greater than moderator except target is ownself
	// Deny changing permission if not admin
	// Deny changing access if not greater than moderator
	if (issuerPermission == account_const.PERMISSION_MOD &&
		accountCurrent.Permission >= account_const.PERMISSION_MOD &&
		notSelf) ||
		(accountChange.Permission != accountCurrent.Permission && notAdmin) ||
		((accountChange.Access != gen.AccountStructAccess{}) && notMod) {
		return response.NewPermissionError(), nil
	}
	// Deny changing password/totp/mail if not admin except target is ownself
	if notSelfOrAdmin {
		if accountChange.Password != "" ||
			accountChange.TotpEnabled != accountCurrent.TotpEnabled ||
			accountChange.Mail != "" {
			return response.NewPermissionError(), nil
		}
	}
	// Update using input
	col := s.md.Database("accounts").Collection("users")
	if err := accountCurrent.UpdateDisplayID(col, accountChange.DisplayID); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdateName(col, accountChange.Name); err != nil {
		return response.NewLockedErrorWithMessage(err.Error()), nil
	}
	if err := accountCurrent.UpdatePassword(accountChange.OldPassword, accountChange.Password); err != nil {
		return response.NewRequestErrorWithMessage(err.Error()), nil
	}
	// Update current instance (they don't return errors since already validated)
	accountCurrent.UpdateDescription(accountChange.Description)
	accountCurrent.UpdatePermission(accountChange.Permission)
	accountCurrent.UpdateApiSeq(accountChange.ApiSeq)
	accountCurrent.UpdateMail(accountChange.Mail)
	accountCurrent.UpdateFavorite(accountChange.Favorite)
	accountCurrent.UpdateAccess(accountChange.Access)
	accountCurrent.UpdateIpfs(accountChange.Ipfs)
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
	notMod := issuerPermission < account_const.PERMISSION_MOD
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
	if issuerPermission == account_const.PERMISSION_USER {
		account.AccountStatus = account_const.STATUS_DELETED_BY_SELF
	} else {
		account.AccountStatus = account_const.STATUS_DELETED_BY_MOD
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
		return response.NewUnauthorizedError(), nil
	}
	if err := account.ValidatePassword(accountPassword); err != nil {
		return response.NewUnauthorizedError(), nil
	}
	// Deny if account deleted
	if account.AccountStatus != account_const.STATUS_ACTIVE {
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
