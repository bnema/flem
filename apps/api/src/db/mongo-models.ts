import { Movie } from "@flem/types";
import { Schema, model } from 'mongoose';

// Create a new movie schema for mongoDB and define the model in the same line
export const MovieModel = model<Movie>('Movie', new Schema({
  id: Number,
  title: String,
  genres: Array,
  overview: String,
  release_date: String,
  spoken_languages: Array,
  vote_average: Number,
  poster_path: String,
}));
