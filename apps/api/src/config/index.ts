import dotenv from "dotenv";
// Load token from .env.local
dotenv.config();

// CORS config
export const allowedOrigins = ["http://localhost:3000"];


// TMDB config
export const TMDB_API_KEY = process.env.TMDB_API_KEY;
export const TMDB_API_URL = "https://api.themoviedb.org/3";