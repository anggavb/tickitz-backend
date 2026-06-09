ALTER TABLE "movie_cinemas"
ADD COLUMN "showtime_id" bigint NOT NULL,
ADD COLUMN "price" int NOT NULL,
ADD COLUMN "created_at" timestamp DEFAULT (now()),
ADD COLUMN "updated_at" timestamp;

ALTER TABLE "movie_cinemas" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id") DEFERRABLE INITIALLY IMMEDIATE;
COMMENT ON COLUMN "movie_cinemas"."price" IS 'ticket price for this show';