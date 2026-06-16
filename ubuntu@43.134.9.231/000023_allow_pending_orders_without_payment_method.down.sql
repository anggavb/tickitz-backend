UPDATE "orders"
SET "payment_method_id" = 1
WHERE "payment_method_id" IS NULL;

ALTER TABLE "orders"
ALTER COLUMN "payment_method_id" SET NOT NULL;
