ALTER TABLE "showtimes"
DROP COLUMN IF EXISTS "movie_cinema_id",
DROP COLUMN IF EXISTS "show_time",
DROP COLUMN IF EXISTS "price",
ADD COLUMN "showtime" time NOT NULL;

COMMENT ON COLUMN "showtimes"."showtime" IS '08:30:00';