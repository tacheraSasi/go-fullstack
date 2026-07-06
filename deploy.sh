#!/usr/bin/env sh
set -e

APP_NAME="go-fullstack"
COMPOSE_FILE="compose.yaml"

echo "==> Deploying $APP_NAME ..."

# 1. Pull latest code (uncomment if you want this step)
# echo "==> Pulling latest code..."
# git pull origin main

# 2. Build and start (or restart) services in detached mode
echo "==> Building and starting services..."
docker compose -f "$COMPOSE_FILE" up --build -d

# 3. Wait a moment and check health
echo "==> Waiting for service to become healthy..."
sleep 3

# Poll the health endpoint up to 10 times
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
echo "   Check the logs with: docker compose -f $COMPOSE_FILE logs"
exit 1
