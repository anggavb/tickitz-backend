CREATE TYPE "roleUser" AS ENUM (
  'user',
  'admin'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "firstname" varchar,
  "lastname" varchar,
  "phone_number" char(14),
  "photo" varchar,
  "point" int DEFAULT 0,
  "activation_token" varchar,
  "token_expire_at" timestamp,
  "verified_at" timestamp,
  "role" "roleUser" DEFAULT 'user',
  "loyalty_tier_id" bigint,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);

COMMENT ON COLUMN "users"."loyalty_tier_id" IS 'nullable if role admin';

ALTER TABLE "users" ADD FOREIGN KEY ("loyalty_tier_id") REFERENCES "loyalty_tiers" ("id") DEFERRABLE INITIALLY IMMEDIATE;
