# Deployment

Deploys this service to a plain Ubuntu/Debian VPS using Docker Compose: Postgres,
the Go backend, and nginx (serving the built frontend + reverse-proxying `/api/*` to
the backend) in one stack, with Let's Encrypt for HTTPS.

## 1. Provision the server

Any VPS with a public IP works (DigitalOcean, Hetzner, EC2, ...). Minimum ~1GB RAM.

```bash
# On the server, as a non-root sudo user:
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker "$USER"   # log out/in once for this to take effect

sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

## 2. Get the code and configure it

```bash
git clone https://github.com/PurePras/deliveryService.git
cd deliveryService
cp .env.production.example .env
```

Edit `.env`:
- `POSTGRES_PASSWORD` — a real password (also update it inside `DATABASE_URL`).
- `JWT_SECRET` — generate with `openssl rand -base64 32`.
- `DOMAIN` / `CORS_ALLOWED_ORIGIN` — your domain, once you have one pointed at this
  server's IP (an A record). Without a domain yet, leave `CORS_ALLOWED_ORIGIN` as the
  server's IP over http and skip step 4 below until you have one.
- `LETSENCRYPT_EMAIL` — used for Let's Encrypt renewal/expiry notices.

## 3. First boot (HTTP only)

```bash
docker compose -f docker-compose.prod.yml up -d --build db backend nginx
docker compose -f docker-compose.prod.yml logs -f backend   # confirm it started cleanly
curl http://localhost/api/health
```

This alone is a fully working deployment over plain HTTP — everything through step 3
works with no domain name at all. Do step 4 whenever a domain is ready.

## 4. Enable HTTPS (once DNS for `$DOMAIN` points at this server)

```bash
./scripts/init-letsencrypt.sh
```

This bootstraps a real Let's Encrypt certificate and switches nginx to the HTTPS
config (see `nginx/app-ssl.conf.template`), then starts the `certbot` service, which
keeps the certificate renewed for as long as the stack runs. After this, also update
`.env`: set `CORS_ALLOWED_ORIGIN=https://$DOMAIN`, then
`docker compose -f docker-compose.prod.yml up -d backend` to pick it up.

## 5. Deploying updates

```bash
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

The backend applies pending migrations itself on startup (`database.RunMigrations` in
`cmd/api/main.go`) — no separate migration step needed for a normal deploy.

## 6. Database backups

```bash
./scripts/backup-db.sh          # dumps to ./data/backups, prunes anything over 14 days old
./scripts/restore-db.sh <file>  # restores a dump — destructive, asks for confirmation
```

Automate nightly backups with cron (`crontab -e`):

```cron
0 2 * * * cd /path/to/deliveryService && ./scripts/backup-db.sh >> /var/log/shri-ram-backup.log 2>&1
```

Copy `./data/backups` off the server regularly (e.g. `rsync` to another host) — backups
sitting only on the same disk as the database don't survive that disk failing.

## 7. Monitoring

The backend exposes Prometheus metrics at `/metrics` (request count and latency by
route — see `internal/middleware/metrics.go`), reachable only from other containers on
the compose network, never through nginx or the internet.

```bash
docker compose -f docker-compose.prod.yml -f docker-compose.monitoring.yml up -d
ssh -L 3000:localhost:3000 <user>@<server>   # then open http://localhost:3000
```

Grafana comes pre-provisioned with a Prometheus datasource and a starter "Shri Ram
API" dashboard (request rate, error rate, p95 latency). Default login is
`admin` / `admin` (from `docker-compose.monitoring.yml`'s `GF_SECURITY_ADMIN_PASSWORD`)
— change it on first login.

`docker compose -f docker-compose.prod.yml ps` shows each service's container-level
health (backend via its `-healthcheck` mode, db via `pg_isready`) for a quick check
that doesn't need the full stack.

## 8. CI

`.github/workflows/ci.yml` runs on every push/PR to `main`: backend build, vet, unit
tests, and the integration suite against a real Postgres service container; frontend
lint, typecheck, and build. It's verification only — nothing here deploys
automatically. Come back to this once you're ready to wire up CD.
