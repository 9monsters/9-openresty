# DNS 管理

## 功能概览

管理 DNS 提供商凭证、域名区域、DNS 分组和动态 DNS。

## DNS 凭证

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/dns-credentials` | 凭证列表 |
| POST | `/api/dns-credentials` | 添加凭证 |
| PUT | `/api/dns-credentials/:id` | 更新凭证 |
| DELETE | `/api/dns-credentials/:id` | 删除凭证 |

支持的提供商：Cloudflare、Aliyun DNS 等。

> **安全**: `secretKey` 字段在 API 响应中始终为 `-`（脱敏）。

## DNS 域名

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/dns-domains` | 域名列表 |
| POST | `/api/dns-domains` | 添加域名 |
| DELETE | `/api/dns-domains/:id` | 删除域名 |

## DNS 分组

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/dns-groups` | 分组列表 |
| POST | `/api/dns-groups` | 创建分组 |
| PUT | `/api/dns-groups/:id` | 更新分组 |

分组支持配置路由策略（轮询、加权、地理位置）。

## DDNS（动态 DNS）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/dns-ddns` | DDNS 列表 |
| POST | `/api/dns-ddns` | 创建 DDNS 条目 |

自动更新家庭和动态端点的 IP 地址。