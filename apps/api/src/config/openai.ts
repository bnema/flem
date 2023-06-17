// Path : apps/api/src/config/openai.ts
import dotenv from "dotenv";

// Load token from .env.local
dotenv.config();

// OpenAI config
export const OPENAI_API_KEY = process.env.OPENAI_API_KEY;
export const OPENAI_URL = "https://api.openai.com/v1/chat/completions"
export const model  = 'gpt-3.5-turbo-0613';