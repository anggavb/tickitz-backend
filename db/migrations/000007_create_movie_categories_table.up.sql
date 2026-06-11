CREATE TABLE "movie_categories" (
  "movie_id" bigint NOT NULL,
  "category_id" bigint NOT NULL
);

CREATE UNIQUE INDEX ON "movie_categories" ("movie_id", "category_id");

ALTER TABLE "movie_categories" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "movie_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") DEFERRABLE INITIALLY IMMEDIATE;
