#!/usr/bin/env sh
set -e

APP_NAME="go-fullstack"
MODE="${1:-docker}"   # pass "pm2" to deploy without Docker

# ── PM2 / bare-metal deploy ──────────────────────────────────────────────
if [ "$MODE" = "pm2" ]; then
  echo "==> Deploying $APP_NAME via Go build + PM2 ..."

  # 1. Ensure required directories exist
  mkdir -p bin data logs

  # 2. Build the binary
  echo "==> Building binary..."
  CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/api ./cmd/app/main.go

  # 3. Create a .env if it doesn't exist (won't overwrite existing)
  if [ ! -f .env ]; then
    cat > .env << EOF
GIN_MODE=release
SERVER_PORT=8080
DB_TYPE=sqlite
DB_PATH=./data/core.db
JWT_SECRET=change-this-to-a-random-secret
JWT_EXPIRES_IN=24
LOG_FILE_PATH=./logs/app.log
EOF
    echo "==> Created default .env file — edit JWT_SECRET before production use"
  fi

  # 4. Start / restart via PM2
  if pm2 describe "$APP_NAME" > /dev/null 2>&1; then
    echo "==> Restarting PM2 process..."
    pm2 restart "$APP_NAME"
  else
    echo "==> Starting PM2 process..."
    pm2 start ecosystem.config.js --env production
  fi

  # 5. Save PM2 process list so it resurrects on server reboot
  pm2 save

  echo "==> Done! $APP_NAME is running on http://localhost:8080"
  echo "   Manage it:  pm2 status | pm2 logs $APP_NAME | pm2 stop $APP_NAME"
  exit 0
fi

# ── Docker deploy (default) ──────────────────────────────────────────────
echo "==> Deploying $APP_NAME via Docker..."

docker compose -f compose.yaml up --build -d

echo "==> Waiting for service to become healthy..."
sleep 3

i=0
while [ $i -lt 10 ]; do
  status=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health 2>/dev/null || echo "000")
  if [ "$status" = "200" ]; then
    echo "==> $APP_NAME is healthy! (HTTP $status)"
    exit 0
  fi
  printf "   Attempt %d/10 – status %s, retrying...\n" $((i + 1)) "$status"
  sleep 2
  i=$((i + 1))
done

echo "==> WARNING: Health check did not pass within 20 seconds."
echo "   Check the logs with: docker compose logs"
exit 1
