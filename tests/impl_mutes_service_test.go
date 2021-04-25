package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/impl"
	"github.com/UsagiBooru/accounts-server/utils/server"
)

func GetMutesServer() (*httptest.Server, func(), bool) {
	db, shutdown, isParallel := GetDatabaseConnection()
	MutesApiService := impl.NewMutesApiImplService(db)
	MutesApiController := gen.NewMutesApiController(MutesApiService)
	router := server.NewRouterWithInject(MutesApiController)
	return httptest.NewServer(router), shutdown, isParallel
}

func TestGetMuteSuccess(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1/mutes/1", nil)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateMuteSuccess(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newMute := gen.MuteStruct{
		AccountID:  1,
		TargetType: "artist",
		TargetID:   2,
	}
	user_json, err := json.Marshal(newMute)
	if err != nil {
		log.Fatal("Convert struct to json failed.")
	}
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/1/mutes",
		bytes.NewBuffer(user_json),
	)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteMuteSuccess(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodDelete,
		"/accounts/1/mutes/1",
		nil,
	)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
