import { FastifyInstance } from "fastify";
import { allowedOrigins } from "../config";
import { publicOriginCheck, privateOriginCheck } from "../middlewares";

export const registerRoutes = (fastify: FastifyInstance) => {
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
};
