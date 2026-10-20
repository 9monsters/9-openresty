# 配置文件参考

通过 `-c config.json` 指定。所有字段可选。

```json
{
  "host": "0.0.0.0",
  "port": 3000,
  "paths": {
    "nginxBin": "/usr/local/opt/openresty/bin/openresty",
    "nginxConfDir": "/usr/local/etc/openresty",
    "nginxConfPath": "/usr/local/etc/openresty/nginx.conf",
    "nginxLogDir": "/usr/local/var/log/nginx",
    "nginxPIDPath": "/usr/local/var/run/openresty.pid"
  }
}
```

## 字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `host` | string | `127.0.0.1` | HTTP 监听地址 |
| `port` | int | `3000` | HTTP 监听端口 |
| `paths.nginxBin` | string | brew 路径 | Nginx 可执行文件 |
| `paths.nginxConfDir` | string | brew 路径 | Nginx 配置目录 |
| `paths.nginxConfPath` | string | brew 路径 | 主配置文件 |
| `paths.nginxLogDir` | string | brew 路径 | 日志目录 |
| `paths.nginxPIDPath` | string | brew 路径 | PID 文件路径 |

## 优先级

```
CLI 参数 > 配置文件 > 环境变量 > 默认值
```