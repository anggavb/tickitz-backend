CREATE UNIQUE INDEX IF NOT EXISTS order_details_showtime_id_seat_id_idx
ON "order_details" ("showtime_id", "seat_id");
