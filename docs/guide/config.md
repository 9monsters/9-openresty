# 配置

## 优先级

`CLI 参数 > 配置文件 > 环境变量 > 默认值`

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DATABASE_URL` | `postgres://openresty:...@192.168.5.100:5432/openresty` | PostgreSQL 连接串 |
| `HOST` | `127.0.0.1` | 监听地址 |
| `PORT` | `3000` | HTTP 端口 |
| `JWT_SIGNING_KEY` | 自动生成 (ephemeral) | JWT 签名密钥 |

## JSON 配置文件

通过 `-c config.json` 指定：

```json
{
  "host": "0.0.0.0",
  "port": 3000,
  "paths": {
    "nginxBin": "/usr/local/bin/openresty",
    "nginxConfDir": "/usr/local/etc/openresty",
    "nginxConfPath": "/usr/local/etc/openresty/nginx.conf",
    "nginxLogDir": "/usr/local/var/log/nginx",
    "nginxPIDPath": "/usr/local/var/run/openresty.pid"
  }
}
```

所有字段均可选。未配置的字段回退到环境变量或默认值。

## CLI 参数

```bash
./bin/openresty-admin \
  --host 0.0.0.0 \
  --port 3000 \
  --nginx-bin /usr/local/bin/openresty \
  --nginx-conf-dir /etc/nginx/openresty \
  --nginx-conf-path /etc/nginx/openresty/nginx.conf \
  --nginx-log-dir /var/log/nginx \
  --nginx-pid-path /var/run/nginx.pid \
  --grpc-addr :9090
```

## 运行时配置

Nginx 路径可通过 Web UI 的 **系统设置** 页面实时修改，无需重启。修改后立即生效于 Host Adapter。