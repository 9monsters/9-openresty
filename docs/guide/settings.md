# 系统设置

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/settings` | 获取所有设置 |
| PUT | `/api/settings` | 更新设置 |

## 设置项

```json
{
  "paths": {
    "nginxBin": "/usr/local/opt/openresty/bin/openresty",
    "nginxConfDir": "/usr/local/etc/openresty",
    "nginxConfPath": "/usr/local/etc/openresty/nginx.conf",
    "nginxLogDir": "/usr/local/var/log/nginx",
    "nginxPIDPath": "/usr/local/var/run/openresty.pid"
  },
  "theme": "auto",
  "language": "zh-CN",
  "autoBackup": true,
  "backupIntervalHours": 24,
  "backupRetentionCount": 30
}
```

## 优先级

```
CLI 参数 > 配置文件 > 环境变量 > 默认值
```

## 路径热更新

修改 paths 后即时生效，无需重启服务。系统会验证路径有效性。