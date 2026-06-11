INSERT INTO users (
  id,
  email,
  password,
  firstname,
  lastname,
  point,
  activation_token,
  token_expire_at,
  verified_at,
  role,
  loyalty_tier_id,
  updated_at
) VALUES
  (
    1,
    'admin@tickitz.test',
    '$argon2id$v=19$m=32768,t=2,p=1$V1IDAabTVgFYHQ/mMg8vgw$h2b28j57mCKX1kKQu8YSVSHWbASQSISq9489Nyw4jFc',
    'Tickitz',
    'Admin',
    0,
    NULL,
    NULL,
    now(),
    'admin',
    NULL,
    now()
  ),
  (
    2,
    'user.verified@tickitz.test',
    '$argon2id$v=19$m=32768,t=2,p=1$V1IDAabTVgFYHQ/mMg8vgw$h2b28j57mCKX1kKQu8YSVSHWbASQSISq9489Nyw4jFc',
    'Verified',
    'User',
    650,
    NULL,
    NULL,
    now(),
    'user',
    2,
    now()
  ),
  (
    3,
    'user.unverified@tickitz.test',
    '$argon2id$v=19$m=32768,t=2,p=1$V1IDAabTVgFYHQ/mMg8vgw$h2b28j57mCKX1kKQu8YSVSHWbASQSISq9489Nyw4jFc',
    'Unverified',
    'User',
    0,
    'seed-activation-token-user-unverified',
    now() + interval '24 hours',
    NULL,
    'user',
    1,
    now()
  )
ON CONFLICT (id) DO UPDATE SET
  email = EXCLUDED.email,
  password = EXCLUDED.password,
  firstname = EXCLUDED.firstname,
  lastname = EXCLUDED.lastname,
  point = EXCLUDED.point,
  activation_token = EXCLUDED.activation_token,
  token_expire_at = EXCLUDED.token_expire_at,
  verified_at = EXCLUDED.verified_at,
  role = EXCLUDED.role,
  loyalty_tier_id = EXCLUDED.loyalty_tier_id,
  updated_at = EXCLUDED.updated_at;

SELECT setval(pg_get_serial_sequence('users', 'id'), COALESCE((SELECT MAX(id) FROM users), 1), true);
