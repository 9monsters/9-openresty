# 前端构建

## 开发

```bash
pnpm --dir web dev
# Vite HMR on http://localhost:3001
# API proxy → http://localhost:3000
```

## 生产构建

```bash
pnpm --dir web build
```

默认产物输出到 `internal/embedweb/dist/`；设置 `STANDALONE=true` 时输出到 `web/dist/`，供独立部署。

## 嵌入 Go 二进制

Vite 默认将产物直接写入 `internal/embedweb/dist/`，`go build` 通过 `//go:embed dist` 嵌入，无需额外拷贝步骤。

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