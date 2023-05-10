import { connectDB } from "../config/mongodb";
import { MovieModel } from "./mongo-models";
import {Movie} from '@flem/types';

// Function to save movies in mongoDB
export const saveMovie = async (data: Movie) => {
async function saveMovie() {

        await connectDB();
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
        console.log('Movie saved');
    }
    
    await saveMovie();
}

// Function to get movies from mongoDB
export const getMovie = async (movieId: number) => {
    async function getMovie() {
        await connectDB();
    
        const movie = await MovieModel.findOne({ id: movieId });
    
        return movie;
    }

    return await getMovie();
}