package dto

import (
	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
)

// AccountListResponse wraps a paginated list of AccountDTOResponse items
// for OpenAPI schema generation.
type AccountListResponse struct {
	responsetypes.ListResponse[AccountDTOResponse]
} // @name AccountListResponse

// ClusterListResponse wraps a paginated list of ClusterDTOResponse items
// for OpenAPI schema generation.
type ClusterListResponse struct {
	responsetypes.ListResponse[ClusterDTOResponse]
} // @name ClusterListResponse

// ActionListResponse wraps a paginated list of ActionDTOResponse items
// for OpenAPI schema generation.
type ActionListResponse struct {
	responsetypes.ListResponse[ActionDTOResponse]
} // @name ActionListResponse

// InstanceListResponse wraps a paginated list of InstanceDTOResponse items
// for OpenAPI schema generation.
type InstanceListResponse struct {
	responsetypes.ListResponse[InstanceDTOResponse]
} // @name InstanceListResponse

// ExpenseListResponse wraps a paginated list of ExpenseDTOResponse items
// for OpenAPI schema generation.
type ExpenseListResponse struct {
	responsetypes.ListResponse[ExpenseDTOResponse]
} // @name ExpenseListResponse

// SystemEventListResponse wraps a paginated list of SystemEventDTOResponse items
// for OpenAPI schema generation.
type SystemEventListResponse struct {
	responsetypes.ListResponse[SystemEventDTOResponse]
} // @name SystemEventListResponse

// ClusterEventListResponse wraps a paginated list of ClusterEventDTOResponse items
// for OpenAPI schema generation.
type ClusterEventListResponse struct {
	responsetypes.ListResponse[ClusterEventDTOResponse]
} // @name ClusterEventListResponse
