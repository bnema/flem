// Path: apps\api\features\tmdb\routes.tsx
import { FastifyInstance } from "fastify";
import {
  searchMoviesByTitle,
  getMovieDetails,
  getMinMaxMovieID,
} from "./requests";
import { Movie } from "@flem/types";

export const registerTmdbRoutes = (fastify: FastifyInstance) => {
  fastify.post<{ Body: { titles: string[] } }>(
    "/v1/tmdb/query/title",
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
  "/v1/tmdb/query/id",
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


fastify.get('/v1/tmdb/random10', async (request, reply) => {
  try {
    // Call the 'getMinMaxMovieID' function to get the minimum and maximum movie IDs
    const { minID, maxID } = await getMinMaxMovieID();

    // Generate an array of 10 random movie IDs between the minimum and maximum values
    const randomIds = Array.from({ length: 10 }, () => Math.floor(Math.random() * (maxID - minID + 1) + minID));

    // Use Promise.all to fetch the details of all 10 movies
    const results = await Promise.all(
      randomIds.map(async (id) => {
        const details = await getMovieDetails(id);
        return details;
      }),
    );

    // Send the movie details as a JSON response
    reply.send(results);
  } catch (err) {
    // If there is an error, log it to the console and send a 500 error response
    console.error(err);
    reply.status(500).send({ error: 'Something went wrong' });
  }
});
};