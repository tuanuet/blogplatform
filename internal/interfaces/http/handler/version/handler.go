package version

import (
	"net/http"
	"strconv"

	"github.com/aiagent/internal/application/dto"
	"github.com/aiagent/internal/domain/repository"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VersionHandler struct {
	versionService service.VersionService
	blogService    service.BlogService
}

func NewVersionHandler(versionService service.VersionService, blogService service.BlogService) *VersionHandler {
	return &VersionHandler{
		versionService: versionService,
		blogService:    blogService,
	}
}

func (h *VersionHandler) List(c *gin.Context) {
	blogID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}
	userID := userIDVal.(uuid.UUID)

	pagination := repository.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	blog, err := h.blogService.GetByID(c.Request.Context(), blogID, &userID)
	if err != nil {
		if err == service.ErrBlogNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == service.ErrBlogAccessDenied {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	if blog.AuthorID != userID {
		response.Forbidden(c, "Access denied")
		return
	}

	result, err := h.versionService.ListVersions(c.Request.Context(), blogID, pagination)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	versions := make([]dto.VersionResponse, len(result.Data))
	for i, v := range result.Data {
		editorBrief := dto.UserBriefResponse{
			ID: v.EditorID,
		}
		if v.Editor != nil {
			editorBrief.Name = v.Editor.Name
			editorBrief.Email = v.Editor.Email
		}

		versions[i] = dto.VersionResponse{
			ID:            v.ID,
			VersionNumber: v.VersionNumber,
			Title:         v.Title,
			Excerpt:       v.Excerpt,
			Status:        string(v.Status),
			Visibility:    string(v.Visibility),
			Editor:        editorBrief,
			ChangeSummary: v.ChangeSummary,
			CreatedAt:     v.CreatedAt,
		}
	}

	response.SuccessWithMeta(c, versions, &response.Meta{
		Page:       result.Page,
		PageSize:   result.PageSize,
		Total:      result.Total,
		TotalPages: result.TotalPages,
	})
}

func (h *VersionHandler) Get(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version ID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}
	userID := userIDVal.(uuid.UUID)

	version, err := h.versionService.GetVersion(c.Request.Context(), versionID)
	if err != nil {
		if err == service.ErrVersionNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	blog, err := h.blogService.GetByID(c.Request.Context(), version.BlogID, &userID)
	if err != nil {
		if err == service.ErrBlogNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == service.ErrBlogAccessDenied {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	if blog.AuthorID != userID {
		response.Forbidden(c, "Access denied")
		return
	}

	editorBrief := dto.UserBriefResponse{
		ID: version.EditorID,
	}
	if version.Editor != nil {
		editorBrief.Name = version.Editor.Name
		editorBrief.Email = version.Editor.Email
	}

	var categoryResp *dto.CategoryResponse
	if version.Category != nil {
		categoryResp = &dto.CategoryResponse{
			ID:   version.Category.ID,
			Name: version.Category.Name,
			Slug: version.Category.Slug,
		}
	}

	tagsResp := make([]dto.TagResponse, 0)
	for _, t := range version.Tags {
		tagsResp = append(tagsResp, dto.TagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	resp := dto.VersionDetailResponse{
		ID:            version.ID,
		VersionNumber: version.VersionNumber,
		Title:         version.Title,
		Slug:          version.Slug,
		Excerpt:       version.Excerpt,
		Content:       version.Content,
		ThumbnailURL:  version.ThumbnailURL,
		Status:        string(version.Status),
		Visibility:    string(version.Visibility),
		Category:      categoryResp,
		Tags:          tagsResp,
		Editor:        editorBrief,
		ChangeSummary: version.ChangeSummary,
		CreatedAt:     version.CreatedAt,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *VersionHandler) Create(c *gin.Context) {
	blogID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	editorIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}
	editorID := editorIDVal.(uuid.UUID)

	var req dto.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// We need to fetch the blog first
	// Note: We use GetByID with nil viewerID because we are snapshotting internal state,
	// checking access might be relevant but usually the service layer handles permission.
	// However, VersionService.CreateVersion does NOT check permissions, it just creates.
	// So we should probably check if the user has access to the blog.
	// But `GetByID` checks access if viewerID is provided.
	// The problem is `GetByID` returns `ErrBlogAccessDenied` if check fails.
	// If the user is the author, they have access.
	// We should pass `editorID` as viewerID.

	blog, err := h.blogService.GetByID(c.Request.Context(), blogID, &editorID)
	if err != nil {
		if err == service.ErrBlogNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == service.ErrBlogAccessDenied {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	// Also ensure user is author. Only author can create versions?
	// Or maybe editors?
	// `VersionService.DeleteVersion` checks if requester is editor or author.
	// Ideally `Create` should also enforce this.
	// `blogService.GetByID` checks if user can view. Viewing doesn't mean editing.
	// But the API requirements didn't specify strict RBAC for creating version.
	// Assuming if they can fetch it (and we might add logic later), we proceed.
	// Actually, usually only Author can update blog, so only Author should create version.
	if blog.AuthorID != editorID {
		response.Forbidden(c, "only author can create versions")
		return
	}

	version, err := h.versionService.CreateVersion(c.Request.Context(), blog, editorID, req.ChangeSummary)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Map response
	editorBrief := dto.UserBriefResponse{
		ID: version.EditorID,
	}
	if version.Editor != nil {
		editorBrief.Name = version.Editor.Name
		editorBrief.Email = version.Editor.Email
	}

	resp := dto.VersionResponse{
		ID:            version.ID,
		VersionNumber: version.VersionNumber,
		Title:         version.Title,
		Excerpt:       version.Excerpt,
		Status:        string(version.Status),
		Visibility:    string(version.Visibility),
		Editor:        editorBrief,
		ChangeSummary: version.ChangeSummary,
		CreatedAt:     version.CreatedAt,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *VersionHandler) Restore(c *gin.Context) {
	blogID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid blog ID")
		return
	}

	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version ID")
		return
	}

	editorIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}
	editorID := editorIDVal.(uuid.UUID)

	blog, err := h.versionService.RestoreVersion(c.Request.Context(), blogID, versionID, editorID)
	if err != nil {
		if err == service.ErrVersionNotFound {
			response.NotFound(c, "version not found")
			return
		}
		if err == service.ErrVersionMismatch {
			response.BadRequest(c, err.Error())
			return
		}
		if err == service.ErrBlogNotFound {
			response.NotFound(c, "blog not found")
			return
		}
		if err == service.ErrBlogAccessDenied {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	// Map Blog to BlogResponse
	// We need to map manually as we don't have a mapper service here
	// and we don't want to duplicate code too much.
	// Ideally we use a helper.

	authorBrief := dto.UserBriefResponse{ID: blog.AuthorID}
	if blog.Author != nil {
		authorBrief.Name = blog.Author.Name
		authorBrief.Email = blog.Author.Email
	}

	var catResp *dto.CategoryResponse
	if blog.Category != nil {
		catResp = &dto.CategoryResponse{
			ID:   blog.Category.ID,
			Name: blog.Category.Name,
			Slug: blog.Category.Slug,
		}
	}

	tagsResp := make([]dto.TagResponse, 0)
	for _, t := range blog.Tags {
		tagsResp = append(tagsResp, dto.TagResponse{
			ID:   t.ID,
			Name: t.Name,
			Slug: t.Slug,
		})
	}

	resp := dto.BlogResponse{
		ID:            blog.ID,
		AuthorID:      blog.AuthorID,
		Author:        &authorBrief,
		CategoryID:    blog.CategoryID,
		Category:      catResp,
		Title:         blog.Title,
		Slug:          blog.Slug,
		Excerpt:       blog.Excerpt,
		Content:       blog.Content,
		ThumbnailURL:  blog.ThumbnailURL,
		Status:        blog.Status,
		Visibility:    blog.Visibility,
		PublishedAt:   blog.PublishedAt,
		Tags:          tagsResp,
		UpvoteCount:   blog.UpvoteCount,
		DownvoteCount: blog.DownvoteCount,
		CreatedAt:     blog.CreatedAt,
		UpdatedAt:     blog.UpdatedAt,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *VersionHandler) Delete(c *gin.Context) {
	versionID, err := uuid.Parse(c.Param("versionId"))
	if err != nil {
		response.BadRequest(c, "invalid version ID")
		return
	}

	requesterIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}
	requesterID := requesterIDVal.(uuid.UUID)

	if err := h.versionService.DeleteVersion(c.Request.Context(), versionID, requesterID); err != nil {
		if err == service.ErrVersionNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == service.ErrBlogAccessDenied {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
