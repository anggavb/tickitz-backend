CREATE TABLE "movie_cinemas" (
  "id" bigserial PRIMARY KEY,
  "movie_id" bigint NOT NULL,
  "cinema_id" bigint NOT NULL,
  "start_date" timestamp NOT NULL,
  "end_date" timestamp NOT NULL
);

CREATE UNIQUE INDEX ON "movie_cinemas" ("movie_id", "cinema_id");

ALTER TABLE "movie_cinemas" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "movie_cinemas" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;
