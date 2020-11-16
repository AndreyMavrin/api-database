CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "about" varchar NOT NULL,
  "email" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "nickname" varchar NOT NULL
);