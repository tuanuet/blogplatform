# Ranking System Design

## 1. Overview

Hệ thống **Ranking** (Xếp hạng) được thiết kế để định danh và vinh danh những người dùng có hoạt động sôi nổi và tăng trưởng tốt trong cộng đồng. Thay vì chỉ xếp hạng dựa trên tổng số lượng người theo dõi (Followers) tích lũy - điều có thể gây bất lợi cho những người dùng mới nhưng tài năng - hệ thống sử dụng mô hình **"Velocity Score"** (Điểm tốc độ).

Mô hình này tập trung vào **tăng trưởng** và **hoạt động** trong cửa sổ thời gian **30 ngày gần nhất**.

## 2. Architecture (Kiến trúc)

Hệ thống tuân theo kiến trúc **Clean Architecture** (Layered Architecture):

- **Handler Layer** (`internal/interfaces/http/handler/ranking`):
  - Expose API endpoints cho Client (`/ranking/top`, `/ranking/me`).
- **UseCase Layer** (`internal/application/usecase/ranking`):
  - `RankingUseCase`: Điều phối logic nghiệp vụ, gọi xuống Service.
- **Domain Service Layer** (`internal/domain/service`):
  - `RankingService`: Chứa logic cốt lõi tính toán điểm số (Business Rules).
  - `RankingJob`: Quản lý các tiến trình chạy ngầm (Background Jobs).
- **Repository Layer**:
  - Tương tác với cơ sở dữ liệu (PostgreSQL).
- **Infrastructure**:
  - Cron Job / Scheduler: Kích hoạt việc tính toán lại bảng xếp hạng định kỳ.

## 3. Core Logic: Velocity Score

Điểm số xếp hạng (**Composite Score**) được tính toán dựa trên công thức có trọng số, đánh giá 2 yếu tố chính trong 30 ngày:

### Công thức tổng quát
```
CompositeScore = (FollowerGrowthRate * 0.6) + (BlogPostVelocity * 0.4)
```

### Chi tiết các thành phần

#### A. Follower Growth Rate (Trọng số 0.6)
Đo lường tốc độ tăng trưởng người theo dõi.
- **Logic**: So sánh `CurrentFollowers` với `Followers30DaysAgo`.
- **Normalization (Chuẩn hóa)**:
  - Để đảm bảo công bằng cho các tài khoản nhỏ (ví dụ: tăng từ 1 lên 2 follow là 100% nhưng không nên được đánh giá cao hơn tăng từ 1000 lên 1100), hệ thống sử dụng ngưỡng tối thiểu (`MinFollowersForRate` = 100).
  - Nếu `Previous < 100`: Sử dụng tăng trưởng tuyệt đối được chuẩn hóa.
  - Nếu `Previous >= 100`: Sử dụng % tăng trưởng thông thường.
  - **Capping**: Giới hạn mức tăng trưởng tối đa (1000%) để tránh gian lận (gaming).

#### B. Blog Post Velocity (Trọng số 0.4)
Đo lường sự đóng góp nội dung.
- **Logic**: Tổng số bài viết trong 30 ngày / 30.
- **Ý nghĩa**: Khuyến khích người dùng duy trì thói quen viết bài đều đặn.

## 4. Data Flow & Processing

Để tối ưu hiệu năng và tránh tính toán nặng nề mỗi khi có request, hệ thống sử dụng cơ chế **Batch Processing** (Xử lý theo lô) định kỳ.

### Quy trình hàng ngày (Daily Job)

1.  **Snapshot Collection**: Hệ thống liên tục hoặc định kỳ lưu `UserFollowerSnapshot` để làm mốc so sánh.
2.  **Calculation (Tính toán)**:
    - Job chạy (ví dụ: lúc 00:00).
    - Duyệt qua danh sách người dùng.
    - Tính `CompositeScore` dựa trên dữ liệu snapshot và số bài viết.
    - Cập nhật hoặc tạo mới bản ghi trong `UserVelocityScore`.
3.  **Ranking (Sắp xếp)**:
    - Sắp xếp toàn bộ `UserVelocityScore` theo điểm giảm dần.
    - Gán Rank Position (Hạng 1, 2, 3...) cho từng user.
4.  **Archiving (Lưu lịch sử)**:
    - Lưu snapshot trạng thái xếp hạng hiện tại vào `UserRankingHistory`.
    - Dữ liệu này dùng để hiển thị biểu đồ lịch sử hoặc tính toán sự thay đổi hạng (lên/xuống hạng so với hôm qua).

## 5. Data Model

### Entities chính

| Entity | Mô tả |
| :--- | :--- |
| `UserVelocityScore` | Lưu điểm số hiện tại, thứ hạng hiện tại và các chỉ số thành phần (Growth, Velocity) của User. Đây là bảng dùng để query cho API lấy Top Ranking. |
| `UserRankingHistory` | Lưu trữ lịch sử thứ hạng theo thời gian (Timeseries data). Dùng để vẽ biểu đồ và theo dõi xu hướng. |
| `UserFollowerSnapshot` | Lưu số lượng follower của user tại các thời điểm cụ thể. Dữ liệu thô dùng để tính toán Growth Rate. |

## 6. API Reference

### Public APIs
- **GET /ranking/top**
  - Lấy danh sách Top N người dùng có điểm cao nhất.
  - Hỗ trợ phân trang.
- **GET /ranking/me** (hoặc `/ranking/users/{id}`)
  - Lấy thông tin xếp hạng chi tiết của một user.
  - Bao gồm: Hạng hiện tại, điểm số, lịch sử thay đổi hạng trong 30 ngày.

### Internal/Admin APIs
- **POST /admin/ranking/recalculate**
  - Trigger thủ công việc tính toán lại toàn bộ bảng xếp hạng (hữu ích khi deploy logic mới hoặc fix bug dữ liệu).
