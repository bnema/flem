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

export const getMinMaxMovieID = async () => {
  // Fetch the data for the first page of movies sorted by popularity in ascending order
  const minResponse = await fetch(
    `${TMDB_API_URL}/discover/movie?api_key=${TMDB_API_KEY}&sort_by=popularity.asc&page=1`
  );
  // Fetch the data for the 500th page of movies sorted by popularity in ascending order
  const maxResponse = await fetch(
    `${TMDB_API_URL}/discover/movie?api_key=${TMDB_API_KEY}&sort_by=popularity.asc&page=500`
  );

  // If either of the API requests failed, throw an error with the corresponding status texts
  if (!minResponse.ok || !maxResponse.ok) {
    throw new Error(`TMDB API request failed: ${minResponse.statusText}, ${maxResponse.statusText}`);
  }

  // Parse the response data into JSON format
  const minData = await minResponse.json();
  const maxData = await maxResponse.json();

  // If either of the responses did not contain the expected 'results' field, throw an error
  if (!minData.results || !maxData.results) {
    throw new Error("Unexpected response from TMDB API");
  }

  // Get the minimum and maximum movie IDs from the response data
  const minID = minData.results[0].id;
  const maxID = maxData.results[maxData.results.length - 1].id;

  // Return an object containing the minimum and maximum IDs
  return { minID, maxID };
};
