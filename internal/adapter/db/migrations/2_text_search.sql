-- +migrate Up
-- +migrate StatementBegin
DO
$$BEGIN
    CREATE TEXT SEARCH CONFIGURATION custom_english ( COPY = pg_catalog.english );
    CREATE TEXT SEARCH DICTIONARY custom_english_dict (
        Template = snowball,
        Language = english
    );

    ALTER TEXT SEARCH CONFIGURATION custom_english
    ALTER MAPPING FOR asciiword, asciihword, hword_asciipart,
                      word, hword, hword_part, numword
    WITH custom_english_dict;
EXCEPTION
   WHEN unique_violation THEN
      NULL;  -- ignore error
END;$$;
-- +migrate StatementEnd
-- +migrate Down
DROP TEXT SEARCH CONFIGURATION IF EXISTS custom_english;
DROP TEXT SEARCH DICTIONARY custom_english_dict;