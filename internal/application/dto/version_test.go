package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateVersionRequest_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		request  CreateVersionRequest
		expected string
	}{
		{
			name:     "empty request",
			request:  CreateVersionRequest{},
			expected: `{}`,
		},
		{
			name: "with change summary",
			request: CreateVersionRequest{
				ChangeSummary: "Updated introduction",
			},
			expected: `{"changeSummary":"Updated introduction"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))
		})
	}
}

func TestCreateVersionRequest_JSONDeserialization(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected CreateVersionRequest
	}{
		{
			name:     "empty JSON",
			json:     `{}`,
			expected: CreateVersionRequest{},
		},
		{
			name: "with change summary",
			json: `{"changeSummary":"Updated introduction"}`,
			expected: CreateVersionRequest{
				ChangeSummary: "Updated introduction",
			},
		},
		{
			name: "with null change summary",
			json: `{"changeSummary":null}`,
			expected: CreateVersionRequest{
				ChangeSummary: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result CreateVersionRequest
			err := json.Unmarshal([]byte(tt.json), &result)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionResponse_JSONSerialization(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	editorID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")
	createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	excerpt := "Test excerpt"
	changeSummary := "Initial version"

	response := VersionResponse{
		ID:            id,
		VersionNumber: 1,
		Title:         "Test Title",
		Excerpt:       &excerpt,
		Status:        "published",
		Visibility:    "public",
		Editor: UserBriefResponse{
			ID:    editorID,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		ChangeSummary: &changeSummary,
		CreatedAt:     createdAt,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	// Verify all expected fields are present
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, id.String(), result["id"])
	assert.Equal(t, float64(1), result["versionNumber"])
	assert.Equal(t, "Test Title", result["title"])
	assert.Equal(t, "Test excerpt", result["excerpt"])
	assert.Equal(t, "published", result["status"])
	assert.Equal(t, "public", result["visibility"])
	assert.NotNil(t, result["editor"])
	assert.Equal(t, "Initial version", result["changeSummary"])
	assert.NotNil(t, result["createdAt"])
}

func TestVersionResponse_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"versionNumber": 1,
		"title": "Test Title",
		"excerpt": "Test excerpt",
		"status": "published",
		"visibility": "public",
		"editor": {
			"id": "660e8400-e29b-41d4-a716-446655440001",
			"name": "John Doe",
			"email": "john@example.com"
		},
		"changeSummary": "Initial version",
		"createdAt": "2024-01-01T12:00:00Z"
	}`

	var result VersionResponse
	err := json.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), result.ID)
	assert.Equal(t, 1, result.VersionNumber)
	assert.Equal(t, "Test Title", result.Title)
	assert.NotNil(t, result.Excerpt)
	assert.Equal(t, "Test excerpt", *result.Excerpt)
	assert.Equal(t, "published", result.Status)
	assert.Equal(t, "public", result.Visibility)
	assert.Equal(t, uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"), result.Editor.ID)
	assert.Equal(t, "John Doe", result.Editor.Name)
	assert.NotNil(t, result.ChangeSummary)
	assert.Equal(t, "Initial version", *result.ChangeSummary)
}

func TestVersionResponse_OptionalFields(t *testing.T) {
	// Test with nil optional fields
	jsonData := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"versionNumber": 1,
		"title": "Test Title",
		"status": "draft",
		"visibility": "private",
		"editor": {
			"id": "660e8400-e29b-41d4-a716-446655440001",
			"name": "John Doe",
			"email": "john@example.com"
		},
		"createdAt": "2024-01-01T12:00:00Z"
	}`

	var result VersionResponse
	err := json.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	assert.Nil(t, result.Excerpt)
	assert.Nil(t, result.ChangeSummary)
}

func TestVersionDetailResponse_JSONSerialization(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	editorID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")
	categoryID := uuid.MustParse("770e8400-e29b-41d4-a716-446655440002")
	tagID := uuid.MustParse("880e8400-e29b-41d4-a716-446655440003")
	createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	thumbnailURL := "https://example.com/image.jpg"
	excerpt := "Test excerpt"
	changeSummary := "Initial version"

	response := VersionDetailResponse{
		ID:            id,
		VersionNumber: 1,
		Title:         "Test Title",
		Slug:          "test-title",
		Excerpt:       &excerpt,
		Content:       "Test content here",
		ThumbnailURL:  &thumbnailURL,
		Status:        "published",
		Visibility:    "public",
		Category: &CategoryResponse{
			ID:          categoryID,
			Name:        "Tech",
			Slug:        "tech",
			Description: strPtr("Technology articles"),
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		},
		Tags: []TagResponse{
			{
				ID:        tagID,
				Name:      "Go",
				Slug:      "go",
				CreatedAt: createdAt,
			},
		},
		Editor: UserBriefResponse{
			ID:    editorID,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		ChangeSummary: &changeSummary,
		CreatedAt:     createdAt,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, id.String(), result["id"])
	assert.Equal(t, float64(1), result["versionNumber"])
	assert.Equal(t, "Test Title", result["title"])
	assert.Equal(t, "test-title", result["slug"])
	assert.Equal(t, "Test excerpt", result["excerpt"])
	assert.Equal(t, "Test content here", result["content"])
	assert.Equal(t, "https://example.com/image.jpg", result["thumbnailUrl"])
	assert.Equal(t, "published", result["status"])
	assert.Equal(t, "public", result["visibility"])
	assert.NotNil(t, result["category"])
	assert.NotNil(t, result["tags"])
	assert.NotNil(t, result["editor"])
	assert.Equal(t, "Initial version", result["changeSummary"])
	assert.NotNil(t, result["createdAt"])
}

func TestVersionDetailResponse_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"versionNumber": 1,
		"title": "Test Title",
		"slug": "test-title",
		"excerpt": "Test excerpt",
		"content": "Test content here",
		"thumbnailUrl": "https://example.com/image.jpg",
		"status": "published",
		"visibility": "public",
		"category": {
			"id": "770e8400-e29b-41d4-a716-446655440002",
			"name": "Tech",
			"slug": "tech",
			"createdAt": "2024-01-01T12:00:00Z",
			"updatedAt": "2024-01-01T12:00:00Z"
		},
		"tags": [
			{
				"id": "880e8400-e29b-41d4-a716-446655440003",
				"name": "Go",
				"slug": "go",
				"createdAt": "2024-01-01T12:00:00Z"
			}
		],
		"editor": {
			"id": "660e8400-e29b-41d4-a716-446655440001",
			"name": "John Doe",
			"email": "john@example.com"
		},
		"changeSummary": "Initial version",
		"createdAt": "2024-01-01T12:00:00Z"
	}`

	var result VersionDetailResponse
	err := json.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), result.ID)
	assert.Equal(t, 1, result.VersionNumber)
	assert.Equal(t, "Test Title", result.Title)
	assert.Equal(t, "test-title", result.Slug)
	assert.NotNil(t, result.Excerpt)
	assert.Equal(t, "Test content here", result.Content)
	assert.NotNil(t, result.ThumbnailURL)
	assert.Equal(t, "https://example.com/image.jpg", *result.ThumbnailURL)
	assert.NotNil(t, result.Category)
	assert.Equal(t, "Tech", result.Category.Name)
	assert.Len(t, result.Tags, 1)
	assert.Equal(t, "Go", result.Tags[0].Name)
}

func TestVersionListResponse_JSONSerialization(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	editorID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")
	createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	response := VersionListResponse{
		Data: []VersionResponse{
			{
				ID:            id,
				VersionNumber: 1,
				Title:         "Test Title",
				Status:        "published",
				Visibility:    "public",
				Editor: UserBriefResponse{
					ID:    editorID,
					Name:  "John Doe",
					Email: "john@example.com",
				},
				CreatedAt: createdAt,
			},
		},
		Total:      100,
		Page:       1,
		PageSize:   10,
		TotalPages: 10,
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.NotNil(t, result["data"])
	assert.Equal(t, float64(100), result["total"])
	assert.Equal(t, float64(1), result["page"])
	assert.Equal(t, float64(10), result["pageSize"])
	assert.Equal(t, float64(10), result["totalPages"])
}

func TestVersionListResponse_JSONDeserialization(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"versionNumber": 1,
				"title": "Test Title",
				"status": "published",
				"visibility": "public",
				"editor": {
					"id": "660e8400-e29b-41d4-a716-446655440001",
					"name": "John Doe",
					"email": "john@example.com"
				},
				"createdAt": "2024-01-01T12:00:00Z"
			}
		],
		"total": 100,
		"page": 1,
		"pageSize": 10,
		"totalPages": 10
	}`

	var result VersionListResponse
	err := json.Unmarshal([]byte(jsonData), &result)
	require.NoError(t, err)

	assert.Len(t, result.Data, 1)
	assert.Equal(t, int64(100), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 10, result.TotalPages)
	assert.Equal(t, "Test Title", result.Data[0].Title)
}
