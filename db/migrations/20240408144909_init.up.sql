CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE SEQUENCE IF NOT EXISTS banner_id_seq;
CREATE SEQUENCE IF NOT EXISTS feature_id_seq;
CREATE SEQUENCE IF NOT EXISTS tag_id_seq;
CREATE SEQUENCE IF NOT EXISTS banner_tag_id_seq;

CREATE TABLE IF NOT EXISTS public."user"
(
    id         BIGINT                   DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    login      TEXT UNIQUE                                                       NOT NULL CHECK (login <> '')
        CONSTRAINT max_len_email CHECK (LENGTH(login) <= 256),
    password   TEXT                                                              NOT NULL CHECK (password <> '')
        CONSTRAINT max_len_password CHECK (LENGTH(password) <= 256),
    is_admin   BOOL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()                            NOT NULL
);

CREATE TABLE IF NOT EXISTS public."feature"
(
    id           BIGINT                   DEFAULT NEXTVAL('feature_id_seq'::regclass)  NOT NULL PRIMARY KEY,
    title        TEXT                                                                 NOT NULL CHECK (title <> '')
        CONSTRAINT   max_len_title CHECK (LENGTH(title) <= 150),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL
);

CREATE TABLE IF NOT EXISTS public."tag"
(
    id           BIGINT                   DEFAULT NEXTVAL('tag_id_seq'::regclass)  NOT NULL PRIMARY KEY,
    title        TEXT                                                                 NOT NULL CHECK (title <> '')
        CONSTRAINT   max_len_title CHECK (LENGTH(title) <= 150),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL
);

CREATE TABLE IF NOT EXISTS public."banner"
(
    id           BIGINT                   DEFAULT NEXTVAL('banner_id_seq'::regclass)  NOT NULL PRIMARY KEY,
    author_id    BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    feature_id   BIGINT                                                               NOT NULL REFERENCES public."feature" (id),
    title        TEXT                                                                 NOT NULL CHECK (title <> '')
        CONSTRAINT   max_len_title CHECK (LENGTH(title) <= 150),
    text         TEXT                                                                 NOT NULL CHECK (text <> '')
        CONSTRAINT   max_len_text CHECK (LENGTH(text) <= 1000),
    url          TEXT                                                                 NOT NULL CHECK (url <> '')
        CONSTRAINT max_len_url CHECK (LENGTH(url) <= 256),
    is_active    BOOL                     DEFAULT TRUE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL
);

CREATE TABLE IF NOT EXISTS public."banner_tag"
(
    id           BIGINT                   DEFAULT NEXTVAL('banner_tag_id_seq'::regclass)  NOT NULL PRIMARY KEY,
    banner_id    BIGINT                                                               NOT NULL REFERENCES public."banner" (id),
    tag_id       BIGINT                                                               NOT NULL REFERENCES public."tag" (id)
);

CREATE OR REPLACE FUNCTION updated_at_now()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS verify_updated_at ON public."banner";
CREATE TRIGGER verify_updated_at
    BEFORE UPDATE
    ON public."banner"
    FOR EACH ROW
EXECUTE PROCEDURE updated_at_now();