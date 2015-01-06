-- +goose Up
-- +goose StatementBegin
-- We have to remove the default value of '{}'::json because it's incompatible with jsonb
ALTER TABLE reading ALTER COLUMN readings DROP DEFAULT;
-- Change the column type, and cast existing data to jsonb
ALTER TABLE reading ALTER COLUMN readings TYPE jsonb USING readings::jsonb;
-- Add the default value back, of course using jsonb this time
ALTER TABLE reading ALTER COLUMN readings SET DEFAULT '{}'::jsonb;
-- SQL function to add a key to jsonb field
CREATE OR REPLACE FUNCTION "json_object_set_key"(
  "jsonb"          jsonb,
  "key_to_set"    TEXT,
  "value_to_set"  anyelement
)
  RETURNS jsonb
  LANGUAGE sql
  IMMUTABLE
  STRICT
AS $function$
SELECT COALESCE(
  (SELECT ('{' || string_agg(to_json("key") || ':' || "value", ',') || '}')
     FROM (SELECT *
             FROM jsonb_each("jsonb")
            WHERE "key" <> "key_to_set"
            UNION ALL
           SELECT "key_to_set", "value_to_set") AS "fields"),
  '{}'
)::jsonb
$function$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- We have to remove the default value of '{}'::json because it's incompatible with jsonb
ALTER TABLE reading ALTER COLUMN readings DROP DEFAULT;
-- Change the column type, and cast existing data to jsonb
ALTER TABLE reading ALTER COLUMN readings TYPE json USING readings::json;
-- Add the default value back, of course using jsonb this time
ALTER TABLE reading ALTER COLUMN readings SET DEFAULT '{}'::json;
-- SQL function to add a key to json field
CREATE OR REPLACE FUNCTION "json_object_set_key"(
  "json"          json,
  "key_to_set"    TEXT,
  "value_to_set"  anyelement
)
  RETURNS json
  LANGUAGE sql
  IMMUTABLE
  STRICT
AS $function$
SELECT COALESCE(
  (SELECT ('{' || string_agg(to_json("key") || ':' || "value", ',') || '}')
     FROM (SELECT *
             FROM json_each("json")
            WHERE "key" <> "key_to_set"
            UNION ALL
           SELECT "key_to_set", to_json("value_to_set")) AS "fields"),
  '{}'
)::json
$function$;
-- +goose StatementEnd
