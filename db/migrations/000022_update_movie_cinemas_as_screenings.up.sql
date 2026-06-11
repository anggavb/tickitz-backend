CREATE TABLE "movie_cinemas_new" (
  "id" bigserial PRIMARY KEY,
  "movie_id" bigint NOT NULL,
  "cinema_id" bigint NOT NULL,
  "show_date" date NOT NULL,
  "showtime_id" bigint NOT NULL,
  "price" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);

ALTER TABLE "movie_cinemas_new" ADD FOREIGN KEY ("movie_id") REFERENCES "movies" ("id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movie_cinemas_new" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "movie_cinemas_new" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id") DEFERRABLE INITIALLY IMMEDIATE;

INSERT INTO "movie_cinemas_new" (
  movie_id,
  cinema_id,
  show_date,
  showtime_id,
  price,
  created_at,
  updated_at
)
SELECT
  mc.movie_id,
  mc.cinema_id,
  generated_date::date AS show_date,
  mc.showtime_id,
  mc.price,
  COALESCE(mc.created_at, now()) AS created_at,
  mc.updated_at
FROM "movie_cinemas" mc
CROSS JOIN LATERAL generate_series(mc.start_date::date, mc.end_date::date, interval '1 day') AS generated_date
ON CONFLICT DO NOTHING;

DROP TABLE "movie_cinemas";

ALTER TABLE "movie_cinemas_new" RENAME TO "movie_cinemas";
ALTER SEQUENCE "movie_cinemas_new_id_seq" RENAME TO "movie_cinemas_id_seq";
ALTER INDEX "movie_cinemas_new_pkey" RENAME TO "movie_cinemas_pkey";

CREATE UNIQUE INDEX movie_cinemas_movie_cinema_date_time_unique
ON "movie_cinemas" ("movie_id", "cinema_id", "show_date", "showtime_id");

ALTER TABLE "orders"
ADD COLUMN IF NOT EXISTS "movie_cinema_id" bigint;

ALTER TABLE "order_details"
ADD COLUMN IF NOT EXISTS "movie_cinema_id" bigint;

ALTER TABLE "orders" ADD CONSTRAINT orders_movie_cinema_id_fkey
FOREIGN KEY ("movie_cinema_id") REFERENCES "movie_cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_details" ADD CONSTRAINT order_details_movie_cinema_id_fkey
FOREIGN KEY ("movie_cinema_id") REFERENCES "movie_cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;

CREATE INDEX IF NOT EXISTS orders_movie_cinema_id_index
ON "orders" ("movie_cinema_id");

CREATE INDEX IF NOT EXISTS order_details_movie_cinema_id_index
ON "order_details" ("movie_cinema_id");

CREATE UNIQUE INDEX IF NOT EXISTS order_details_movie_cinema_seat_unique
ON "order_details" ("movie_cinema_id", "seat_id")
WHERE "movie_cinema_id" IS NOT NULL;

