CREATE TABLE "order_details" (
  "id" bigserial PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "seat_id" bigint NOT NULL,
  "showtime_id" bigint NOT NULL,
  "price" int NOT NULL
);

CREATE UNIQUE INDEX ON "order_details" ("order_id", "seat_id");

CREATE UNIQUE INDEX ON "order_details" ("showtime_id", "seat_id");

ALTER TABLE "order_details" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_details" ADD FOREIGN KEY ("seat_id") REFERENCES "seats" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_details" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id") DEFERRABLE INITIALLY IMMEDIATE;
