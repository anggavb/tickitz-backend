INSERT INTO cinemas (id, location_id, name, logo) VALUES
  (1, 1, 'Tickitz Plaza Senayan', '/cinemas/tickitz-plaza-senayan.jpg'),
  (2, 2, 'Tickitz Paris Van Java', '/cinemas/tickitz-paris-van-java.jpg'),
  (3, 3, 'Tickitz Tunjungan Plaza', '/cinemas/tickitz-tunjungan-plaza.jpg')
ON CONFLICT (id) DO UPDATE SET
  location_id = EXCLUDED.location_id,
  name = EXCLUDED.name,
  logo = EXCLUDED.logo;

SELECT setval(pg_get_serial_sequence('cinemas', 'id'), COALESCE((SELECT MAX(id) FROM cinemas), 1), true);

WITH seat_rows(row_label, row_index) AS (
  VALUES ('A', 1), ('B', 2), ('C', 3), ('D', 4), ('E', 5), ('F', 6), ('G', 7)
),
seat_numbers(number) AS (
  SELECT generate_series(1, 14)
),
seed_seats AS (
  SELECT
    cinemas.id AS cinema_id,
    seat_rows.row_label AS "row",
    seat_numbers.number,
    CASE
      WHEN seat_rows.row_label = 'F' AND seat_numbers.number IN (10, 11) THEN 'love_nest'::"seatType"
      ELSE 'regular'::"seatType"
    END AS type
  FROM cinemas
  CROSS JOIN seat_rows
  CROSS JOIN seat_numbers
  WHERE cinemas.id IN (1, 2, 3)
)
INSERT INTO seats (cinema_id, "row", "number", "type")
SELECT cinema_id, "row", number, type
FROM seed_seats
ON CONFLICT (cinema_id, "row", "number") DO UPDATE SET
  "type" = EXCLUDED."type";

SELECT setval(pg_get_serial_sequence('seats', 'id'), COALESCE((SELECT MAX(id) FROM seats), 1), true);
