package impl_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/impl"
	"github.com/UsagiBooru/accounts-server/utils/server"
	"github.com/UsagiBooru/accounts-server/utils/tests"
)

func GetMutesServer() (*httptest.Server, func(), bool) {
	db, shutdown, isParallel := tests.GetDatabaseConnection()
	MutesApiService := impl.NewMutesApiImplService(db)
	MutesApiController := gen.NewMutesApiController(MutesApiService)
	router := server.NewRouterWithInject(MutesApiController)
	return httptest.NewServer(router), shutdown, isParallel
}

func TestAddMuteSuccess(t *testing.T) {
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
	user_json, _ := json.Marshal(newMute)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/1/mutes",
		bytes.NewBuffer(user_json),
	)
	req = tests.SetAdminUserHeader(req)
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
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestGetMuteSuccess(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1/mutes/1", nil)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}
