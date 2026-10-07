<script setup lang="ts">
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { getTerminalWsUrl } from '@/api/terminal'
import '@xterm/xterm/css/xterm.css'

const container = ref<HTMLDivElement>()
let term: Terminal | null = null
let ws: WebSocket | null = null

const cols = ref(80)
const rows = ref(24)

onMounted(() => {
  if (!container.value) return
  term = new Terminal({ cursorBlink: true, fontSize: 14, fontFamily: 'Menlo, Monaco, monospace', theme: { background: '#1e1e1e', foreground: '#d4d4d4' } })
  const fit = new FitAddon()
  term.loadAddon(fit)
  term.open(container.value)
  fit.fit()
  cols.value = term.cols
  rows.value = term.rows

  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const url = `${protocol}//${location.host}${getTerminalWsUrl()}`
  ws = new WebSocket(url)

  ws.onopen = () => term?.writeln('\x1b[32m● Connected to terminal\x1b[0m')
  ws.onmessage = (e) => term?.write(typeof e.data === 'string' ? e.data : new TextDecoder().decode(e.data))
  ws.onclose = () => term?.writeln('\n\x1b[31m● Disconnected\x1b[0m')
  ws.onerror = () => term?.writeln('\n\x1b[31m● Connection error\x1b[0m')

  term.onData((data) => ws?.send(data))
})

onUnmounted(() => { ws?.close(); term?.dispose() })
</script>

<template>
  <div>
    <h3 class="text-xl font-semibold mb-4">终端</h3>
    <div class="bg-gray-900 rounded-lg overflow-hidden">
      <div ref="container" class="p-2" :style="{ height: rows * 18 + 'px' }" />
    </div>
  </div>
</template>
