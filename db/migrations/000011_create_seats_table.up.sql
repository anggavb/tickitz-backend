CREATE TYPE "seatType" AS ENUM (
  'regular',
  'love_nest'
);

CREATE TABLE "seats" (
  "id" bigserial PRIMARY KEY,
  "cinema_id" bigint NOT NULL,
  "row" char(1) NOT NULL,
  "number" int NOT NULL,
  "type" "seatType" NOT NULL DEFAULT 'regular',
  "created_at" timestamp DEFAULT (now())
);

COMMENT ON COLUMN "seats"."row" IS 'A,B,C,D,etc.';

COMMENT ON COLUMN "seats"."number" IS '1-14';

ALTER TABLE "seats" ADD FOREIGN KEY ("cinema_id") REFERENCES "cinemas" ("id") DEFERRABLE INITIALLY IMMEDIATE;
