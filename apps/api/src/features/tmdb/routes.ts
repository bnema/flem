// Path: apps\api\features\tmdb\routes.tsx
import { FastifyInstance } from "fastify";
import {
  searchMoviesByTitle,
  getMovieDetails,
  getMinMaxMovieID,
} from "./requests";
import { Movie } from "@flem/types";

export const registerTmdbRoutes = (fastify: FastifyInstance) => {
  // Route to return all the movies with a given title
  fastify.post<{ Body: { titles: string[] } }>(
    "/v1/tmdb/movies/post/title",
    async (request, reply) => {
      try {
        const { titles } = request.body;
        const results = await Promise.all(
          titles.map(async (title) => {
            const movies = await searchMoviesByTitle(title);

            return Promise.all(
              movies.map(async (movie: Movie) => {
                const details = await getMovieDetails(movie.id);
                return details;
              })
            );
          })
        );

        reply.send(results.flat());
      } catch (err) {
        console.error(err);
        reply.status(500).send({ error: "Something went wrong" });
      }
    }
  );

  fastify.post<{ Body: { ids: number[] } }>(
    // Route to return all the movies from a given list of IDs
    "/v1/tmdb/movies/post/ids",
    async (request, reply) => {
      try {
        const { ids } = request.body;
        const results = await Promise.all(
          ids.map(async (id) => {
            const details = await getMovieDetails(id);
            return details;
          })
        );

        reply.send(results);
      } catch (err) {
        console.error(err);
        reply.status(500).send({ error: "Something went wrong" });
      }
    }
  );

fastify.get("/v1/tmdb/random10", async (request, reply) => {
  try {
    const { minID, maxID } = await getMinMaxMovieID();
    const results = [];

    while (results.length < 10) {
      const id = Math.floor(Math.random() * (maxID - minID + 1) + minID);
      const details = await getMovieDetails(id);

      // If 'details' is not undefined, push it to the results array
      if (details) {
        results.push(details);
      } else {
        console.log(`Invalid movie ${id}, refetching...`);
      }
    }

    reply.send(results);
  } catch (err) {
    console.error(err);
    reply.status(500).send({ error: "Something went wrong" });
  }
});
};
