INSERT INTO loyalty_tiers (id, name, min_point) VALUES
  (1, 'Nominee', 0),
  (2, 'Laurels', 500),
  (3, 'Silver Screen', 1500),
  (4, 'Golden Statuette', 3000),
  (5, 'Grand Prix / Palme d''Or', 5000)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  min_point = EXCLUDED.min_point;

SELECT setval(pg_get_serial_sequence('loyalty_tiers', 'id'), COALESCE((SELECT MAX(id) FROM loyalty_tiers), 1), true);

INSERT INTO categories (id, name) VALUES
  (1, 'Action'),
  (2, 'Adventure'),
  (3, 'Animation'),
  (4, 'Comedy'),
  (5, 'Drama'),
  (6, 'Horror'),
  (7, 'Romance'),
  (8, 'Sci-Fi'),
  (9, 'Thriller')
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name;

SELECT setval(pg_get_serial_sequence('categories', 'id'), COALESCE((SELECT MAX(id) FROM categories), 1), true);

INSERT INTO casts (id, name) VALUES
  (1, 'Aria Mahendra'),
  (2, 'Nadia Putri'),
  (3, 'Raka Pratama'),
  (4, 'Sinta Lestari'),
  (5, 'Dimas Wardana'),
  (6, 'Clara Wijaya'),
  (7, 'Bima Santoso'),
  (8, 'Maya Kirana'),
  (9, 'Reza Alamsyah'),
  (10, 'Tara Adinata'),
  (11, 'Kevin Hartono'),
  (12, 'Luna Permata')
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name;

SELECT setval(pg_get_serial_sequence('casts', 'id'), COALESCE((SELECT MAX(id) FROM casts), 1), true);

INSERT INTO locations (id, name) VALUES
  (1, 'Purwokerto'),
  (2, 'Jakarta'),
  (3, 'Bandung'),
  (4, 'Surabaya')
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name;

SELECT setval(pg_get_serial_sequence('locations', 'id'), COALESCE((SELECT MAX(id) FROM locations), 1), true);

INSERT INTO payment_methods (id, name, logo, updated_at) VALUES
  (1, 'Google Pay', '/payment/google-pay.png', now()),
  (2, 'VISA', '/payment/visa.png', now()),
  (3, 'gopay', '/payment/gopay.png', now()),
  (4, 'PayPal', '/payment/paypal.png', now()),
  (5, 'DANA', '/payment/dana.png', now()),
  (6, 'BCA', '/payment/bca.png', now()),
  (7, 'BRI', '/payment/bri.png', now()),
  (8, 'OVO', '/payment/ovo.png', now())
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  logo = EXCLUDED.logo,
  updated_at = EXCLUDED.updated_at;

SELECT setval(pg_get_serial_sequence('payment_methods', 'id'), COALESCE((SELECT MAX(id) FROM payment_methods), 1), true);
