DROP TRIGGER IF EXISTS verify_updated_at ON public."banner";
DROP FUNCTION IF EXISTS updated_at_now;

DROP TABLE IF EXISTS public."user" CASCADE;
DROP TABLE IF EXISTS public."banner" CASCADE;
DROP TABLE IF EXISTS public."feature" CASCADE;
DROP TABLE IF EXISTS public."tag" CASCADE;
DROP TABLE IF EXISTS public."banner_tag" CASCADE;

DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS banner_id_seq;
DROP SEQUENCE IF EXISTS feature_id_seq;
DROP SEQUENCE IF EXISTS tag_id_seq;
DROP SEQUENCE IF EXISTS banner_tag_id_seq;