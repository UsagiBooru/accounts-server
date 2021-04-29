package impl_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/models/constmodels"
	"github.com/UsagiBooru/accounts-server/utils/tests"
)

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

func TestCreateAccountBadRequestOnEmptyMainField(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	newAccount := gen.AccountStruct{
		Name:      "デバッグアカウント",
		DisplayID: "debugaccount",
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
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateAccountBadRequestOnEmptyInviteField(t *testing.T) {
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
		Invite:    gen.AccountStructInvite{},
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

func TestEditAccountInternalErrorOnEmptyHeader(t *testing.T) {
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
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEditAccountBadRequestOnInvalidMail(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Mail: "mail_address",
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/1",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEditAccountNotFoundOnInvalidId(t *testing.T) {
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
		"/accounts/404",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestEditAccountBadRequestOnChangeInvite(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Invite: gen.AccountStructInvite{
			InviteID:     999,
			Code:         "newcode",
			InvitedCount: 999,
		},
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/1",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetAdminUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEditAccountForbiddenOnChangeDifferentUserFromNormal(t *testing.T) {
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
	req = tests.SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEditAccountForbiddenOnChangeGreaterUser(t *testing.T) {
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
	req = tests.SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEditAccountForbiddenOnModChangeOwnPermission(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Permission: constmodels.PERMISSION_ADMIN,
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/2",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEditAccountForbiddenOnChangeAccessFromNormal(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Access: gen.AccountStructAccess{
			CanInvite:      true,
			CanLike:        true,
			CanComment:     true,
			CanCreatePost:  true,
			CanEditPost:    true,
			CanApprovePost: true,
		},
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/3",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEditAccountForbiddenOnChangeMailFromMod(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Mail: "mail_change@example.com",
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/3",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestEditAccountConflictOnChangeDisplayId(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		DisplayID: "domao",
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/2",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestEditAccountConflictOnChangeName(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		Name: "ドマオー",
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/2",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestEditAccountBadRequestOnChangeWithWrongPassword(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.AccountStruct{
		OldPassword: "domao",
		Password:    "KafuuChino",
	}
	req_json, _ := json.Marshal(editAccount)
	req := httptest.NewRequest(
		http.MethodPatch,
		"/accounts/2",
		bytes.NewBuffer(req_json),
	)
	req = tests.SetModUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAccountMeInternalErrorOnEmptyHeader(t *testing.T) {
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
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLoginWithFormUnAuthorizedOnInvalidId(t *testing.T) {
	s, shutdown, isParallel := GetAccountsServer()
	if isParallel {
		t.Parallel()
	}
	defer s.Close()
	defer shutdown()
	editAccount := gen.PostLoginWithFormRequest{
		Id:       "omadosan",
		Password: tests.PASSWORD,
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

func TestDeleteAccountForbiddenFromNormal(t *testing.T) {
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
	req = tests.SetNormalUserHeader(req)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	t.Log(rec.Body)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
