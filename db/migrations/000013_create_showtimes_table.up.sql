CREATE TABLE "showtimes" (
  "id" bigserial PRIMARY KEY,
  "movie_cinema_id" bigint NOT NULL,
  "show_time" timestamp NOT NULL,
  "price" int NOT NULL
);

COMMENT ON COLUMN "showtimes"."show_time" IS '2026-06-10 19:30:00';

COMMENT ON COLUMN "showtimes"."price" IS 'ticket price for this show';

ALTER TABLE "showtimes" ADD FOREIGN KEY ("movie_cinema_id") REFERENCES "movie_cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;
