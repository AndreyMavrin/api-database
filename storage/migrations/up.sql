CREATE UNLOGGED TABLE "users" (
  "about" varchar NOT NULL,
  "email" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "nickname" varchar PRIMARY KEY
);

CREATE UNLOGGED TABLE "forums" (
  "username" varchar NOT null,
  "posts" BIGINT DEFAULT 0,
  "threads" int DEFAULT 0,
  "slug" varchar PRIMARY KEY,
  "title" varchar NOT NULL,
  FOREIGN KEY ("username") REFERENCES "users" (nickname)
);

CREATE UNLOGGED TABLE "threads" (
  "author" varchar NOT NULL,
  "created" timestamp with time zone default now(),
  "forum" varchar NOT NULL,
  "id" SERIAL PRIMARY KEY,
  "message" varchar NOT NULL,
  "slug" varchar UNIQUE,
  "title" varchar NOT NULL,
  "votes" int DEFAULT 0,
  FOREIGN KEY (author) REFERENCES "users" (nickname),
  FOREIGN KEY (forum) REFERENCES "forums" (slug)
);

CREATE UNLOGGED TABLE "posts" (
  "author" varchar NOT NULL,
  "created" timestamp with time zone default now(),
  "forum" varchar NOT NULL,
  "id" BIGSERIAL PRIMARY KEY,
  "is_edited" BOOLEAN DEFAULT false,
  "message" varchar NOT NULL,
  "parent" BIGINT DEFAULT 0,
  "thread" int,
  "path" BIGINT[] DEFAULT ARRAY []::INTEGER[],
  
  FOREIGN KEY (author) REFERENCES "users" (nickname),
  FOREIGN KEY (forum) REFERENCES "forums" (slug),
  FOREIGN KEY (thread) REFERENCES "threads" (id),
  FOREIGN KEY (parent) REFERENCES "posts" (id)
);

CREATE UNLOGGED TABLE "votes" (
  "nickname" varchar NOT NULL,
  "voice" int,
  "thread" int,
  
   FOREIGN KEY (nickname) REFERENCES "users" (nickname),
   FOREIGN KEY (thread) REFERENCES "threads" (id),
   UNIQUE (nickname, thread)
);

CREATE OR REPLACE FUNCTION update_threads_count() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE forums SET Threads=(Threads+1) WHERE LOWER(slug)=LOWER(NEW.forum);
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent IS NULL) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT path FROM posts WHERE id = new.parent INTO parent_path;
        SELECT thread FROM posts WHERE id = parent_path[1] INTO first_parent_thread;
        IF NOT FOUND OR first_parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'parent is from different thread';
        end if;

        NEW.path := NEW.path || parent_path || new.id;
    end if;
    UPDATE forums SET Posts=Posts + 1 WHERE lower(forums.slug) = lower(new.forum);
    RETURN new;
end
$update_path$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION insert_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE threads SET votes=(votes+NEW.voice) WHERE id=NEW.thread;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE threads SET votes=(votes+NEW.voice*2) WHERE id=NEW.thread;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;


CREATE TRIGGER add_thread_to_forum
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_threads_count();

CREATE TRIGGER path_update_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();

CREATE TRIGGER add_vote
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes();

CREATE TRIGGER edit_vote
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();

CREATE INDEX post_thread_id_index ON posts (thread, id);
CREATE INDEX post_path_id_index ON posts (id, (posts.path));
CREATE INDEX post_author_index ON posts (lower(author));
CREATE INDEX post_forum_index ON posts (lower(forum));

CREATE INDEX forum_slug_lower_index ON forums (lower(forums.Slug));

CREATE INDEX users_nickname_lower_index ON users (lower(users.nickname));
CREATE INDEX users_email_index ON users (lower(users.email));

CREATE INDEX thread_id_index ON threads (id);

CREATE INDEX vote_nickname ON votes (nickname, thread);