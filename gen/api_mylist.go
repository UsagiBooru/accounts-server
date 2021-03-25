/*
 * UsagiBooru Accounts API
 *
 * アカウント関連API
 *
 * API version: 2.0
 * Contact: dsgamer777@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A MylistApiController binds http requests to an api service and writes the service results to the http response
type MylistApiController struct {
	service MylistApiServicer
}

// NewMylistApiController creates a default api controller
func NewMylistApiController(s MylistApiServicer) Router {
	return &MylistApiController{ service: s }
}

// Routes returns all of the api route for the MylistApiController
func (c *MylistApiController) Routes() Routes {
	return Routes{ 
		{
			"CreateMylist",
			strings.ToUpper("Post"),
			"/accounts/{accountID}/mylists",
			c.CreateMylist,
		},
		{
			"GetUserMylists",
			strings.ToUpper("Get"),
			"/accounts/{accountID}/mylists",
			c.GetUserMylists,
		},
	}
}

// CreateMylist - Create user mylist
func (c *MylistApiController) CreateMylist(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	accountID, err := parseInt32Parameter(params["accountID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	mylistStruct := &MylistStruct{}
	if err := json.NewDecoder(r.Body).Decode(&mylistStruct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	result, err := c.service.CreateMylist(r.Context(), accountID, *mylistStruct)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// GetUserMylists - Get user mylists
func (c *MylistApiController) GetUserMylists(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	accountID, err := parseInt32Parameter(params["accountID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	result, err := c.service.GetUserMylists(r.Context(), accountID)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}
