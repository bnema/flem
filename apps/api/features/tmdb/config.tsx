// Path: apps\api\features\tmdb\config.tsx
import dotenv from "dotenv";
// Load token from .env.local
dotenv.config();

export const TMDB_API_KEY = process.env.TMDB_API_KEY;
export const TMDB_API_URL = "https://api.themoviedb.org/3";