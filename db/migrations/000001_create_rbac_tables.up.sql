CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE "profiles" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "profile_picture" varchar NOT NULL,
  "profile_background" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "mobile" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "username" varchar UNIQUE NOT NULL,
  "hashedpassword" varchar NOT NULL,
  "profile" uuid NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "roles" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "role" varchar UNIQUE NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "account_roles" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "account_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "task" varchar UNIQUE NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "role_tasks" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "task_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "helptext" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "task_id" uuid NOT NULL,
  "helptext" varchar UNIQUE NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "menu" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar UNIQUE NOT NULL,
  "task_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "navigation" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" varchar UNIQUE NOT NULL,
  "task_id" uuid NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("profile") REFERENCES "profiles" ("id");

ALTER TABLE "account_roles" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "account_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_tasks" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "role_tasks" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "helptext" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "menu" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "navigation" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

CREATE UNIQUE INDEX ON "account_roles" ("account_id", "role_id");

CREATE UNIQUE INDEX ON "role_tasks" ("task_id", "role_id");
