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

func GetMutesServer() *httptest.Server {
	MutesApiService := impl.NewMutesApiImplService()
	MutesApiController := gen.NewMutesApiController(MutesApiService)
	router := server.NewRouterWithInject(MutesApiController)
	return httptest.NewServer(router)
}

func TestGetMute(t *testing.T) {
	s := GetMutesServer()
	defer s.Close()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1/mutes/1", nil)
	req = SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateMute(t *testing.T) {
	s := GetMutesServer()
	defer s.Close()
	newMute := gen.MuteStruct{
		MuteID:     1,
		AccountID:  1,
		TargetType: "artist",
		TargetID:   1,
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

func TestDeleteMute(t *testing.T) {
	s := GetMutesServer()
	defer s.Close()
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
