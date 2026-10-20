# 用户与认证

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/users/login` | 用户登录（返回 JWT） |
| GET | `/api/users` | 用户列表 |
| POST | `/api/users` | 创建用户 |
| PUT | `/api/users/:id` | 更新用户 |
| PUT | `/api/users/:id/password` | 修改密码 |
| GET | `/api/users/profile` | 当前用户信息 |

## JWT 认证

- 算法：HMAC-SHA256
- 有效期：24 小时
- 签名密钥：`JWT_SIGNING_KEY` 环境变量
- 中间件自动拦截未授权请求（除 `/api/users/login`）

## 用户结构

```json
{
  "id": 1,
  "username": "admin",
  "email": "admin@example.com",
  "role": "admin",
  "createdAt": "2026-01-01T00:00:00Z"
}
```

## 角色

| 角色 | 权限 |
|------|------|
| `admin` | 全部权限 |
| `editor` | 配置读写，无用户管理 |
| `viewer` | 只读 |