CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "hashed_password" varchar NOT NULL,
    "phone" varchar DEFAULT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- ALTER TABLE "users" ADD CONSTRAINT "email_unqiue" UNIQUE ("email");

CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "email" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "expires_at" timestamptz NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "is_logged_out" boolean NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- ALTER TABLE "sessions" ADD FOREIGN KEY ("email") REFERENCES "users" ("email");

CREATE TABLE "channels" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "memberships" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigserial NOT NULL REFERENCES users(id),
    "channel_id" bigserial NOT NULL REFERENCES channels(id),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);