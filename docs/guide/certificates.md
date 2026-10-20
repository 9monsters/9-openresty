# 证书管理

## 功能概览

管理 SSL/TLS 证书，支持手动导入和 ACME 自动签发。

## 证书 CRUD

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/certificates` | 证书列表 |
| POST | `/api/certificates` | 导入证书 |
| PUT | `/api/certificates/:id` | 更新证书信息 |
| DELETE | `/api/certificates/:id` | 删除证书 |

证书行结构（列表页 PEM 字段自动脱敏）：

```json
{
  "id": 1,
  "name": "example.com",
  "domains": ["example.com", "*.example.com"],
  "issuer": "Let's Encrypt",
  "status": "valid",
  "autoRenew": true,
  "notBefore": "2026-01-01T00:00:00Z",
  "notAfter": "2026-04-01T00:00:00Z"
}
```

## ACME 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/acme-users` | ACME 账户列表 |
| POST | `/api/acme-users` | 创建 ACME 账户 |
| PUT | `/api/acme-users/:id` | 更新 ACME 账户 |

默认 provider: Let's Encrypt，server: `https://acme-v02.api.letsencrypt.org/directory`。