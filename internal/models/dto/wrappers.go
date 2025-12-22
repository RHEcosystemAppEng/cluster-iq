package dto

import (
	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
)

type AccountListResponse struct {
	responsetypes.ListResponse[AccountDTOResponse]
} // @name AccountListResponse

type ClusterListResponse struct {
	responsetypes.ListResponse[ClusterDTOResponse]
} // @name ClusterListResponse

type ActionListResponse struct {
	responsetypes.ListResponse[ActionDTOResponse]
} // @name ActionListResponse

type InstanceListResponse struct {
	responsetypes.ListResponse[InstanceDTOResponse]
} // @name InstanceListResponse

type ExpenseListResponse struct {
	responsetypes.ListResponse[ExpenseDTOResponse]
} // @name ExpenseListResponse

type SystemEventListResponse struct {
	responsetypes.ListResponse[SystemEventDTOResponse]
} // @name SystemEventListResponse

type ClusterEventListResponse struct {
	responsetypes.ListResponse[ClusterEventDTOResponse]
} // @name ClusterEventListResponse
