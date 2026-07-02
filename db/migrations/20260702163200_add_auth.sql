-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "email" text NOT NULL,
  "password_hash" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "token" text NOT NULL,
  "user_id" bigint NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("token"),
  CONSTRAINT "refresh_tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Modify "reviews" table
ALTER TABLE "public"."reviews" ADD COLUMN "user_id" bigint NOT NULL, ADD CONSTRAINT "reviews_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create index "reviews_user_id_idx" to table: "reviews"
CREATE INDEX "reviews_user_id_idx" ON "public"."reviews" ("user_id");
