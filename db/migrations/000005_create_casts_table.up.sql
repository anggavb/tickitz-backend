CREATE TABLE "casts" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE INDEX ON "casts" ("name");
