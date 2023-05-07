// Path: apps\api\src\index.tsx
// Import Fastify and necessary types from the 'fastify' package
import Fastify,{ FastifyInstance, RouteOptions } from 'fastify';
// Needed types for messages and responses
import { IncomingMessage, ServerResponse } from "http";
import { AddressInfo } from "net";

// Import middie (Fastify middleware support) and the cors package
import middie, { NextHandleFunction } from "@fastify/middie";
import cors from "cors";

// Create a Fastify instance with the logger 
const fastify = Fastify({
  logger: true,
});

// Allowed origins for incoming requests
const allowedOrigins = ["http://localhost:3000"];

// Function to set up middleware for the Fastify instance
const setupMiddleware = async () => {
  // Register the middie plugin to support middleware
  await fastify.register(middie);

  // Add CORS middleware with specific origins and credentials set to true
  fastify.use(
    cors({
      origin: (origin, callback) => {
        // Allow requests with no origin (like mobile apps or curl requests)
        if (!origin || allowedOrigins.includes(origin)) {
          callback(null, true);
        } else {
          callback(new Error("Not allowed by CORS"));
        }
      },
      credentials: true,
    })
  );
};

// Middleware for public routes
const publicOriginCheck: RouteOptions['preHandler'] = (request, reply, done) => {
  if (!request.headers.origin || allowedOrigins.includes(request.headers.origin)) {
    done();
  } else {
    reply.code(403).send({ error: 'Forbidden' });
  }
};

// Middleware for private routes (private means that the origin must be in the allowedOrigins array)
const privateOriginCheck: RouteOptions['preHandler'] = (request, reply, done) => {
  if (allowedOrigins.includes(request.headers.origin || '')) {
    done();
  } else {
    reply.code(403).send({ error: 'Forbidden' });
  }
};

// Define route for health checks
fastify.get('/v1/healthcheck', async (request, reply) => {
  if (!request.headers.origin || allowedOrigins.includes(request.headers.origin)) {
    // Public route, return simple response
    reply.send({
      status: 'ok',
    });
  } else if (allowedOrigins.includes(request.headers.origin || '')) {
    // Private route, return detailed server information
    fastify.server.getConnections((error, count) => {
      if (error) {
        reply.send({
          status: 'error',
          error: error.message,
        });
      } else {
        reply.send({
          status: 'ok',
          uptime: process.uptime(),
          connections: count,
        });
      }
    });
  } else {
    // Origin is not allowed, return 403 Forbidden
    reply.code(403).send({ error: 'Forbidden' });
  }
});

// Define a private route for testing purposes
fastify.get(
  '/v1/private',
  { preHandler: privateOriginCheck },
  async (request, reply) => {
    reply.send({
      status: 'ok',
      message: 'This is a private route',
    });
  }
);


// Function to start the Fastify server
const start = async () => {
  try {
    // Set up the middleware
    await setupMiddleware();

    // Start the Fastify server on port 3333 and avoid null
    await fastify.listen({ port: 3333 });

    const address = fastify.server.address();
    if (address !== null) {
      const port = typeof address === "string" ? address : address.port;
      console.log(`Server started on port ${port}`);
    } else {
      console.error("Failed to get the server address: address is null.");
    }
  } catch (err) {
    // Log any errors that occurred while starting the server
    console.error(err);
  }
};
// Call the start function to start the server
start();
