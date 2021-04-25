package tests

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/impl"
	"github.com/UsagiBooru/accounts-server/utils/server"
)

var parallelFlag = flag.Bool("docker", false, "Set true to use parallel test(Local), otherwise to simple test(CI)")

func GetAccountsServer() (*httptest.Server, func(), bool) {
	var db *mongo.Client
	var shutdown func()
	var err error
	var isParallel bool
	if *parallelFlag {
		server.Debug("Using mongo container")
		db, shutdown, err = GenerateMongoTestContainer()
		isParallel = true
	} else {
		server.Debug("Using mongo server")
		conf := server.GetConfig()
		db = server.NewMongoDBClient(conf.MongoHost, conf.MongoUser, conf.MongoPass)
		shutdown = func() {}
		err = nil
		isParallel = false
	}
	if err != nil {
		server.Fatal(err.Error())
	}
	if err := ReGenerateDatabase(db); err != nil {
		server.Fatal(err.Error())
	}
	AccountsApiService := impl.NewAccountsApiImplService(db, JWT_SECRET)
	AccountsApiController := gen.NewAccountsApiController(AccountsApiService)
	router := server.NewRouterWithInject(AccountsApiController)
	return httptest.NewServer(router), shutdown, isParallel
}

func TestGetAccountSuccessOnValid(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAccountNotFoundOnInvalidId(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/404", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAccountNotFoundOnDeletedIdFromUser(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/4", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAccountSuccessOnDeletedIdFromMod(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/4", nil)
	req = SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateAccountSuccessOnValid(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newAccount := gen.AccountStruct{
		Name:      "デバッグアカウント",
		DisplayID: "debugaccount",
		Password:  "debugaccount",
		Mail:      "mail@example.com",
		Invite: gen.AccountStructInvite{
			Code: "devcode1",
		},
	}
	user_json, _ := json.Marshal(newAccount)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		bytes.NewBuffer(user_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateAccountBadRequestOnInvalidCode(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newAccount := gen.AccountStruct{
		Name:      "デバッグアカウント",
		DisplayID: "debugaccount",
		Password:  "debugaccount",
		Mail:      "mail@example.com",
		Invite: gen.AccountStructInvite{
			Code: "invalidcode",
		},
	}
	user_json, _ := json.Marshal(newAccount)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		bytes.NewBuffer(user_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateAccountBadRequestOnInvalidMail(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newAccount := gen.AccountStruct{
		Name:      "デバッグアカウント",
		DisplayID: "debugaccount",
		Password:  "debugaccount",
		Mail:      "mailaddress",
		Invite: gen.AccountStructInvite{
			Code: "devcode1",
		},
	}
	user_json, _ := json.Marshal(newAccount)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts",
		bytes.NewBuffer(user_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEditAccountSuccessOnChangeName(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Name: "デバッグアカウント2",
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/1",
		bytes.NewBuffer(req_json),
	)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAccountMeFromAdmin(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts/me",
		nil,
	)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAccountMeFromMod(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts/me",
		nil,
	)
	req = SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAccountMeFromNormal(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts/me",
		nil,
	)
	req = SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoginWithFormSuccessOnValid(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.PostLoginWithFormRequest{
		Id:       "domao",
		Password: PASSWORD,
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/login/form",
		bytes.NewBuffer(req_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoginWithFormNotFoundOnInvalidId(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.PostLoginWithFormRequest{
		Id:       "omadosan",
		Password: PASSWORD,
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/login/form",
		bytes.NewBuffer(req_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginWithFormUnAuthorizedOnInvalidPass(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	loginAccount := gen.PostLoginWithFormRequest{
		Id:       "domao",
		Password: "invalid-pass",
	}
	req_json, _ := json.Marshal(loginAccount)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/login/form",
		bytes.NewBuffer(req_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteAccountFromAdmin(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodDelete,
		"/accounts/1",
		nil,
	)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteAccountFromMod(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodDelete,
		"/accounts/1",
		nil,
	)
	req = SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteAccountFromNormal(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodDelete,
		"/accounts/1",
		nil,
	)
	req = SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
