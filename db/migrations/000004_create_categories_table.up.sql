CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE INDEX ON "categories" ("name");
