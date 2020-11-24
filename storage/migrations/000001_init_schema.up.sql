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
  "user" varchar NOT NULL
);