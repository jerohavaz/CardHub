
#!/bin/sh
set -eu

PRIVATE_KEY="/keys/private/jwt-private.pem"
PUBLIC_KEY="/keys/public/jwt-public.pem"

# Distroless non-root UID/GID
APP_UID=65532
APP_GID=65532

umask 077

# Keep existing keys across container restarts.
if [ -f "$PRIVATE_KEY" ] && [ -f "$PUBLIC_KEY" ]; then
    echo "JWT keys already initialized."

elif [ -e "$PRIVATE_KEY" ] || [ -e "$PUBLIC_KEY" ]; then
    echo "ERROR: Incomplete JWT key pair." >&2
    exit 1

else
    echo "Generating Ed25519 JWT keys..."

    openssl genpkey \
        -algorithm ED25519 \
        -out "$PRIVATE_KEY"

    openssl pkey \
        -in "$PRIVATE_KEY" \
        -pubout \
        -out "$PUBLIC_KEY"
fi

# Verify that the public key matches the private key.
openssl pkey \
    -in "$PRIVATE_KEY" \
    -pubout \
    -out /tmp/jwt-public-check.pem

if ! cmp -s /tmp/jwt-public-check.pem "$PUBLIC_KEY"; then
    echo "ERROR: JWT key pair does not match." >&2
    exit 1
fi

rm /tmp/jwt-public-check.pem

# Private volume: readable by Identity only.
chown "$APP_UID:$APP_GID" "$PRIVATE_KEY"
chmod 600 "$PRIVATE_KEY"

# Public volume: readable by Gateway.
chown "$APP_UID:$APP_GID" "$PUBLIC_KEY"
chmod 644 "$PUBLIC_KEY"

echo "JWT keys ready."
