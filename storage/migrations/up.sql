CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE "users" (
  "about" TEXT,
  "email" CITEXT UNIQUE,
  "fullname" TEXT NOT NULL,
  "nickname" CITEXT PRIMARY KEY
);

CREATE UNLOGGED TABLE "forums" (
  "username" CITEXT NOT null,
  "posts" BIGINT DEFAULT 0,
  "threads" int DEFAULT 0,
  "slug" CITEXT PRIMARY KEY,
  "title" TEXT NOT NULL,
  FOREIGN KEY ("username") REFERENCES "users" (nickname)
);

CREATE UNLOGGED TABLE "threads" (
  "author" CITEXT NOT NULL,
  "created" timestamp with time zone default now(),
  "forum" CITEXT NOT NULL,
  "id" SERIAL PRIMARY KEY,
  "message" TEXT NOT NULL,
  "slug" CITEXT UNIQUE,
  "title" TEXT NOT NULL,
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
  "message" TEXT NOT NULL,
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
    fullname TEXT NOT NULL,
    about    TEXT,
    email    CITEXT,
    slug     citext NOT NULL,
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
DECLARE
    m_fullname CITEXT;
    m_about    CITEXT;
    m_email CITEXT;
BEGIN
    SELECT fullname, about, email FROM users WHERE nickname = NEW.author INTO m_fullname, m_about, m_email;
    INSERT INTO users_forum (nickname, fullname, about, email, slug)
    VALUES (NEW.author, m_fullname, m_about, m_email, NEW.forum) on conflict do nothing;
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

CREATE INDEX post_first_parent_thread_index ON posts ((posts.path[1]), thread);
CREATE INDEX post_first_parent_id_index ON posts ((posts.path[1]), id);
CREATE INDEX post_first_parent_index ON posts ((posts.path[1]));
CREATE INDEX post_path_index ON posts ((posts.path));
CREATE INDEX post_thread_index ON posts (thread);
CREATE INDEX post_thread_id_index ON posts (thread, id);
CREATE INDEX post_path_id_index ON posts (id, (posts.path));

CREATE INDEX forum_slug_lower_index ON forums (lower(forums.Slug));

CREATE INDEX users_nickname_lower_index ON users (lower(users.nickname));
CREATE INDEX users_email_index ON users (lower(users.email));

CREATE INDEX users_forum_forum_user_index ON users_forum (lower(users_forum.Slug), nickname);
CREATE INDEX users_forum_user_index ON users_forum (nickname);
CREATE INDEX users_forum_forum_index ON users_forum ((users_forum.Slug));

CREATE INDEX thread_slug_lower_index ON threads (lower(slug));
CREATE INDEX thread_slug_index ON threads (slug);
CREATE INDEX thread_slug_id_index ON threads (lower(slug), id);
CREATE INDEX thread_forum_lower_index ON threads (lower(forum));
CREATE INDEX thread_id_forum_index ON threads (id, forum);
CREATE INDEX thread_created_index ON threads (created);

CREATE INDEX vote_nickname ON votes (lower(nickname), thread, voice);