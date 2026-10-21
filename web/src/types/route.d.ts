import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
    hiddenInSidebar?: boolean
    hideChildren?: boolean
    modules?: string[]
  }
}

export {}