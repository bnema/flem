// Path: apps\api\features\tmdb\requests.tsx
import { TMDB_API_KEY, TMDB_API_URL } from "../../config";
import { Movie } from "@flem/types";
import { saveMovie } from "../../db/mongo-handlers";


export const searchMoviesByTitle = async (title: string) => {
  const response = await fetch(
    `${TMDB_API_URL}/search/movie?api_key=${TMDB_API_KEY}&query=${encodeURIComponent(title)}`
  );

  const data = await response.json();
  console.log(data);
  return data.results;
};

export const getMovieDetails = async (movieId: number) => {
  const response = await fetch(
    `${TMDB_API_URL}/movie/${movieId}?api_key=${TMDB_API_KEY}`
  );

  const data = await response.json();
  
  await saveMovie(data);

  return data;
};