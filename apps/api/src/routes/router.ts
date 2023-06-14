import { FastifyInstance } from "fastify";
import fastifySwagger from "@fastify/swagger";
import { allowedOrigins } from "../config";
import { publicOriginCheck, privateOriginCheck } from "../middlewares";
import { registerTmdbRoutes } from "../features/tmdb/routes";
import { registerAIRoutes } from "./ai";


export const registerRoutes = async (fastify: FastifyInstance) => {
  console.log("Registering routes...");

  await fastify.register(fastifySwagger, {
    routePrefix: "/v1/documentation",
    swagger: {
      info: {
        title: "AI Movies API",
        description: "AI movie suggestions API",
        version: "1.0.0",
      }
      // ...other swagger options...
    },
    exposeRoute: true,
  });

  console.log("Registered swagger");

  // Define route for health checks
  fastify.get("/v1/health", async (request, reply) => {
    if (
      !request.headers.origin ||
      allowedOrigins.includes(request.headers.origin)
    ) {
      // Public route, return simple response
      reply.send({
        status: "ok",
      });
    } else if (allowedOrigins.includes(request.headers.origin || "")) {
      // Private route, return detailed server information
      fastify.server.getConnections((error, count) => {
        if (error) {
          reply.send({
            status: "error",
            error: error.message,
          });
        } else {
          reply.send({
            status: "ok",
            uptime: process.uptime(),
            connections: count,
          });
        }
      });
    } else {
      // Origin is not allowed, return 403 Forbidden
      reply.code(403).send({ error: "Forbidden" });
    }
  });

  console.log("Registered health check");

  fastify.get(
    "/v1/private",
    { preHandler: privateOriginCheck },
    async (request, reply) => {
      reply.send({
        status: "ok",
        message: "This is a private route",
      });
    }
  );

  console.log("Registered private");

  await registerTmdbRoutes(fastify);

  console.log("Registered TMDB routes");

  registerAIRoutes(fastify);

  console.log("Registered AI routes");
};
