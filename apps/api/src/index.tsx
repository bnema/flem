// ESM
import Fastify from 'fastify'

const fastify = Fastify({
  logger: false
})

fastify.get('/', async (request, reply) => {
  return { hello: 'world' }
})

/**
 * Run the server!
 */
const start = async () => {
  try {
    await fastify.listen({ port: 3333 })
  } catch (err) {
    console.error(err)
  }
}
start()