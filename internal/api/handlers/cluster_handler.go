package handlers

import (
	"net/http"
	"strconv"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/mappers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	service services.ClusterService
	repo    repositories.ClusterRepository
	// TODO
	// eventService *events.EventService
}

type clusterFilterParams struct {
	Status   string `form:"status"`
	Provider string `form:"provider"`
	Region   string `form:"region"`
	Account  string `form:"account"`
}

type listClustersRequest struct {
	dto.PaginationRequest
	Filters clusterFilterParams `form:"inline"`
}

func (f *clusterFilterParams) ToRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.Status != "" {
		filters["status"] = f.Status
	}
	if f.Provider != "" {
		filters["provider"] = f.Provider
	}
	if f.Region != "" {
		filters["region"] = f.Region
	}
	if f.Account != "" {
		filters["account_name"] = f.Account
	}
	return filters
}

func NewClusterHandler(service services.ClusterService, repo repositories.ClusterRepository) *ClusterHandler {
	return &ClusterHandler{service: service, repo: repo}
}

// ==================== Clusters Handlers ====================

// ListClusters handles the request for obtaining the entire Cluster list
//
//	@Summary		Obtain every Cluster
//	@Description	Returns a list of Clusters with a single instance filtered by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ClusterListResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/clusters [get]
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	var req listClustersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		// TODO maybe jsut "Invalid format for query parameters. Please check the documentation."
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.ToRepoFilters(),
	}

	clusters, total, err := h.repo.ListClusters(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve clusters"))
		return
	}

	clusterDTOs := mappers.ToClusterDTOs(clusters)
	response := dto.NewListResponse(clusterDTOs, total)
	// a.logger.Debug("Retrieving complete clusters inventory")

	//clusters, err := h.repo.GetClusters()
	//if err != nil {
	//	a.logger.Error("Can't retrieve Clusters list", zap.Error(err))
	//	c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
	//	return
	//}
	//
	//c.PureJSON(http.StatusOK, NewClusterListResponse(clusters))
	// TODO. REview, should be a good option to avoid JSON body unpack to get numbers
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.PureJSON(http.StatusOK, response)
}
