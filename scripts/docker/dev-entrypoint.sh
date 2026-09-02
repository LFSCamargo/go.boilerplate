#!/bin/sh
set -e

# Email npm install is deferred so Air can bind :5000 immediately.
# Recover/verify emails need tsx; install later with: docker compose exec app sh -c 'cd emails && npm install'
if [ ! -x /app/emails/node_modules/.bin/tsx ]; then
	echo "email workspace not installed yet; SMTP templates will fail until npm install runs"
fi

exec air -c /app/.air.toml
