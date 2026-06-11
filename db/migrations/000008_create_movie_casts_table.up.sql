CREATE TABLE "movie_casts" (
  "movie_id" bigint NOT NULL,
  "cast_id" bigint NOT NULL
);

CREATE UNIQUE INDEX ON "movie_casts" ("movie_id", "cast_id");

ALTER TABLE "movie_casts" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "movie_casts" ADD FOREIGN KEY ("cast_id") REFERENCES "casts" ("id") DEFERRABLE INITIALLY IMMEDIATE;
