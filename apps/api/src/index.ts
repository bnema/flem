import { createServer } from "./server";
import { registerRoutes } from "./routes";
import { AddressInfo } from "net";

const fastify = createServer();
registerRoutes(fastify);

const start = async () => {
  try {
    console.log("Starting server...");
    await fastify.listen({ port: 3333 });

    const address = fastify.server.address();
    if (address !== null) {
      const port = typeof address === "string" ? address : (address as AddressInfo).port;
      console.log(`Server started on port ${port}`);
    } else {
      console.error("Failed to get the server address: address is null.");
    }
  } catch (err) {
    console.error(err);
  }
};

start();
