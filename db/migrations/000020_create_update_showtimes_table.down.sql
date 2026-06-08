ALTER TABLE "showtimes"
DROP COLUMN IF EXISTS "showtime",
ADD COLUMN "movie_cinema_id" INTEGER,
ADD COLUMN "show_time" TIMESTAMP,
ADD COLUMN "price" NUMERIC;
