import { FastifyInstance } from "fastify";
import { Movie, SummaryItemMovie } from "@flem/types";
import {
  searchMoviesByTitle,
  getMovieDetails,
} from "../features/tmdb/requests";
import { getMoviesFromGPT3 } from "../features/openai/handlers";

export const createSummaryFromMoviesAndSendItToGPT3 = async (movies: Movie[]) => {
  const summariesForGPT3: SummaryItemMovie[] = movies.map((movie) => {
    return {
      id: movie.id,
      title: movie.title,
      release_date: movie.release_date,
      genres: movie.genres,
    };
  });


  const suggestedMoviesResponseFromGPT3 = await getMoviesFromGPT3(summariesForGPT3);
  
  console.log(suggestedMoviesResponseFromGPT3);
  return suggestedMoviesResponseFromGPT3;
};


export const registerAIRoutes = (fastify: FastifyInstance) => {
  fastify.post<{ Body: { ids: number[] } }>(
    "/v1/ai/movies/post/ids",
    async (request, reply) => {
      try {
        const ids = request.body.ids;
        const movies: Movie[] = (
          await Promise.all(ids.map((id) => getMovieDetails(id, "english")))
        ).filter((movie) => movie !== undefined);

        const detailedMovies = await createSummaryFromMoviesAndSendItToGPT3(movies);
        reply.send(detailedMovies);

      } catch (err) {
        console.error(err);
        reply.status(500).send({ error: "Something went wrong" });
      }
    }
  );
};
