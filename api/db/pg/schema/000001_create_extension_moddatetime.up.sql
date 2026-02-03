--- https://stackoverflow.com/a/76092474
DO $SETUP_MODDATETIME$
BEGIN
  CREATE EXTENSION IF NOT EXISTS "moddatetime";
EXCEPTION WHEN INSUFFICIENT_PRIVILEGE THEN
  -- Unable to use moddatetime extension, so create a function instead
  CREATE OR REPLACE FUNCTION moddatetime() RETURNS trigger LANGUAGE plpgsql AS $moddatetime$
    DECLARE
      colname name;
    BEGIN
      IF (TG_NARGS = 1) THEN
        colname = TG_ARGV[0];
      ELSE
        RAISE EXCEPTION 'moddatetime(colname) requires one argument';
      END IF;

      RETURN json_populate_record(NEW, json_build_object(colname, NOW()));
    END;
  $moddatetime$;
END $SETUP_MODDATETIME$;
