// Path: apps\api\features\tmdb\requests.tsx
import { TMDB_API_KEY, TMDB_API_URL } from "../../config";
import { saveMovie, getMovie } from "../../db/mongo-handlers";


export const searchMoviesByTitle = async (title: string) => {
  const response = await fetch(
    `${TMDB_API_URL}/search/movie?api_key=${TMDB_API_KEY}&query=${encodeURIComponent(title)}`
  );

  const data = await response.json();
  return data.results;
};

export const getMovieDetails = async (movieId: number) => {
  try {
    // Check if the movie is already in the database to avoid unnecessary API calls
    const movie = await getMovie(movieId);
    if (movie) {
      return movie;
    }
    // If the movie is not in the database, fetch it from TMDB and save it in the database
    const response = await fetch(
      `${TMDB_API_URL}/movie/${movieId}?api_key=${TMDB_API_KEY}`
    );

    const data = await response.json();

    await saveMovie(data);

    return data;
  } catch (error) {
    console.error(error);
    throw new Error(`Failed to get movie details for id ${movieId}`);
  }
};

export const getMinMaxMovieID = async () => {
  // Fetch the data for the first page of movies sorted by popularity in ascending order and disable adult content
  const minResponse = await fetch(
    `${TMDB_API_URL}/discover/movie?api_key=${TMDB_API_KEY}&sort_by=popularity.asc&page=1&include_adult=false`
  );
  // Fetch the data for the 500th page of movies sorted by popularity in ascending order and disable adult content
  const maxResponse = await fetch(
    `${TMDB_API_URL}/discover/movie?api_key=${TMDB_API_KEY}&sort_by=popularity.asc&page=500&include_adult=false`
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
