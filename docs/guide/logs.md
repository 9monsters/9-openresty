# 日志与终端

## Nginx 日志

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/logs/access` | access.log（最近 200 行） |
| GET | `/api/logs/error` | error.log（最近 200 行） |
| GET | `/api/logs/files` | 日志文件列表 |
| GET | `/api/logs/file?path=` | 指定日志文件内容 |

前端 LogViewer 支持自动刷新。

### 查询参数

| 参数 | 说明 |
|------|------|
| `lines` | 返回行数（默认 200） |
| `path` | 指定日志文件路径 |

## Web Terminal

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/terminal/ws` | WebSocket 双向 shell |

- 服务端启动 `/bin/bash`，cwd = nginx config 目录
- 前端使用 xterm.js
- 无认证裸终端 — 生产环境应限制 loopback 访问