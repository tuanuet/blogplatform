# Feature: Blog Version History

## Objective
Cho phép authors lưu lịch sử chỉnh sửa bài viết, xem và khôi phục về các phiên bản cũ.

## Requirements

### Functional Requirements
- [ ] Auto-save version khi create blog (Version 1 - initial state)
- [ ] Auto-save version mỗi khi blog được update (post-update snapshot)
- [ ] Manual-save version khi author muốn tạo checkpoint (POST /versions)
- [ ] List tất cả versions của 1 blog (phân trang)
- [ ] Xem chi tiết 1 version (full snapshot)
- [ ] Restore blog về version cũ (tạo version mới từ version cũ)
- [ ] Xóa version (chỉ author hoặc admin)
- [ ] Giới hạn 50 versions/blog (configurable)

### Non-Functional Requirements
- [ ] Versions được lưu trong bảng riêng, không ảnh hưởng performance query blog chính
- [ ] Auto-cleanup khi vượt quá 50 versions (xóa version cũ nhất)
- [ ] Audit trail: ai edit, khi nào, change summary

## Technical Context

### New Database Schema

```sql
-- Blog versions table
CREATE TABLE blog_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blog_id UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    version_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    excerpt TEXT,
    content TEXT NOT NULL,
    thumbnail_url VARCHAR(500),
    status blog_status NOT NULL,
    visibility blog_visibility NOT NULL,
    category_id UUID REFERENCES categories(id),
    editor_id UUID NOT NULL REFERENCES users(id),
    change_summary TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(blog_id, version_number)
);

-- Index for performance
CREATE INDEX idx_blog_versions_blog_id ON blog_versions(blog_id);
CREATE INDEX idx_blog_versions_blog_id_created_at ON blog_versions(blog_id, created_at DESC);

-- Blog version tags (many-to-many)
CREATE TABLE blog_version_tags (
    version_id UUID NOT NULL REFERENCES blog_versions(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (version_id, tag_id)
);
```

### New Entities

```go
// internal/domain/entity/blog_version.go
type BlogVersion struct {
    ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    BlogID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"blogId"`
    VersionNumber int            `gorm:"not null" json:"versionNumber"`
    Title         string         `gorm:"size:255;not null" json:"title"`
    Slug          string         `gorm:"size:255;not null" json:"slug"`
    Excerpt       *string        `gorm:"type:text" json:"excerpt,omitempty"`
    Content       string         `gorm:"type:text;not null" json:"content"`
    ThumbnailURL  *string        `gorm:"size:500" json:"thumbnailUrl,omitempty"`
    Status        BlogStatus     `gorm:"type:blog_status;not null" json:"status"`
    Visibility    BlogVisibility `gorm:"type:blog_visibility;not null" json:"visibility"`
    CategoryID    *uuid.UUID     `gorm:"type:uuid" json:"categoryId,omitempty"`
    EditorID      uuid.UUID      `gorm:"type:uuid;not null" json:"editorId"`
    ChangeSummary *string        `gorm:"type:text" json:"changeSummary,omitempty"`
    CreatedAt     time.Time      `gorm:"not null;default:now()" json:"createdAt"`
    
    // Relationships
    Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    Editor   *User     `gorm:"foreignKey:EditorID" json:"editor,omitempty"`
    Tags     []Tag     `gorm:"many2many:blog_version_tags" json:"tags,omitempty"`
}

func (BlogVersion) TableName() string {
    return "blog_versions"
}
```

### New Repository Interface

```go
// internal/domain/repository/blog_version_repository.go
type BlogVersionRepository interface {
    Create(ctx context.Context, version *entity.BlogVersion) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.BlogVersion, error)
    FindByBlogID(ctx context.Context, blogID uuid.UUID, pagination Pagination) (*PaginatedResult[entity.BlogVersion], error)
    GetNextVersionNumber(ctx context.Context, blogID uuid.UUID) (int, error)
    Delete(ctx context.Context, id uuid.UUID) error
    CountByBlogID(ctx context.Context, blogID uuid.UUID) (int64, error)
    DeleteOldest(ctx context.Context, blogID uuid.UUID, keep int) error
}
```

### New Service Interface

```go
// internal/domain/service/version_service.go
type VersionService interface {
    CreateVersion(ctx context.Context, blog *entity.Blog, editorID uuid.UUID, changeSummary string) (*entity.BlogVersion, error)
    ListVersions(ctx context.Context, blogID uuid.UUID, pagination repository.Pagination) (*repository.PaginatedResult[entity.BlogVersion], error)
    GetVersion(ctx context.Context, versionID uuid.UUID) (*entity.BlogVersion, error)
    RestoreVersion(ctx context.Context, blogID, versionID, editorID uuid.UUID) (*entity.Blog, error)
    DeleteVersion(ctx context.Context, versionID uuid.UUID, requesterID uuid.UUID) error
}
```

### API Endpoints

```
GET    /api/v1/blogs/:id/versions              - List versions
GET    /api/v1/blogs/:id/versions/:versionId   - Get version detail  
POST   /api/v1/blogs/:id/versions              - Manual checkpoint
POST   /api/v1/blogs/:id/versions/:versionId/restore  - Restore
DELETE /api/v1/blogs/:id/versions/:versionId   - Delete version
```

### DTOs

```go
// internal/application/dto/version.go

// CreateVersionRequest - Manual checkpoint
 type CreateVersionRequest struct {
    ChangeSummary string `json:"changeSummary,omitempty"`
}

// VersionResponse
 type VersionResponse struct {
    ID            uuid.UUID          `json:"id"`
    VersionNumber int                `json:"versionNumber"`
    Title         string             `json:"title"`
    Excerpt       *string            `json:"excerpt,omitempty"`
    Status        string             `json:"status"`
    Visibility    string             `json:"visibility"`
    Editor        UserSummary        `json:"editor"`
    ChangeSummary *string            `json:"changeSummary,omitempty"`
    CreatedAt     time.Time          `json:"createdAt"`
}

// VersionDetailResponse - Full content
 type VersionDetailResponse struct {
    ID            uuid.UUID          `json:"id"`
    VersionNumber int                `json:"versionNumber"`
    Title         string             `json:"title"`
    Slug          string             `json:"slug"`
    Excerpt       *string            `json:"excerpt,omitempty"`
    Content       string             `json:"content"`
    ThumbnailURL  *string            `json:"thumbnailUrl,omitempty"`
    Status        string             `json:"status"`
    Visibility    string             `json:"visibility"`
    Category      *CategoryResponse  `json:"category,omitempty"`
    Tags          []TagResponse      `json:"tags,omitempty"`
    Editor        UserSummary        `json:"editor"`
    ChangeSummary *string            `json:"changeSummary,omitempty"`
    CreatedAt     time.Time          `json:"createdAt"`
}

// VersionListResponse
 type VersionListResponse struct {
    Data       []VersionResponse `json:"data"`
    Total      int64             `json:"total"`
    Page       int               `json:"page"`
    PageSize   int               `json:"pageSize"`
    TotalPages int               `json:"totalPages"`
}
```

### Impacted Services

1. **BlogService**
   - `Create()` - Thêm: tạo version 1 sau khi create blog
   - `Update()` - Thêm: auto-save version sau khi update

2. **New: VersionService**
   - `CreateVersion()` - Tạo version mới
   - `ListVersions()` - List versions theo blog
   - `GetVersion()` - Get version detail
   - `RestoreVersion()` - Restore blog từ version
   - `DeleteVersion()` - Xóa version

### Implementation Flow

**1. Create Blog:**
```
BlogService.Create(blog)
  → blogRepo.Create(blog)
  → versionService.CreateVersion(blog, authorID, "Initial version")
  → Kiểm tra limit 50 versions
  → Nếu > 50: xóa versions cũ nhất
```

**2. Update Blog:**
```
BlogService.Update(blog)
  → Lấy blog hiện tại
  → blogRepo.Update(blog)
  → versionService.CreateVersion(updatedBlog, editorID, "Auto-saved")
  → Kiểm tra limit 50 versions
```

**3. Manual Checkpoint:**
```
POST /blogs/:id/versions
  → Handler gọi versionService.CreateVersion(blog, currentUserID, request.ChangeSummary)
```

**4. Restore Version:**
```
POST /blogs/:id/versions/:vid/restore
  → Handler lấy version
  → Tạo blog mới từ version data
  → BlogService.Update(newBlogData)
  → Tự động tạo version mới ("Restored from version X")
```

## Acceptance Criteria

1. **Auto-save on create:** Khi tạo blog mới, tự động tạo version 1
2. **Auto-save on update:** Mỗi lần update blog, tự động tạo version mới
3. **Manual checkpoint:** Author có thể tạo version thủ công qua API
4. **List versions:** Có thể xem danh sách versions với phân trang
5. **View version detail:** Có thể xem chi tiết 1 version (full content)
6. **Restore:** Có thể restore blog về bất kỳ version nào (tạo version mới từ version cũ)
7. **Delete:** Author có thể xóa version của mình, admin có thể xóa mọi version
8. **Limit enforcement:** System không lưu quá 50 versions/blog (auto-cleanup)
9. **Audit trail:** Mỗi version lưu editor_id và thời gian

## Constraints

- Chỉ lưu versions cho blogs đã tồn tại (blog_id phải valid)
- Không cho phép edit version đã tạo (immutable)
- Restore tạo version mới, không ghi đè versions cũ
- Giới hạn 50 versions có thể configurable qua config

## Notes for Implementation

1. **Database Migration:** Cần tạo migration cho 2 bảng mới
2. **Auto-save logic:** Thêm vào BlogService.Create và BlogService.Update
3. **Permission check:** 
   - List/Get versions: Author hoặc Admin
   - Create version: Author hoặc Collaborator (nếu có)
   - Restore: Author hoặc Admin
   - Delete: Author hoặc Admin
4. **Change Summary:** Optional, cho phép null/empty
5. **Tags:** Khi restore, cần restore cả tags của version đó
