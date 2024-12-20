import PocketBase from 'pocketbase';
// Tous les types de données dont le front et le back ont besoin


// TMDB Movie type
export type Movie = {
  language: string;
  adult: boolean;
  backdrop_path: string;
  belongs_to_collection: null | unknown;
  director: string;
  budget: number;
  genres: {
    id: number;
    name: string;
  }[];
  homepage: string;
  id: number;
  imdb_id: string;
  original_language: string;
  original_title: string;
  overview: string;
  popularity: number;
  poster_path: string;
  production_companies: {
    id: number;
    logo_path: string | null;
    name: string;
    origin_country: string;
  }[];
  production_countries: {
    iso_3166_1: string;
    name: string;
  }[];
  release_date: string;
  revenue: number;
  runtime: number;
  spoken_languages: {
    english_name: string;
    iso_639_1: string;
    name: string;
  }[];
  status: string;
  tagline: string;
  title: string;
  video: boolean;
  vote_average: number;
  vote_count: number;
};

export type GPTPrompt = {
  role: string;
  content: string;
};

export type GPTResponse = {
  choices: {
    message: {
      content: string
    }
  }[]
};

export type PocketBaseInstance = PocketBase;

export interface SummaryItemMovie {
  id: number;
  title: string;
  release_date: string;
  genres: { id: number; name: string }[];
}

export type TranslationResponse = {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: [
    {
      text: string;
      index: number;
      logprobs: null;
      finish_reason: string;
      json: any;
    }
  ];
};