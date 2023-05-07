
// Tous les types de données dont le front et le back ont besoin

export type User = {
    id: number;
    username: string;
    email: string;
    password: string;
    role: string;
    createdAt: string;
    updatedAt: string;
};

export type UserList = User[];

export type Movie = {
    id: number;
    title: string;
    description: string;
    releaseDate: string;
    duration: number;
    actors: string;
    director: string;
    poster: string;
    createdAt: string;
    updatedAt: string;
};

export type MovieList = Movie[];

export type TVShow = {
    id: number;
    title: string;
    description: string;
    releaseDate: string;
    actors: string;
    director: string;
    poster: string;
    created: string;
    updated: string;
};

export type TVShowList = TVShow[];

export type TVShowHasSeasons = {
    id: number;
    tvShowId: number;
    seasonNumber: number;
    createdAt: string;
    updatedAt: string;
};

export type UserHasMovies = {
    id: number;
    userId: number;
    movieId: number;
    createdAt: string;
    updatedAt: string;
};

export type UserHasTVShows = {
    id: number;
    userId: number;
    tvShowId: number;
    createdAt: string;
    updatedAt: string;
};

