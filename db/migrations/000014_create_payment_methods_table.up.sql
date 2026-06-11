CREATE TABLE "payment_methods" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "logo" varchar NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp
);
