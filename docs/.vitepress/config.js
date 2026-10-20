import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'OpenResty Admin',
  description: 'Web-based single-machine OpenResty/NGINX management platform',
  lang: 'zh-CN',
  lastUpdated: true,
  base: '/9-openresty/',

  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: '指南', link: '/guide/' },
      { text: '部署', link: '/deploy/' },
      { text: '参考', link: '/reference/' },
      { text: '架构', link: '/architecture/' },
    ],

    sidebar: {
      '/guide/': [
        { text: '快速开始', items: [
          { text: '概览', link: '/guide/' },
          { text: '安装', link: '/guide/install' },
          { text: '配置', link: '/guide/config' },
        ]},
        { text: '功能模块', items: [
          { text: '仪表盘', link: '/guide/dashboard' },
          { text: '站点管理', link: '/guide/sites' },
          { text: 'Stream / Upstream', link: '/guide/proxy' },
          { text: '证书管理', link: '/guide/certificates' },
          { text: 'DNS 管理', link: '/guide/dns' },
          { text: '访问控制', link: '/guide/access' },
          { text: 'Lua 脚本', link: '/guide/lua' },
          { text: 'WAF 管理', link: '/guide/waf' },
          { text: '日志', link: '/guide/logs' },
          { text: '备份与恢复', link: '/guide/backup' },
          { text: '系统设置', link: '/guide/settings' },
          { text: '用户与认证', link: '/guide/users' },
        ]},
      ],
      '/deploy/': [
        { text: '部署指南', items: [
          { text: '概览', link: '/deploy/' },
          { text: '数据库初始化', link: '/deploy/database' },
          { text: '二进制部署', link: '/deploy/binary' },
          { text: 'systemd 守护', link: '/deploy/systemd' },
          { text: '前端构建', link: '/deploy/frontend' },
        ]},
      ],
      '/reference/': [
        { text: '参考文档', items: [
          { text: 'API 接口', link: '/reference/api' },
          { text: 'CLI 命令', link: '/reference/cli' },
          { text: '配置文件', link: '/reference/config-file' },
          { text: '环境变量', link: '/reference/env' },
          { text: '数据库表结构', link: '/reference/schema' },
        ]},
      ],
      '/architecture/': [
        { text: '架构决策', items: [
          { text: '概览', link: '/architecture/' },
          { text: 'ADR-003 JWT 认证', link: '/architecture/adr-003-jwt-auth' },
          { text: 'ADR-004 PostgreSQL', link: '/architecture/adr-004-postgresql' },
          { text: 'ADR-007 可观测性', link: '/architecture/adr-007-observability' },
          { text: 'ADR-008 实施路线', link: '/architecture/adr-008-roadmap' },
          { text: 'ADR-009 数据库 Schema', link: '/architecture/adr-009-schema' },
        ]},
        { text: '工程结构', items: [
          { text: '项目布局', link: '/architecture/openresty-admin-architecture' },
          { text: '模块清单', link: '/architecture/modules' },
          { text: '实现细节', link: '/architecture/implementation-details' },
        ]},
      ],
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com' }
    ],

    search: {
      provider: 'local'
    },

    editLink: {
      pattern: 'https://github.com'
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright 2026 OpenResty Admin'
    }
  }
})