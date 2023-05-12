import { checkBlacklist } from "../../config/filters";
import { Movie } from "@flem/types";

export const validateMovieData = async (data: Movie, movieId: number) => {
  // If the movie does not have an id, title, overview or poster_path, we do not save or return the movie
  if (
    !data.id ||
    !data.title ||
    !data.overview ||
    !data.poster_path ||
    !data.genres
  ) {
    console.log(`Movie ${movieId} does not have an id or title or overview`);
    return false;
  } else if (data.adult) {
    console.log(`Movie ${movieId} is an adult movie`);
    return false;
  } else {
    // Pass the data to the blacklist filter function checkBlacklist
    const blacklistWords = await checkBlacklist(data);

    // If the blacklist filter is not empty, then we do not save or return the movie
    if (blacklistWords.length > 0) {
      console.log(
        `Movie ${movieId} contains the following blacklisted words: ${blacklistWords.join(
          ", "
        )}`
      );
      return false;
    }
  }

  return true;
};
