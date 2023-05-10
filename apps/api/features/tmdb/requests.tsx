// Path: apps\api\features\tmdb\requests.tsx
import { TMDB_API_KEY, TMDB_API_URL } from "./config";


export const searchMoviesByTitle = async (title: string) => {
  const response = await fetch(
    `${TMDB_API_URL}/search/movie?api_key=${TMDB_API_KEY}&query=${encodeURIComponent(title)}`
  );

  const data = await response.json();
  console.log(data);
  return data.results; // Make sure to return the results array
};

export const getMovieDetails = async (movieId: number) => {
  const response = await fetch(
    `${TMDB_API_URL}/movie/${movieId}?api_key=${TMDB_API_KEY}`
  );

  const data = await response.json();
  // data is a json object with one or multiple movies details
    // We need to store the details in Redis but first we need to convert it to an array
    const movies = Array.isArray(data) ? data : [data];

  return data;
};