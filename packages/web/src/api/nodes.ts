import api from './client'

export interface Node {
  id: number
  name: string
  host: string
  port: number
  status: string
  version: string
  uptime: string
  environment?: string
}

export function listNodes(): Promise<{ nodes: Node[] }> {
  return api.get('/nodes')
}

export function createNode(data: Partial<Node>): Promise<{ ok: boolean; node: Node }> {
  return api.post('/nodes', data)
}

export function updateNode(id: number, data: Partial<Node>): Promise<{ ok: boolean; node: Node }> {
  return api.put(`/nodes/${id}`, data)
}

export function deleteNode(id: number): Promise<{ ok: boolean }> {
  return api.delete(`/nodes/${id}`)
}

export function testNode(id: number): Promise<{ ok: boolean; reachable: boolean; latency: string }> {
  return api.post(`/nodes/${id}/test`)
}