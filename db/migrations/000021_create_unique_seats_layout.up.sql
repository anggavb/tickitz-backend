CREATE UNIQUE INDEX IF NOT EXISTS seats_cinema_row_number_unique
ON "seats" ("cinema_id", "row", "number");

