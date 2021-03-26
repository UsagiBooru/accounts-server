package api_tester

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	openapi "github.com/UsagiBooru/accounts-server/gen"
	impl "github.com/UsagiBooru/accounts-server/impl"
)

func GetAccountsServer() *httptest.Server {
	AccountsApiService := impl.NewAccountsApiImplService()
	AccountsApiController := openapi.NewAccountsApiController(AccountsApiService)
	router := openapi.NewRouter(AccountsApiController)
	return httptest.NewServer(router)
}

func TestGetAccount(t *testing.T) {
	s := GetAccountsServer()
	defer s.Close()
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	rec := httptest.NewRecorder()
	s.Config.Handler.ServeHTTP(rec, req)
	// t.Log(rec.Body)
	assert.Equal(t, http.StatusOK, rec.Code)
}
