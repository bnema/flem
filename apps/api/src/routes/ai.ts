import { FastifyInstance } from "fastify";
import { Movie } from "@flem/types";
import { searchMoviesByTitle, getMovieDetails } from "../features/tmdb/requests";
import { getMoviesFromGPT3 } from "../features/openai/handlers";

export const registerAIRoutes = (fastify: FastifyInstance) => {
  fastify.post<{ Body: { ids: number[] } }>(
    "/v1/ai/movies/post/ids",
    async (request, reply) => {
        try {
            // 1 - The user submit a list of movie IDs
            const ids = request.body.ids; 

            // 2 - We retrieve the movies from the database
            const movies: Movie[] = (await Promise.all(ids.map(id => getMovieDetails(id, 'english')))).filter(movie => movie !== undefined);;
            // 3 - We create a small summary of each movie
            const summaries = movies.map(movie => {
                return `${movie.title} - ${new Date(movie.release_date).getFullYear()} (${movie.genres.map(genre => genre.name).join(', ')})`;
            });

            // 4 - We add all the summaries into the chatgpt's prompt and we ask to retrieve 20 more movies based on those added summaries
            const prompt = {
              'role': 'system',
              'content': `Based on these summaries, please suggest 20 more movies and returns it as a formatted JSON: ${summaries.join('\n')}`,
            };

            // 5 - We retrieve the movies from GPT-3 response
            const suggestedMovies = await getMoviesFromGPT3(prompt);
            const suggestedMoviesString = suggestedMovies.join("");


            // 6 Catch the beginning and the end "{...}" of the JSON inside the GPT-3 response and create a JSON object
            const suggestedMoviesCreateJSON = suggestedMoviesString
  .split("{")
  .filter((item) => item.includes("}"))
  .map((item) => `{${item}`);


  const suggestedMoviesJSON = JSON.parse(suggestedMoviesCreateJSON[0]);
            // 7 - We search by the title of each movie to retrieve more details
            const detailedMovies = await Promise.all(
              suggestedMoviesJSON.map(async (movie) => {
                const movies = await searchMoviesByTitle(movie);
                return movies[0];
              })
            );



            reply.send(detailedMovies);
        } catch (err) {
            console.error(err);
            reply.status(500).send({ error: "Something went wrong" });
        }
        }
  );
};
