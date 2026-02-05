# Blog System Design

## 1. Overview

Hệ thống **Blog** cung cấp nền tảng cho người dùng tạo, quản lý và chia sẻ nội dung bài viết. Nó hỗ trợ đầy đủ quy trình vòng đời của một bài viết từ lúc soạn thảo, xuất bản cho đến khi tương tác với cộng đồng.

## 2. Architecture (Kiến trúc)

Hệ thống tuân theo kiến trúc **Clean Architecture** (Layered Architecture):

- **UseCase Layer** (`internal/application/usecase/blog`):
  - `BlogUseCase`: Điều phối logic nghiệp vụ chính, tương tác với `BlogService`.
- **Domain Service Layer** (`internal/domain/service`):
  - `BlogService`: Chứa logic cốt lõi (Business Rules) như kiểm tra quyền sở hữu, validate slug, xử lý logic xuất bản.
- **Domain Entity Layer** (`internal/domain/entity`):
  - `Blog`: Entity chính đại diện cho bài viết.
  - `BlogReaction`, `Category`, `Tag`: Các entity liên quan.
- **Repository Layer** (`internal/domain/repository`):
  - `BlogRepository`: Interface tương tác với cơ sở dữ liệu.

## 3. Core Features & Logic

### A. Blog Management (CRUD)
- **Create**: Tạo bài viết mới với trạng thái mặc định là `Draft`. Hỗ trợ gắn thẻ (Tags) và danh mục (Category).
- **Update**: Cập nhật nội dung, tiêu đề, thumbnail, v.v. Chỉ tác giả (Author) mới có quyền cập nhật.
- **Delete**: Xóa bài viết (Soft delete).
- **Read**:
  - Xem chi tiết theo `ID` hoặc `Slug`.
  - Hỗ trợ `ViewerID` để cá nhân hóa nội dung (ví dụ: hiển thị trạng thái đã like/dislike của người xem).

### B. Publishing Flow (Quy trình xuất bản)
- **Status**: Bài viết có các trạng thái: `Draft`, `Published`, `Archived`.
- **Scheduling**: Hỗ trợ hẹn giờ xuất bản thông qua trường `PublishedAt`.
- **Visibility**:
  - `Public`: Ai cũng có thể xem.
  - `Private`: Chỉ tác giả mới thấy (hoặc nhóm được cấp quyền).

### C. Engagement (Tương tác)
- **Reactions**: Người dùng có thể `Upvote` hoặc `Downvote` bài viết.
- **Logic**:
  - Mỗi người dùng chỉ có 1 trạng thái reaction duy nhất cho 1 bài viết tại 1 thời điểm.
  - Hệ thống tính toán tổng `UpvoteCount` và `DownvoteCount`.

### D. Organization & Discovery
- **Categorization**: Bài viết thuộc về 1 `Category`.
- **Tagging**: Bài viết có thể gắn nhiều `Tag` để dễ dàng tìm kiếm.
- **Search & Filtering**:
  - Tìm kiếm theo từ khóa (Title, Content).
  - Lọc theo Author, Category, Status.
  - Lọc theo thời gian xuất bản (ẩn các bài hẹn giờ chưa đến giờ G đối với người xem thường).

## 4. Data Model

### Entities chính

| Entity | Mô tả |
| :--- | :--- |
| `Blog` | Chứa thông tin chính: Title, Content, Slug, Status, AuthorID, CategoryID, Stats (Upvote/Downvote). |
| `BlogReaction` | Lưu trạng thái tương tác của User với Blog (Up/Down). |
| `Category` | Danh mục bài viết (1-n với Blog). |
| `Tag` | Thẻ bài viết (n-n với Blog). |

### Trạng thái (Enum)

- **BlogStatus**: `draft`, `published`, `archived`.
- **BlogVisibility**: `public`, `private`, `internal`.

## 5. API Reference (Internal Design)

### Public / Protected Operations (via UseCase)

- `Create(ctx, authorID, req)`
- `GetByID(ctx, id, viewerID)`
- `GetBySlug(ctx, authorID, slug, viewerID)`
- `List(ctx, params, viewerID)`
- `Update(ctx, id, authorID, req)`
- `Delete(ctx, id, authorID)`
- `Publish(ctx, id, authorID, req)`
- `Unpublish(ctx, id, authorID)`
- `React(ctx, id, userID, req)`
