// filters.ts 
import { Movie } from "@flem/types";

const blacklist = require("./blacklist.json");

type Blacklist = {
  [language: string]: string[];
};

// Convert the blacklist to a set of sets for faster lookup
const blacklistSets: { [language: string]: Set<string> } = {};
for (let language in blacklist) {
  blacklistSets[language] = new Set(blacklist[language]);
}

export const checkBlacklist = async (data: Movie) => {
  const title = data.title.toLowerCase();
  const overview = data.overview.toLowerCase();

  const titleWords = title.split(" ");
  const overviewWords = overview.split(" ");

  const words = [...titleWords, ...overviewWords];

  let blacklistWords = new Set<string>();

  // Loop through all keys in the blacklist object
  for (let language in blacklistSets) {
    // Get the list of blacklisted words for the current language
    const languageBlacklist = blacklistSets[language];

    // Filter the words that are in the language's blacklist
    const filteredWords = words.filter((word) => languageBlacklist.has(word));

    // Add the filtered words to the overall list of blacklisted words
    filteredWords.forEach((word) => blacklistWords.add(word));
  }

  return Array.from(blacklistWords);
};