CREATE TYPE "provider_name" AS ENUM (
  'google',
  'facebook'
);

CREATE TABLE "user_oauths" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "provider" provider_name NOT NULL,
  "access_token" varchar NOT NULL,
  "refresh_token" varchar,
  "expires_at" date
);

ALTER TABLE "user_oauths" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
