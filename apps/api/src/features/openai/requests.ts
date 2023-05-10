import { OPENAI_API_KEY, OPENAI_URL, model } from "../../config/openai";
import { MovieModel } from "../../db/mongo-models";
import { Movie } from "@flem/types";

export async function translateToFrench(data: Movie) {
  const prompt = `Translate this JSON Movie object into French while respecting the original structure ${JSON.stringify(
    data
  )}`;

  const body = JSON.stringify({
    messages: prompt,
    model: model,
    stream: false,
  });

  try {
    const response = await fetch(OPENAI_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${OPENAI_API_KEY}`,
      },
      body: body,
    });
    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}
