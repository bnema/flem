// Path : apps/api/src/db/mongo-handlers.ts
import { connectDB } from "../config/mongodb";
import { MovieModel } from "./mongo-models";
import {Movie} from '@flem/types';
import { translateToFrench } from "../features/openai/requests";

// Function to save movies in mongoDB
export const saveMovie = async (data: Movie, language: string ='en') => {
async function saveMovie() {
        await connectDB();

        // Before saving we check if the movie is already in the database
        const movieExists = await MovieModel.exists({ id: data.id });
        if (movieExists) {
            console.log(`Movie ${data.id} already exists in the database`);
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
        });
    
        await movie.save();
        console.log(` Movie ${data.id} saved in the database`);
    
        
    }
    await saveMovie();


}

// Function to get movies from MongoDB
export const getMovie = async (movieId: number) => {
  await connectDB();

  const movie = await MovieModel.findOne({ id: movieId });

  return movie;
}