import { OPENAI_API_KEY, OPENAI_URL, model } from "../../config/openai";
import { Movie, GPTPrompt, GPTResponse, SummaryItemMovie, TranslationResponse } from "@flem/types";
import { saveMovieInFrench, getMovie } from "../../db/mongo-handlers";

async function fetchTranslation(
  prompt: GPTPrompt,
  userContent: GPTPrompt
): Promise<TranslationResponse> {
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
  const startIndex = messageContent.indexOf("{");
  const endIndex = messageContent.lastIndexOf("}");
  const jsonContent = messageContent.substring(startIndex, endIndex + 1);
  return JSON.parse(jsonContent);
}

export async function translateMovieToFrench(data: Movie) {
  const existingMovie = await getMovie(data.id, "french");
  if (existingMovie) {
    console.log(`Movie ${data.id} in french already exists in the database`);
    return existingMovie;
  }

  const prompt: GPTPrompt = {
    role: "system",
    content: "You are a helpful assistant translator.",
  };

  const userContent: GPTPrompt = {
    role: "user",
    content: `Translate this JSON Movie object into French while respecting the original structure ${JSON.stringify(
      data
    )}. Please respond with a JSON-formatted structure but do not translate the keys.`,
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
        throw new Error(
          `Failed to translate movie to French after ${maxAttempts} attempts: ${error}`
        );
      }
    }
  }
}

export async function getMoviesFromGPT3(
  summariesForGPT3: SummaryItemMovie[]
): Promise<Movie[]> {
  const systemContent: GPTPrompt = {
    role: "system",
    content: "You are a helpful assistant that suggests movies.",
  };

  const userContent: GPTPrompt = {
    role: "user",
    content: `Here are some movies I like: ${JSON.stringify(
      summariesForGPT3
    )}. Based on my movie preferences, please suggest 5 more movies that I might like. The response must be formatted as a single JSON array of movie objects, each having the following properties: "id" (only if you know the themoviedb.org one, otherwise leave it out), "title", "release_date", and "genres". The "genres" property is an array itself containing the genre strings. Please refer to the example format below:

    [
      {
        "id": (only if you know themoviedb.org one, otherwise leave it out)),
        "title": "The Matrix",
        "release_date": "1999-03-30",
        "genres": [
          "Action",
          "Science Fiction"
        ]
      },
      {
        "id": (only if you know the themoviedb.org one, otherwise leave it out)),
        "title": "Lord of the Rings",
        "release_date": "2001-12-19",
        "genres": [
          "Adventure",
          "Fantasy"
        ]
      }
    ]

Please note that all movie recommendations must be contained within the same array, not as separate entities and you should not include any other information in the response otherwise you will break the format.
`,
  };

  // Fetch translation from OpenAI
  const response = await fetchTranslation(systemContent, userContent);
  // Extract the message content from the response
  const messageContent = response.choices[0].message.content;

  console.log(messageContent);

  // Extract the JSON content from the message
  const startIndex = messageContent.indexOf("[");
  const endIndex = messageContent.lastIndexOf("]");
  const jsonContent = messageContent.substring(startIndex, endIndex + 1);
  // Parse the JSON content into an array of movies
  const movies: Movie[] = JSON.parse(jsonContent);

  // Return the movies
  return movies;
}
