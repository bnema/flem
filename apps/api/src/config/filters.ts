// filters.ts 
import { Movie } from "@flem/types";

export const checkBlacklist = async (data: Movie) => {
    const blacklist = require("./blacklist.json");
    
    const title = data.title.toLowerCase();
    const overview = data.overview.toLowerCase();
    
    const titleWords = title.split(" ");
    const overviewWords = overview.split(" ");
    
    const words = [...titleWords, ...overviewWords];
    
    let blacklistWords: string[] = [];

    // Loop through all keys in the blacklist object
    for (let language in blacklist) {
        // Get the list of blacklisted words for the current language
        const languageBlacklist = blacklist[language];

        // Filter the words that are in the language's blacklist
        const filteredWords = words.filter((word) => languageBlacklist.includes(word));

        // Add the filtered words to the overall list of blacklisted words
        blacklistWords = [...blacklistWords, ...filteredWords];
    }

    return blacklistWords;
};
