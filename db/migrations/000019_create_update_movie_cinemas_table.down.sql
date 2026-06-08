ALTER TABLE "movie_cinemas"
DROP COLUMN IF EXISTS "showtime_id",
DROP COLUMN IF EXISTS "price",
DROP COLUMN IF EXISTS "created_at",
DROP COLUMN IF EXISTS "updated_at";