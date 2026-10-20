# 前端构建

## 开发

```bash
pnpm --filter @openresty-admin/web dev
# Vite HMR on http://localhost:3001
# API proxy → http://localhost:3000
```

## 生产构建

```bash
pnpm --filter @openresty-admin/web build
```

产物在 `packages/web/dist/`，包含 `index.html` + `assets/` 目录（约 2.1MB）。

## 嵌入 Go 二进制

```bash
bash scripts/copy-web-dist.sh
go build -o bin/openresty-admin ./cmd/openresty-admin/
```

`scripts/copy-web-dist.sh` 将 `packages/web/dist/` 复制到 `internal/embedweb/dist/`，Go 编译时通过 `//go:embed dist` 嵌入二进制。

## 技术栈

| 层 | 技术 |
|----|------|
| 框架 | Vue 3 + Composition API |
| 路由 | Vue Router 4 (Hash History) |
| 状态 | Pinia |
| UI | antdv-next |
| CSS | UnoCSS + iOS 26 Liquid Glass tokens |
| 图表 | vue-echarts (ECharts 5) |
| 终端 | xterm.js |
| 构建 | Vite 5 |