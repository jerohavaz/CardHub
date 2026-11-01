#!/bin/sh
set -eu
BASE="${CARDHUB_BASE_URL:-http://127.0.0.1:8080}"

printf 'Checking gateway health...\n'
i=0
until curl --fail --silent "$BASE/healthz" >/dev/null; do
  i=$((i + 1))
  if [ "$i" -ge 30 ]; then echo 'Gateway did not become ready' >&2; exit 1; fi
  sleep 1
done
printf 'Checking generated OpenAPI...\n'
curl --fail --silent "$BASE/openapi.json" | grep -q '"/v1/examples"'
printf 'Creating example...\n'
i=0
created=""
until created="$(curl --fail --silent -X POST "$BASE/v1/examples" -H 'Content-Type: application/json' -d '{"name":"Pikachu"}')"; do
  i=$((i + 1))
  if [ "$i" -ge 30 ]; then echo 'Example service unavailable' >&2; exit 1; fi
  sleep 1
done
printf '%s\n' "$created"
id="$(printf '%s' "$created" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')"
if [ -z "$id" ]; then echo 'Creation did not return id' >&2; exit 1; fi
printf 'Fetching %s...\n' "$id"
fetched="$(curl --fail --silent "$BASE/v1/examples/$id")"
printf '%s\n' "$fetched"
printf '%s' "$fetched" | grep -q '"name":"Pikachu"'
# Check status propagation for a valid UUID that does not exist.
missing="$(curl --silent -o /dev/null -w '%{http_code}' "$BASE/v1/examples/00000000-0000-0000-0000-000000000000")"
[ "$missing" = 404 ] || { echo "Expected 404, received $missing" >&2; exit 1; }
printf 'Smoke test passed.\n'
