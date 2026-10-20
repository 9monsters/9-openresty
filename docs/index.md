---
layout: home
hero:
  name: "OpenResty Admin"
  text: "NGINX 集群管理平台"
  tagline: Web-based multi-machine OpenResty/NGINX management platform with Go + PostgreSQL + Vue 3
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/
    - theme: alt
      text: 部署指南
      link: /deploy/
features:
  - title: 多机集群管理
    details: Hub-Agent 架构，gRPC 通信，统一管理数十台 NGINX 节点。增量配置下发，失败自动回滚。
  - title: 完整的 Web UI
    details: Vue 3 + iOS 26 设计语言，20+ 管理页面，实时仪表盘、站点管理、证书、DNS、Lua 编辑。
  - title: PostgreSQL 持久化
    details: 所有模块数据均落库，30+ 张类型化表，GORM AutoMigrate，零 JSON 大字段。
  - title: JWT 认证
    details: HMAC-SHA256 签名，24h 有效期，中间件自动拦截未授权请求。
  - title: Prometheus 可观测
    details: /metrics 端点暴露 HTTP/gRPC/Deploy 指标，slog JSON 结构化日志 + trace_id 分布式追踪。
  - title: 单二进制部署
    details: Go 编译 + Vue 嵌入，一个 24MB 文件即可运行完整管理平台。systemd 守护开箱即用。