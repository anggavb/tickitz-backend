CREATE TABLE "cinemas" (
  "id" bigserial PRIMARY KEY,
  "location_id" bigint NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "cinemas" ADD FOREIGN KEY ("location_id") REFERENCES "location" ("id") DEFERRABLE INITIALLY IMMEDIATE;
