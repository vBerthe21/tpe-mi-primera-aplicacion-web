-- Create "genre" table
CREATE TABLE "genre" (
  "id" serial NOT NULL,
  "name_genre" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "movie" table
CREATE TABLE "movie" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "author" character varying(255) NOT NULL,
  "duration" integer NOT NULL,
  "score" numeric(3,1) NULL,
  "review" character varying(500) NULL,
  "description" character varying(500) NULL,
  PRIMARY KEY ("id")
);
-- Create "movie_genre" table
CREATE TABLE "movie_genre" (
  "movie_id" integer NOT NULL,
  "genre_id" integer NOT NULL,
  CONSTRAINT "pk_movie_genre" PRIMARY KEY ("movie_id", "genre_id"),
  CONSTRAINT "fk_movie_genre_genre" FOREIGN KEY ("genre_id") REFERENCES "genre" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_movie_genre_movie" FOREIGN KEY ("movie_id") REFERENCES "movie" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
