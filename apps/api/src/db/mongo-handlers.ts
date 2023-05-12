// Path : apps/api/src/db/mongo-handlers.ts
import { connectDB } from "../config/mongodb";
import { MovieModel } from "./mongo-models";
import {Movie} from '@flem/types';

// Function to save movies in mongoDB
export const saveMovie = async (data: Movie, language: string ='english') => {
  await connectDB();

  // Before saving we check if the movie is already in the database
  const movieExists = await MovieModel.exists({ id: data.id, language: language });
  if (movieExists) {
    console.log(`Movie ${data.id} in ${language} already exists in the database`);
    return;
  }
  
  const movie = new MovieModel({
    id: data.id,
    title: data.title,
    genres: data.genres,
    overview: data.overview,
    release_date: data.release_date,
    spoken_languages: data.spoken_languages,
    vote_average: data.vote_average,
    poster_path: data.poster_path,
    language: language,  // Add this line to set the language of the movie data
  });

  await movie.save();
  console.log(`Movie ${data.id} in ${language} saved in the database`);
}

// Function to translate and save movie in French
export const saveMovieInFrench = async (data: Movie) => {
  await connectDB();

  // Save the translated movie in MongoDB
  await saveMovie(data, 'french');
}

// Function to get movie from MongoDB
export const getMovie = async (movieId: number, language: string) => {
  await connectDB();

  const movie = await MovieModel.findOne({ id: movieId, language: language });


  return movie;
}

