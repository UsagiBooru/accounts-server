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

// TODO: Fix flag. This is not working properly. (Always parallel)
var parallelFlag = flag.Bool("parallel", true, "Set true to use parallel test(Local), otherwise to simple test(CI)")

func GetAccountsServer() (*httptest.Server, func(), bool) {
	var db *mongo.Client
	var shutdown func()
	var err error
	var isParallel bool
	if *parallelFlag {
		db, shutdown, err = GenerateMongoTestContainer()
		isParallel = true
	} else {
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

func SetAdminUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "1")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(request.PermissionAdmin))
	return req
}

func SetModUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "2")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(request.PermissionModerator))
	return req
}

func SetNormalUserHeader(req *http.Request) *http.Request {
	req.Header.Set("x-consumer-user-id", "3")
	req.Header.Set("x-consumer-user-permission", strconv.Itoa(request.PermissionUser))
	return req
}

func TestMain(m *testing.M) {
	// server.Debug("Resetting database...")
	err := ReGenerateTestDatabase()
	if err != nil {
		server.Error(err.Error())
	}
	// server.Debug("Reset database success.")

	m.Run()
}

func TestGetAccount(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateAccount(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
	newAccount := gen.AccountStruct{
		Name:      "デバッグアカウント",
		DisplayID: "debugaccount",
		Password:  "debugaccount",
		Mail:      "mail@example.com",
		Invite: gen.AccountStructInvite{
			Code: "devcode1",
		},
	}
	user_json, err := json.Marshal(newAccount)
	if err != nil {
		log.Fatal("Convert struct to json failed.")
	}
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

func TestEditAccount(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
	editAccount := gen.AccountStruct{
		Name: "デバッグアカウント2",
	}
	req_json, err := json.Marshal(editAccount)
	if err != nil {
		log.Fatal("Convert struct to json failed.")
	}
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

func TestGetAccountMe(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
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

func TestLoginWithForm(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
	editAccount := gen.PostLoginWithFormRequest{
		Id:       "domao",
		Password: PASSWORD,
	}
	req_json, err := json.Marshal(editAccount)
	if err != nil {
		log.Fatal("Convert struct to json failed.")
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/login/form",
		bytes.NewBuffer(req_json),
	)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteAccount(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
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
