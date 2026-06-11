INSERT INTO movies (id, name, release_date, duration_in_minute, director_name, synopsis, image, updated_at, slug) VALUES
  (1, 'Echoes of Jakarta', '2026-05-12', 118, 'Rangga Wiratama', 'A sound engineer follows a mysterious recording through the late-night streets of Jakarta.', '/movies/echoes-of-jakarta.jpg', now(), 'echoes-of-jakarta'),
  (2, 'Orbit Seven', '2026-06-01', 132, 'Mira Anggraini', 'A small crew must repair a failing research station before its orbit decays.', '/movies/orbit-seven.jpg', now(), 'orbit-seven'),
  (3, 'The Last Screening', '2026-04-18', 101, 'Dewi Kartika', 'A closing cinema becomes the center of a haunting tied to its final reel.', '/movies/the-last-screening.jpg', now(), 'the-last-screening'),
  (4, 'Letters from Bandung', '2026-03-22', 109, 'Fajar Nugroho', 'Two strangers trade misplaced letters and discover an old family secret.', '/movies/letters-from-bandung.jpg', now(), 'letters-from-bandung'),
  (5, 'Laugh Track Heroes', '2026-06-07', 96, 'Sari Wulandari', 'A struggling comedy troupe gets one chaotic shot at saving their theater.', '/movies/laugh-track-heroes.jpg', now(), 'laugh-track-heroes'),
  (6, 'River of Ashes', '2026-02-14', 124, 'Bayu Prakoso', 'A detective returns to his hometown to investigate a case everyone wants forgotten.', '/movies/river-of-ashes.jpg', now(), 'river-of-ashes')
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  release_date = EXCLUDED.release_date,
  duration_in_minute = EXCLUDED.duration_in_minute,
  director_name = EXCLUDED.director_name,
  synopsis = EXCLUDED.synopsis,
  image = EXCLUDED.image,
  updated_at = EXCLUDED.updated_at,
  slug = EXCLUDED.slug;

SELECT setval(pg_get_serial_sequence('movies', 'id'), COALESCE((SELECT MAX(id) FROM movies), 1), true);

INSERT INTO movie_categories (movie_id, category_id) VALUES
  (1, 5), (1, 9),
  (2, 2), (2, 8), (2, 1),
  (3, 6), (3, 9),
  (4, 5), (4, 7),
  (5, 4), (5, 5),
  (6, 5), (6, 9), (6, 1)
ON CONFLICT (movie_id, category_id) DO NOTHING;

INSERT INTO movie_casts (movie_id, cast_id) VALUES
  (1, 1), (1, 2), (1, 9),
  (2, 3), (2, 6), (2, 11),
  (3, 4), (3, 8), (3, 12),
  (4, 2), (4, 5), (4, 10),
  (5, 7), (5, 10), (5, 11),
  (6, 1), (6, 6), (6, 9)
ON CONFLICT (movie_id, cast_id) DO NOTHING;
