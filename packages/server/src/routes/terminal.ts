import { FastifyPluginAsync } from 'fastify'
import { spawn } from 'node:child_process'
import { getNginxPaths } from '../services/nginx.js'

export const terminalRoutes: FastifyPluginAsync = async (app) => {
  // WebSocket terminal endpoint
  app.get('/terminal/ws', { websocket: true }, (socket, req) => {
    const paths = getNginxPaths()
    const shell = process.platform === 'win32' ? 'cmd.exe' : '/bin/bash'
    
    const child = spawn(shell, [], {
      cwd: paths.nginxConfDir,
      env: { ...process.env, TERM: 'xterm-256color', PATH: process.env.PATH },
    })

    child.stdout?.on('data', (data: Buffer) => {
      socket.send(data.toString('utf-8'))
    })

    child.stderr?.on('data', (data: Buffer) => {
      socket.send(data.toString('utf-8'))
    })

    child.on('close', () => socket.close())

    socket.on('message', (data: Buffer) => {
      child.stdin?.write(data)
    })

    socket.on('close', () => {
      child.kill()
    })

    socket.send('\r\nWelcome to OpenResty Admin Terminal\r\n\r\n')
  })
}
