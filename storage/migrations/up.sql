CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE "users" (
  "about" varchar,
  "email" CITEXT UNIQUE,
  "fullname" varchar NOT NULL,
  "nickname" CITEXT PRIMARY KEY
);

CREATE UNLOGGED TABLE "forums" (
  "username" CITEXT NOT null,
  "posts" BIGINT DEFAULT 0,
  "threads" int DEFAULT 0,
  "slug" CITEXT PRIMARY KEY,
  "title" varchar NOT NULL,
  FOREIGN KEY ("username") REFERENCES "users" (nickname)
);

CREATE UNLOGGED TABLE "threads" (
  "author" CITEXT NOT NULL,
  "created" timestamp with time zone default now(),
  "forum" CITEXT NOT NULL,
  "id" SERIAL PRIMARY KEY,
  "message" varchar NOT NULL,
  "slug" CITEXT UNIQUE,
  "title" varchar NOT NULL,
  "votes" int DEFAULT 0,
  FOREIGN KEY (author) REFERENCES "users" (nickname),
  FOREIGN KEY (forum) REFERENCES "forums" (slug)
);

CREATE UNLOGGED TABLE "posts" (
  "author" CITEXT NOT NULL,
  "created" timestamp with time zone default now(),
  "forum" CITEXT NOT NULL,
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
  "nickname" CITEXT,
  "voice" int,
  "thread" int,
  
   FOREIGN KEY (nickname) REFERENCES "users" (nickname),
   FOREIGN KEY (thread) REFERENCES "threads" (id),
   UNIQUE (nickname, thread)
);

CREATE UNLOGGED TABLE users_forum
(
    nickname citext NOT NULL,
    Slug     citext NOT NULL,
    FOREIGN KEY (nickname) REFERENCES "users" (nickname),
    FOREIGN KEY (Slug) REFERENCES "forums" (Slug),
    UNIQUE (nickname, Slug)
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
begin
	IF OLD.voice <> NEW.voice THEN
    	UPDATE threads SET votes=(votes+NEW.voice*2) WHERE id=NEW.thread;
    END IF;
    return NEW;
end
$update_users_forum$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_user_forum() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    INSERT INTO users_forum (nickname, Slug) VALUES (NEW.author, NEW.forum) on conflict do nothing;
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

CREATE TRIGGER thread_insert_user_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();

CREATE TRIGGER post_insert_user_forum
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();

CREATE INDEX user_nickname ON users using hash (nickname);
CREATE INDEX user_email ON users using hash (email);

CREATE INDEX forum_slug ON forums using hash (slug);

CREATE UNIQUE INDEX forum_users_unique on users_forum (slug, nickname);
CLUSTER users_forum USING forum_users_unique;

CREATE INDEX thr_slug ON threads using hash (slug);
CREATE INDEX thr_date ON threads (created);
CREATE INDEX thr_forum ON threads using hash (forum);
CREATE INDEX thr_forum_date ON threads (forum, created);

CREATE INDEX post_id_path on posts (id, (path[1]));
CREATE INDEX post_thread_id_path1_parent on posts (thread, id, (path[1]), parent);
CREATE INDEX post_thread_path_id on posts (thread, path, id);
CREATE INDEX post_path1 on posts ((path[1]));
CREATE INDEX post_thread_id on posts (thread, id);
CREATE INDEX post_thr_id ON posts (thread);

CREATE UNIQUE INDEX vote_unique on votes (nickname, thread);