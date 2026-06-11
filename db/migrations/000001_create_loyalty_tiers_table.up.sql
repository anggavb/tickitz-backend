CREATE TABLE "loyalty_tiers" (
  "id" bigserial PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "min_point" int NOT NULL
);
