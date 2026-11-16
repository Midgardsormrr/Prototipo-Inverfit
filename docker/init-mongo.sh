#!/bin/bash
set -eo pipefail

echo "[inverfit-init] Starting MongoDB initialization script..."

APP_USER="${MONGO_APP_USERNAME:-inverfit_app}"
APP_PASS="${MONGO_APP_PASSWORD:-inverfit_app_pass_2026}"
APP_DB="${MONGODB_DATABASE:-inverfit_container}"
ROOT_USER="${MONGO_INITDB_ROOT_USERNAME:-inverfit_admin}"
ROOT_PASS="${MONGO_INITDB_ROOT_PASSWORD:-inverfit_root_pass_2026}"

echo "[inverfit-init] Ensuring application user '$APP_USER' exists for database '$APP_DB'..."

mongosh -u "$ROOT_USER" -p "$ROOT_PASS" --authenticationDatabase admin "$APP_DB" <<EOF
const existingUser = db.getUser("$APP_USER");
if (!existingUser) {
    db.createUser({
        user: "$APP_USER",
        pwd: "$APP_PASS",
        roles: [
            { role: "readWrite", db: "$APP_DB" }
        ]
    });
    print("[inverfit-init] Application user '$APP_USER' created successfully for database '$APP_DB'.");
} else {
    print("[inverfit-init] Application user '$APP_USER' already exists in database '$APP_DB'.");
}
EOF

echo "[inverfit-init] MongoDB initialization completed successfully."
