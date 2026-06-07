CREATE TABLE "movies" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "release_date" date,
  "duration_in_minute" int,
  "director_name" varchar,
  "synopsis" text,
  "image" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);

CREATE INDEX ON "movies" ("name");

COMMENT ON COLUMN "movies"."duration_in_minute" IS '120 = 2 hours';
