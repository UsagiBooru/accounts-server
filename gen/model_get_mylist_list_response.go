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

type GetMylistListResponse struct {

	Pagination PaginationStruct `json:"pagination"`

	Contents []MylistStruct `json:"contents"`
}