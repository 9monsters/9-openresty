# 仪表盘

## 功能概览

仪表盘提供实时系统状态和 Nginx 性能指标的可视化展示。

## 服务器状态

- **CPU** — 型号、核心数、实时负载（1m / 5m / 15m）
- **内存** — 已用 / 总量，百分比进度条
- **Uptime** — 系统运行时间
- **主机信息** — hostname、平台、架构

**API**: `GET /api/system/info`

## Nginx 性能

- **连接统计** — Active / Accepts / Handled / Requests
- **连接状态** — Reading / Writing / Waiting
- **版本信息** — OpenResty 版本、运行状态

**API**: `GET /api/dashboard`

## 动态图表

两个仪表盘页面均支持 **echarts 实时折线图**，可切换刷新间隔：

| 按钮 | 间隔 |
|------|------|
| 1s | 每秒刷新（实时） |
| 5s | 每 5 秒刷新 |
| 10s | 每 10 秒刷新 |

底层通过 `internal/stats/collector.go` 每秒采样系统指标和 Nginx stub_status，存入 600 点环形缓冲区，前端定时拉取渲染。

## 自检

`GET /api/system/self-check` 返回 4 项检查结果：
- Nginx Binary — 是否在 PATH 中
- Config Syntax — `nginx -t` 结果
- Write Permission — config 目录是否可写
- Management API — 服务端口状态