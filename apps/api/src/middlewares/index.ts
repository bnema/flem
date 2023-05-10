import { RouteOptions } from "fastify";
import { allowedOrigins } from "../config";

export const publicOriginCheck: RouteOptions["preHandler"] = (request, reply, done) => {
  if (!request.headers.origin || allowedOrigins.includes(request.headers.origin)) {
    done();
  } else {
    reply.code(403).send({ error: "Forbidden" });
  }
};

export const privateOriginCheck: RouteOptions["preHandler"] = (request, reply, done) => {
  if (allowedOrigins.includes(request.headers.origin || "")) {
    done();
  } else {
    reply.code(403).send({ error: "Forbidden" });
  }
};
