INSERT INTO movie_cinemas (movie_id, cinema_id, show_date, showtime_id, price)
SELECT
  schedule_templates.movie_id,
  schedule_templates.cinema_id,
  generated_date::date AS show_date,
  schedule_templates.showtime_id,
  schedule_templates.price
FROM (
  VALUES
    (1, 1, '2026-06-14'::date, '2026-06-20'::date, 1, 50000),
    (2, 1, '2026-06-14'::date, '2026-06-20'::date, 2, 65000),
    (3, 1, '2026-06-14'::date, '2026-06-20'::date, 3, 80000),
    (4, 2, '2026-06-14'::date, '2026-06-20'::date, 4, 50000)
) AS schedule_templates(movie_id, cinema_id, start_date, end_date, showtime_id, price)
CROSS JOIN LATERAL generate_series(schedule_templates.start_date, schedule_templates.end_date, interval '1 day') AS generated_date
ON CONFLICT (movie_id, cinema_id, show_date, showtime_id) DO UPDATE SET
  price = EXCLUDED.price,
  updated_at = now();

SELECT setval(pg_get_serial_sequence('movie_cinemas', 'id'), COALESCE((SELECT MAX(id) FROM movie_cinemas), 1), true);
