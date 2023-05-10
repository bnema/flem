// Path: apps\api\features\tmdb\routes.tsx
import { FastifyInstance } from 'fastify';
import { searchMoviesByTitle, getMovieDetails} from './requests';
import { Movie } from '@flem/types';

export const registerTmdbRoutes = (fastify: FastifyInstance) => {
  fastify.post<{ Body: { titles: string[] } }>('/v1/tmdb/query/title', async (request, reply) => {
    const { titles } = request.body;
    const results = await Promise.all(
      titles.map(async (title) => {
        const movies = await searchMoviesByTitle(title);

        return Promise.all(
          movies.map(async (movie: Movie) => {
            const details = await getMovieDetails(movie.id);
            return details;
          }),
        );
      }),
    );

    reply.send(results.flat());
  });

   fastify.post<{ Body: { ids: number[] } }>('/v1/tmdb/query/id', async (request, reply) => {
    const { ids } = request.body;
    const results = await Promise.all(
      ids.map(async (id) => {
        const details = await getMovieDetails(id);
        return details;
      }),
    );
    

    reply.send(results);
    });
};