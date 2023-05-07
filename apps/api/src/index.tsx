// ESM
import Fastify from 'fastify';
import {Movie, TVShow, User} from 'types';

const fastify = Fastify({
  logger: false
})

fastify.get('/', async (request, reply) => {
  // Empty
})

fastify.get('/v1/healthcheck', async (request, reply) => {
  fastify.server.getConnections((error, count) => {
    if (error) {
      reply.send({
        status: 'error',
        error: error.message
      })
    } else {
      reply.send({
        status: 'ok',
        uptime: process.uptime(),
        // Show how many connections are currently handled
        connections: count
      })
    }
  })
})

const start = async () => {
  try {
    await fastify.listen({ port: 3333 })
  } catch (err) {
    console.error(err)
  }
}

start()