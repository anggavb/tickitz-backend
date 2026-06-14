WITH schedule_templates AS (
  SELECT *
  FROM (
    VALUES
      (1, 1::bigint, 1::bigint, '2026-06-14'::date, '2026-06-20'::date, 1::bigint, 50000),
      (2, 2::bigint, 1::bigint, '2026-06-14'::date, '2026-06-20'::date, 2::bigint, 65000),
      (3, 3::bigint, 1::bigint, '2026-06-14'::date, '2026-06-20'::date, 3::bigint, 80000),
      (4, 4::bigint, 2::bigint, '2026-06-14'::date, '2026-06-20'::date, 4::bigint, 50000)
  ) AS existing_templates(template_order, movie_id, cinema_id, start_date, end_date, showtime_id, price)

  UNION ALL

  SELECT
    m.id::int AS template_order,
    m.id AS movie_id,
    (((m.id - 1) % 3) + 1)::bigint AS cinema_id,
    CASE
      WHEN m.id BETWEEN 51 AND 60 THEN m.release_date
      ELSE '2026-06-14'::date
    END AS start_date,
    CASE
      WHEN m.id BETWEEN 51 AND 60 THEN (m.release_date + INTERVAL '6 days')::date
      ELSE '2026-06-20'::date
    END AS end_date,
    (((m.id - 1) % 4) + 1)::bigint AS showtime_id,
    CASE (m.id - 1) % 3
      WHEN 0 THEN 50000
      WHEN 1 THEN 65000
      ELSE 80000
    END AS price
  FROM movies m
  WHERE m.id BETWEEN 5 AND 60
)
INSERT INTO movie_cinemas (movie_id, cinema_id, show_date, showtime_id, price)
SELECT
  schedule_templates.movie_id,
  schedule_templates.cinema_id,
  generated_date::date AS show_date,
  schedule_templates.showtime_id,
  schedule_templates.price
FROM schedule_templates
CROSS JOIN LATERAL generate_series(schedule_templates.start_date, schedule_templates.end_date, interval '1 day') AS generated_date
ORDER BY
  schedule_templates.template_order,
  generated_date::date
ON CONFLICT (movie_id, cinema_id, show_date, showtime_id) DO UPDATE SET
  price = EXCLUDED.price,
  updated_at = now();

SELECT setval(pg_get_serial_sequence('movie_cinemas', 'id'), COALESCE((SELECT MAX(id) FROM movie_cinemas), 1), true);
