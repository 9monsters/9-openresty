import { FastifyPluginAsync } from 'fastify'

export const accessListRoutes: FastifyPluginAsync = async (app) => {
  app.get('/access-lists', async () => ({
    lists: [
      { name: 'allow_internal', rules: [
        { directive: 'allow', address: '10.0.0.0/8' },
        { directive: 'allow', address: '172.16.0.0/12' },
        { directive: 'deny', address: 'all' },
      ]},
      { name: 'block_bots', rules: [
        { directive: 'deny', address: '66.249.0.0/16' },
        { directive: 'deny', address: '40.77.0.0/16' },
        { directive: 'allow', address: 'all' },
      ]},
    ]
  }))
}
