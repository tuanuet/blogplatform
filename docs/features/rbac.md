# RBAC System (Role-Based Access Control)

## 1. Overview

Hệ thống **RBAC** quản lý quyền truy cập của người dùng dựa trên vai trò (Role) và quyền hạn (Permission). Nó cho phép phân quyền chi tiết đến từng tài nguyên (Resource).

## 2. Architecture

- **UseCase Layer**:
  - `RoleUseCase` (`internal/application/usecase/role`): Quản lý Roles, gán Role cho User.
  - `PermissionUseCase` (`internal/application/usecase/permission`): Kiểm tra quyền truy cập (CheckPermission).
- **Caching Layer**:
  - Sử dụng Redis để cache quyền hạn và roles nhằm tối ưu hiệu năng (TTL 5-10 phút).
- **Domain Service Layer**:
  - `RoleService`, `PermissionService`: Logic nghiệp vụ cốt lõi.

## 3. Core Concepts

### Roles
- Là tập hợp các quyền hạn.
- Một User có thể có nhiều Role.
- Các Role tiêu chuẩn: `Admin`, `Moderator`, `User`.

### Permissions
- Được định nghĩa trên từng **Resource** (ví dụ: `blog`, `comment`, `user`).
- Quyền hạn dạng Bitmask:
  - `Read` (1)
  - `Create` (2)
  - `Update` (4)
  - `Delete` (8)

### Caching Strategy
- **User Roles Cache**: `rbac:user:{userID}:roles`
- **User Resource Permission Cache**: `rbac:user:{userID}:resource:{resource}`
- **Invalidation**: Tự động xóa cache khi Role hoặc Permission được cập nhật/xóa.

## 4. Features

### A. Role Management
- **CRUD Roles**: Tạo, sửa, xóa, liệt kê Roles.
- **Assign/Remove Role**: Gán hoặc gỡ Role cho User.
- **Set Permissions**: Cấu hình quyền hạn cho Role trên một Resource cụ thể.

### B. Access Control
- **Check Permission**: Kiểm tra xem User có quyền thực hiện hành động (Read/Write/...) trên Resource không.
- **Get User Permissions**: Lấy tổng hợp quyền hạn của User trên một Resource (gộp từ tất cả Roles của User).

## 5. API Reference

### Role APIs
- `POST /roles`
- `GET /roles`
- `GET /roles/{id}`
- `PUT /roles/{id}`
- `DELETE /roles/{id}`
- `POST /roles/{id}/permissions` (Set Permission)

### User Role APIs
- `GET /users/{id}/roles`
- `POST /users/{id}/roles` (Assign Role)
- `DELETE /users/{id}/roles/{roleID}` (Remove Role)
- `GET /users/{id}/permissions` (Check/Get Permissions)
