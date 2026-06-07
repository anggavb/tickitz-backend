CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE "statusOrder" AS ENUM (
  'pending',
  'waiting',
  'paid',
  'cancel'
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "payment_reference" varchar,
  "user_id" bigint NOT NULL,
  "showtime_id" bigint NOT NULL,
  "total_price" int NOT NULL,
  "payment_method_id" bigint NOT NULL,
  "status" "statusOrder" DEFAULT 'pending',
  "expired_at" timestamp NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

COMMENT ON COLUMN "orders"."payment_reference" IS 'virtual account, link e-wallet, null if not yet';

COMMENT ON COLUMN "orders"."expired_at" IS 'created_at + 1 hours';

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD FOREIGN KEY ("showtime_id") REFERENCES "showtimes" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD FOREIGN KEY ("payment_method_id") REFERENCES "payment_methods" ("id") DEFERRABLE INITIALLY IMMEDIATE;
