CREATE TABLE genre (
    id SERIAL PRIMARY KEY,
    name_genre VARCHAR(255) NOT NULL
);

CREATE TABLE movie (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    score DECIMAL(3,1),
    review VARCHAR(500),
    description VARCHAR(500)
);

CREATE TABLE movie_genre(
    movie_id INT NOT NULL,
    genre_id INT NOT NULL,

    CONSTRAINT pk_movie_genre PRIMARY KEY(movie_id, genre_id),

    CONSTRAINT fk_movie_genre_movie
    FOREIGN KEY (movie_id)
    REFERENCES movie(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_movie_genre_genre
    FOREIGN KEY (genre_id)
    REFERENCES genre(id)
    ON DELETE CASCADE
);