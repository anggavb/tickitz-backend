DROP INDEX IF EXISTS order_details_movie_cinema_seat_unique;
DROP INDEX IF EXISTS order_details_movie_cinema_id_index;
DROP INDEX IF EXISTS orders_movie_cinema_id_index;

ALTER TABLE "order_details" DROP CONSTRAINT IF EXISTS order_details_movie_cinema_id_fkey;
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS orders_movie_cinema_id_fkey;

ALTER TABLE "order_details" DROP COLUMN IF EXISTS "movie_cinema_id";
ALTER TABLE "orders" DROP COLUMN IF EXISTS "movie_cinema_id";

CREATE TABLE "movie_cinemas_old" (
  "id" bigserial PRIMARY KEY,
  "movie_id" bigint NOT NULL,
  "cinema_id" bigint NOT NULL,
  "start_date" timestamp NOT NULL,
  "end_date" timestamp NOT NULL,
  "showtime_id" bigint NOT NULL,
  "price" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);

ALTER TABLE "movie_cinemas_old" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movie_cinemas_old" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movie_cinemas_old" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id") DEFERRABLE INITIALLY IMMEDIATE;

INSERT INTO "movie_cinemas_old" (
  movie_id,
  cinema_id,
  start_date,
  end_date,
  showtime_id,
  price,
  created_at,
  updated_at
)
SELECT
  movie_id,
  cinema_id,
  min(show_date)::timestamp AS start_date,
  max(show_date)::timestamp AS end_date,
  showtime_id,
  price,
  min(created_at) AS created_at,
  max(updated_at) AS updated_at
FROM "movie_cinemas"
GROUP BY movie_id, cinema_id, showtime_id, price;

DROP TABLE "movie_cinemas";

ALTER TABLE "movie_cinemas_old" RENAME TO "movie_cinemas";
ALTER SEQUENCE "movie_cinemas_old_id_seq" RENAME TO "movie_cinemas_id_seq";
ALTER INDEX "movie_cinemas_old_pkey" RENAME TO "movie_cinemas_pkey";
