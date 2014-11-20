-- +goose Up
-- +goose StatementBegin
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

-- +goose Down
DROP FUNCTION json_object_set_key(json,TEXT,anyelement);
