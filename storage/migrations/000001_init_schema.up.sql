CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "about" varchar NOT NULL,
  "email" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "nickname" varchar NOT NULL
);

CREATE TABLE "forums" (
  "id" SERIAL PRIMARY KEY,
  "slug" varchar NOT NULL,
  "title" varchar NOT NULL,
  "username" varchar NOT NULL
);

CREATE TABLE "threads" (
  "id" SERIAL PRIMARY KEY,
  "author" varchar NOT NULL,
  "created" varchar NOT NULL,
  "forum" varchar NOT NULL,
  "message" varchar NOT NULL,
  "slug" varchar NOT NULL,
  "title" varchar NOT NULL
);