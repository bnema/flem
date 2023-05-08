import Fastify, { FastifyInstance } from "fastify";
import middie, { NextHandleFunction } from "@fastify/middie";
import cors from "cors";
import { allowedOrigins } from "../config";

export const createServer = (): FastifyInstance => {
  const fastify = Fastify({
    logger: true,
  });

  const setupMiddleware = async () => {
    await fastify.register(middie);

    fastify.use(
      cors({
        origin: (origin, callback) => {
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

  setupMiddleware();

  return fastify;
};