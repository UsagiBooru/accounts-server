/*
 * UsagiBooru Accounts API
 *
 * Accounts related api (required)
 *
 * API version: 2.0
 * Contact: dsgamer777@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gen

//Response return a ImplResponse struct filled
func Response(code int, body interface{}) ImplResponse {
	return ImplResponse{Code: code, Body: body}
}
