import { OPENAI_API_KEY, OPENAI_URL, model } from "../../config/openai";
import { Movie, GPTPrompt, GPTResponse } from "@flem/types";
import { saveMovieInFrench, getMovie } from "../../db/mongo-handlers";


async function fetchTranslation(prompt: GPTPrompt, userContent: GPTPrompt): Promise<TranslationResponse> {
  const body = JSON.stringify({
    messages: [prompt, userContent],
    model: model,
  });

  const response = await fetch(OPENAI_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${OPENAI_API_KEY}`,
    },
    body: body,
  });

  return await response.json();
}

async function extractTranslatedMovie(response: GPTResponse): Promise<Movie> {
  const messageContent = response.choices[0].message.content;
  const startIndex = messageContent.indexOf('{');
  const endIndex = messageContent.lastIndexOf('}');
  const jsonContent = messageContent.substring(startIndex, endIndex + 1);
  return JSON.parse(jsonContent);
}

export async function translateMovieToFrench(data: Movie) {
  const existingMovie = await getMovie(data.id, 'french');
  if (existingMovie) {
    console.log(`Movie ${data.id} in french already exists in the database`);
    return existingMovie;
  }

  const prompt: GPTPrompt = {
    'role': 'system',
    'content': 'You are a helpful assistant translator.',
  };

  const userContent: GPTPrompt = {
    'role': 'user',
    'content': `Translate this JSON Movie object into French while respecting the original structure ${JSON.stringify(data)}. Please respond with a JSON-formatted structure but do not translate the keys.`,
  };

  let attempts = 0;
  const maxAttempts = 3;

  while (attempts < maxAttempts) {
    try {
      const translationResponse = await fetchTranslation(prompt, userContent);
      const translatedMovie = await extractTranslatedMovie(translationResponse);
      await saveMovieInFrench(translatedMovie);
      return translatedMovie;
    } catch (error) {
      attempts++;
      if (attempts === maxAttempts) {
        throw new Error(`Failed to translate movie to French after ${maxAttempts} attempts: ${error}`);
      }
    }
  }
}

export async function getMoviesFromGPT3(prompt: GPTPrompt) {
  // The content of the system message is your prompt and the user message is the question you want GPT-3 to answer.
  const systemContent: GPTPrompt = {
    'role': 'system',
    'content': 'You are a helpful assistant that suggests movies.',
  };

  const userContent: GPTPrompt = {
    'role': 'user',
    'content': `${prompt.content} Can you suggest 20 more movies that are similar?`,
  };

  let attempts = 0;
  const maxAttempts = 3;

  while (attempts < maxAttempts) {
    try {
      // Fetch translation from OpenAI
      const response = await fetchTranslation(systemContent, userContent);
      
      // Extract the message content from the response
      const messageContent = response.choices[0].message.content;

      // Split the message content into lines and remove empty lines
      const lines = messageContent.split('\n').filter(line => line.trim() !== '');

      // The movie suggestions are expected to start from the second line
      const movieSuggestions = lines.slice(1);

      // Return the movie suggestions
      return movieSuggestions;
    } catch (error) {
      attempts++;
      if (attempts === maxAttempts) {
        throw new Error(`Failed to get movie suggestions from GPT-3 after ${maxAttempts} attempts: ${error}`);
      }
    }
  }
}
