#!/usr/bin/env bash
# One-time setup to get nginx serving over HTTPS with a real Let's Encrypt certificate.
#
# Run this from the server, once DNS for $DOMAIN already points at it, after
# `docker compose -f docker-compose.prod.yml up -d` is already running on plain HTTP.
# See docs/DEPLOYMENT.md for the full sequence — this script only handles the
# chicken-and-egg part: nginx can't start with `ssl_certificate` pointing at a file
# that doesn't exist yet, and certbot can't get a real cert without a running nginx to
# answer the HTTP-01 challenge. So: fake a cert just so nginx will start, ask certbot
# for the real one through it, then swap in the real HTTPS config and reload.
set -euo pipefail

cd "$(dirname "$0")/.."
COMPOSE="docker compose -f docker-compose.prod.yml"

if [ -f .env ]; then
	set -a
	# shellcheck disable=SC1091
	source .env
	set +a
fi

: "${DOMAIN:?Set DOMAIN in .env first (see .env.production.example)}"
: "${LETSENCRYPT_EMAIL:?Set LETSENCRYPT_EMAIL in .env first}"

echo "==> Creating a dummy self-signed cert so nginx can start"
# certbot/certbot is Alpine-based but doesn't ship the openssl CLI itself, hence the
# apk add — only needed for this one throwaway cert.
$COMPOSE run --rm --entrypoint sh certbot -c "
  apk add --no-cache openssl >/dev/null &&
  mkdir -p /etc/letsencrypt/live/$DOMAIN &&
  openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
    -keyout /etc/letsencrypt/live/$DOMAIN/privkey.pem \
    -out /etc/letsencrypt/live/$DOMAIN/fullchain.pem \
    -subj '/CN=localhost'
"

echo "==> Installing the HTTPS server config and starting nginx with the dummy cert"
sed "s/\${DOMAIN}/$DOMAIN/g" nginx/app-ssl.conf.template > nginx/conf.d/app.conf
$COMPOSE up -d nginx

echo "==> Deleting the dummy cert so certbot will issue a real one in its place"
$COMPOSE run --rm --entrypoint sh certbot -c "rm -rf /etc/letsencrypt/live/$DOMAIN /etc/letsencrypt/archive/$DOMAIN /etc/letsencrypt/renewal/$DOMAIN.conf"

echo "==> Requesting the real certificate from Let's Encrypt"
$COMPOSE run --rm certbot certonly \
	--webroot -w /var/www/certbot \
	--email "$LETSENCRYPT_EMAIL" -d "$DOMAIN" \
	--rsa-key-size 2048 --agree-tos --non-interactive

echo "==> Reloading nginx with the real certificate"
$COMPOSE exec nginx nginx -s reload

echo "==> Starting the renewal loop"
$COMPOSE up -d certbot

echo "Done. https://$DOMAIN should now be serving over HTTPS."
