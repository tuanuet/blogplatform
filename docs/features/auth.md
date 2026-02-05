# Authentication & User Profile System

## 1. Overview

Hệ thống **Authentication** (Xác thực) và **User Profile** (Hồ sơ người dùng) chịu trách nhiệm quản lý danh tính người dùng, bảo mật truy cập và lưu trữ thông tin cá nhân.

## 2. Architecture (Kiến trúc)

Hệ thống tuân theo kiến trúc Clean Architecture:

- **UseCase Layer**:
  - `AuthUseCase` (`internal/application/usecase/auth`): Xử lý đăng ký, đăng nhập (email/pass, social), xác thực email, logout.
  - `ProfileUseCase` (`internal/application/usecase/profile`): Quản lý thông tin hồ sơ, upload avatar.
- **Domain Service Layer**:
  - `UserService`: Quản lý logic CRUD user.
  - `SocialAuthService`: Tích hợp với các provider OAuth2 (Google, Facebook, Github).
- **Repository Layer**:
  - `UserRepository`: Lưu trữ thông tin User.
  - `SessionRepository`: Quản lý phiên đăng nhập (Token-based hoặc Redis).
  - `SocialAccountRepository`: Lưu liên kết tài khoản mạng xã hội.

## 3. Core Features

### A. Authentication
- **Register**: Đăng ký tài khoản mới bằng Email/Password. Mật khẩu được mã hóa bằng `bcrypt`.
- **Login**: Đăng nhập bằng Email/Password, trả về `SessionID`.
- **Social Login**: Hỗ trợ đăng nhập qua Google, GitHub. Tự động tạo user nếu chưa tồn tại.
- **Email Verification**: Gửi email xác thực khi đăng ký.
- **Logout**: Hủy phiên đăng nhập.

### B. User Profile
- **View Profile**: Xem thông tin cá nhân (Private) hoặc hồ sơ công khai của người khác (Public).
- **Update Profile**: Cập nhật thông tin: Display Name, Bio, Website, Location, Social Handles, Birthday, Gender.
- **Avatar Upload**: Upload ảnh đại diện, validate định dạng (jpg, png, gif, webp) và kích thước (<5MB).

## 4. Security & Data Model

### User Entity
| Field | Mô tả |
| :--- | :--- |
| `Email` | Định danh duy nhất. |
| `PasswordHash` | Mật khẩu mã hóa (Bcrypt). |
| `IsActive` | Trạng thái kích hoạt. |
| `EmailVerifiedAt`| Thời điểm xác thực email. |
| `SocialAccounts` | Danh sách tài khoản MXH liên kết. |

### Session Management
- Sử dụng cơ chế lưu trữ Session ID (trong Redis hoặc DB).
- TTL mặc định: 24h.

## 5. API Reference

### Auth APIs
- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/login/{provider}` (Social Login redirect)
- `GET /auth/callback/{provider}` (Social Login callback)
- `POST /auth/logout`
- `GET /auth/verify-email`

### Profile APIs
- `GET /users/me` (Get own profile)
- `GET /users/{id}` (Get public profile)
- `PUT /users/me` (Update profile)
- `POST /users/me/avatar` (Upload avatar)
