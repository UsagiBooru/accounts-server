package impl_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/utils/tests"
)

func TestAddMuteBadRequestOnInvalidStruct(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newMute := gen.MuteStruct{
		TargetType: "artist",
		TargetID:   -1204,
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
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddMuteInternalErrorOnEmptyHeader(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newMute := gen.MuteStruct{
		TargetType: "artist",
		TargetID:   2,
	}
	user_json, _ := json.Marshal(newMute)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/1/mutes",
		bytes.NewBuffer(user_json),
	)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAddMuteForbiddenOnAccessOtherFromNormal(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newMute := gen.MuteStruct{
		TargetType: "artist",
		TargetID:   2,
	}
	user_json, _ := json.Marshal(newMute)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/1/mutes",
		bytes.NewBuffer(user_json),
	)
	req = tests.SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAddMuteNotFoundOnInvalidId(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newMute := gen.MuteStruct{
		TargetType: "artist",
		TargetID:   2,
	}
	user_json, _ := json.Marshal(newMute)
	req := httptest.NewRequest(
		http.MethodPost,
		"/accounts/404/mutes",
		bytes.NewBuffer(user_json),
	)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestAddMuteConflictedOnExistedMute(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newMute := gen.MuteStruct{
		TargetType: "artist",
		TargetID:   1,
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
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestDeleteMuteOnEmptyHeader(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodDelete, "/accounts/1/mutes/1", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetMuteInternalErrorOnEmptyHeader(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1/mutes/404", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteMuteForbiddenOnAccessOtherFromNormal(t *testing.T) {
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
	req = tests.SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteMuteNotFound(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(
		http.MethodDelete,
		"/accounts/1/mutes/404",
		nil,
	)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetMuteForbiddenOnAccessOtherFromNormal(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1/mutes/1", nil)
	req = tests.SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetMuteNotFoundOnInvalidId(t *testing.T) {
	s, shutdown, isParallel := GetMutesServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1/mutes/404", nil)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
