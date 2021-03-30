package api_tester

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UsagiBooru/accounts-server/gen"
	"github.com/UsagiBooru/accounts-server/impl"
)

func GetAccountsServer() *httptest.Server {
	AccountsApiService := impl.NewAccountsApiImplService()
	AccountsApiController := gen.NewAccountsApiController(AccountsApiService)
	router := gen.NewRouter(AccountsApiController)
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
