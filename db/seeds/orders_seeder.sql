INSERT INTO orders (
    id,
    payment_reference,
    user_id,
    showtime_id,
    movie_cinema_id,
    total_price,
    payment_method_id,
    status,
    expired_at,
    created_at
)
VALUES
(
    gen_random_uuid(),
    'VA-20260611001',
    1,
    1,
    1,
    100000,
    1,
    'paid',
    NOW() + INTERVAL '1 hour',
    NOW() - INTERVAL '5 day'
),
(
    gen_random_uuid(),
    'VA-20260611002',
    1,
    2,
    2,
    50000,
    1,
    'pending',
    NOW() + INTERVAL '1 hour',
    NOW() - INTERVAL '2 hour'
),
(
    gen_random_uuid(),
    'VA-20260611003',
    1,
    3,
    3,
    150000,
    1,
    'paid',
    NOW() + INTERVAL '1 hour',
    NOW() - INTERVAL '2 day'
);