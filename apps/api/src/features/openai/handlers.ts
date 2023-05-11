import { log } from "console";
import { OPENAI_API_KEY, OPENAI_URL, model } from "../../config/openai";
import { Movie } from "@flem/types";
import { saveMovieInFrench } from "../../db/mongo-handlers";

export async function translateMovieToFrench(data: Movie) {
  const prompt = {
    'role': 'system',
    'content': 'You are a helpful assistant translator.',
  };
  const userContent = {
    'role': 'user',
    'content': `Translate this JSON Movie object into French while respecting the original structure ${JSON.stringify(data)}. Please respond with a JSON-formatted structure but do not translate the keys.`,
  };

  const body = JSON.stringify({
    messages: [prompt, userContent],
    model: model,
  });

  let attempts = 0;
  const maxAttempts = 3;

  while (attempts < maxAttempts) {
    try {
      const response = await fetch(OPENAI_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${OPENAI_API_KEY}`,
        },
        body: body,
      });
      const result = await response.json();
      const messageContent = result.choices[0].message.content;
      // From this message, we need to extract the json object who starts at the first '{' and ends at the last '}'
      const startIndex = messageContent.indexOf('{');
      const endIndex = messageContent.lastIndexOf('}');
      const jsonContent = messageContent.substring(startIndex, endIndex + 1);
      const translatedMovie = JSON.parse(jsonContent);
      saveMovieInFrench(translatedMovie);
      return translatedMovie;
    } catch (error) {
      attempts++;
      if (attempts === maxAttempts) {
        throw new Error(`Failed to translate movie to French after ${maxAttempts} attempts: ${error}`);
      }
    }
  }
}
