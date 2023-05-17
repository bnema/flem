import { FastifyInstance } from "fastify";
import { Movie } from "@flem/types";
import { getMoviesByGenreAndDate, getMovieDetails } from "../features/tmdb/requests";
import { getMoviesFromGT3Prompt } from "../features/openai/handlers";

export const registerAIRoutes = (fastify: FastifyInstance) => {
  fastify.post<{ Body: { ids: number[] } }>(
    "/v1/ai/movies/post/ids",
    async (request, reply) => {
        try {
            // 1 - The user submit a list of movie IDs
            const ids = request.body.ids;

            // 2 - We retrieve the movies from the database
            const movies: Movie[] = await Promise.all(ids.map(id => getMovieDetails(id, 'english')));

            // 3 - We create a small summary of each movie
            const summaries = movies.map(movie => {
                return `${movie.title} - ${new Date(movie.release_date).getFullYear()} by ${movie.director} (${movie.genres.map(genre => genre.name).join(', ')})`;
            });

            // 4 - We add all the summaries into the chatgpt's prompt and we ask to retrieve 20 more movies based on those added summaries
            const prompt = {
              'role': 'system',
              'content': `Based on these summaries, please suggest 20 more movies: ${summaries.join('\n')}`,
            };

            // Use your method to communicate with GPT-3 here and retrieve the 20 movies
            // You need to implement "getMoviesFromGPT3Prompt" method which uses GPT-3 to generate a list of movie names
            const suggestedMovies = await getMoviesFromGPT3Prompt(prompt);

            // 5 - We search for the movies in the database and we return them to the user
            const detailedMovies: Movie[] = await Promise.all(suggestedMovies.map(movieName => searchMoviesByTitle(movieName)));

            reply.send(detailedMovies);
        } catch (err) {
            console.error(err);
            reply.status(500).send({ error: "Something went wrong" });
        }
        }
  );
};
