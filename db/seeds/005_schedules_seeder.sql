INSERT INTO showtimes (id, showtime) VALUES
  (1, '10:00:00'),
  (2, '14:30:00'),
  (3, '19:30:00'),
  (4, '10:00:00'),
  (5, '14:30:00'),
  (6, '19:30:00'),
  (7, '10:00:00'),
  (8, '14:30:00'),
  (9, '19:30:00'),
  (10, '10:00:00'),
  (11, '14:30:00'),
  (12, '19:30:00'),
  (13, '10:00:00'),
  (14, '14:30:00'),
  (15, '19:30:00'),
  (16, '10:00:00'),
  (17, '14:30:00'),
  (18, '19:30:00'),
  (19, '10:00:00'),
  (20, '14:30:00'),
  (21, '19:30:00'),
  (22, '10:00:00'),
  (23, '14:30:00'),
  (24, '19:30:00')
ON CONFLICT (id) DO UPDATE SET
  showtime = EXCLUDED.showtime;

SELECT setval(pg_get_serial_sequence('showtimes', 'id'), COALESCE((SELECT MAX(id) FROM showtimes), 1), true);


INSERT INTO movie_cinemas (id, movie_id, cinema_id, start_date, end_date, showtime_id, price) VALUES
  (1, 1, 1, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 1, 10000),
  (2, 2, 1, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 2, 15000),
  (3, 3, 1, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 3, 20000),
  (4, 4, 2, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 4, 10000),
  (5, 5, 2, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 5, 15000),
  (6, 6, 3, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 6, 20000),
  (7, 1, 3, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 7, 10000),
  (8, 2, 3, '2026-06-10 00:00:00', '2026-06-16 23:59:59', 8, 15000)
ON CONFLICT (id) DO UPDATE SET
  movie_id = EXCLUDED.movie_id,
  cinema_id = EXCLUDED.cinema_id,
  start_date = EXCLUDED.start_date,
  end_date = EXCLUDED.end_date;

SELECT setval(pg_get_serial_sequence('movie_cinemas', 'id'), COALESCE((SELECT MAX(id) FROM movie_cinemas), 1), true);