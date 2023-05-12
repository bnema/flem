// Path: apps\api\features\tmdb\requests.tsx
import { TMDB_API_KEY, TMDB_API_URL } from "../../config";
import { saveMovie, getMovie, getMoviesByGenreAndDateFromDB } from "../../db/mongo-handlers";
import { checkBlacklist } from "../../config/filters";
import { Movie } from "@flem/types";
import { translateMovieToFrench } from "../openai/handlers";

export const searchMoviesByTitle = async (title: string) => {
  const response = await fetch(
    `${TMDB_API_URL}/search/movie?api_key=${TMDB_API_KEY}&query=${encodeURIComponent(
      title
    )}`
  );

  const data = await response.json();
  return data.results;
};

export const getMovieDetails = async (movieId: number, language: string) => {
  try {
    // Check if the movie is already in the database to avoid unnecessary API calls
    const movie = await getMovie(movieId, language);
    if (movie) {
      return movie;
    }
    // If the movie is not in the database, fetch it from TMDB and save it in the database
    const response = await fetch(
      `${TMDB_API_URL}/movie/${movieId}?api_key=${TMDB_API_KEY}&include_adult=false`
    );

    const data = await response.json();

    // If the movie does not have an id, title, overview or poster_path, we do not save or return the movie
    if (
      !data.id ||
      !data.title ||
      !data.overview ||
      !data.poster_path ||
      !data.genres
    ) {
      console.log(`Movie ${movieId} does not have an id or title or overview`);
      return;
    } else if (data.adult) {
      console.log(`Movie ${movieId} is an adult movie`);
      return;
    } else {
      // Pass the data to the blacklist filter function checkBlacklist
      const blacklistWords = await checkBlacklist(data);

      // If the blacklist filter is not empty, then we do not save or return the movie
      if (blacklistWords.length > 0) {
        console.log(
          `Movie ${movieId} contains the following blacklisted words: ${blacklistWords.join(
            ", "
          )}`
        );
        return;
      }
    }

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
    throw new Error(
      `TMDB API request failed: ${minResponse.statusText}, ${maxResponse.statusText}`
    );
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

  // Create getMoviesByGenreAndDate(genre, minDate, maxDate, quantity);
export const getMoviesByGenreAndDate = async (
  genre: number,
  minDate: Date,
  maxDate: Date,
  quantity: number
) => {
  const minDateString = minDate.toISOString().split("T")[0];
  const maxDateString = maxDate.toISOString().split("T")[0];

  // Before we fetch we check if the movies are already in the database to avoid unnecessary API calls
  let moviesAlreadyExists = await getMoviesByGenreAndDateFromDB(genre, minDate, maxDate);

  // If the movies are already in the database, we return them

  if (moviesAlreadyExists.length >= quantity) {
    console.log(`Returning ${moviesAlreadyExists.length} movies from the database`);
    return moviesAlreadyExists.slice(0, quantity);
  }

  const response = await fetch(
    `${TMDB_API_URL}/discover/movie?api_key=${TMDB_API_KEY}&with_genres=${genre}&primary_release_date.gte=${minDateString}&primary_release_date.lte=${maxDateString}&include_adult=false`
  );

  const data = await response.json();
  if (!data.results) {
    throw new Error("Unexpected response from TMDB API");
  }
  // Data results is an array of Movie objects
  // We validate the data with the type Movie
  const movies: Movie[] = data.results;

  // for each movie, we translate in French in the background
  movies.forEach(async (movie) => {
      translateMovieToFrench(movie)
        .then(frenchMovie => {

        })
        .catch(err => {
          console.error(`Error translating movie: ${err}`);
        });
    });
  // We return a slice of the array determined by the quantity parameter
  return movies.slice(0, quantity);


};
