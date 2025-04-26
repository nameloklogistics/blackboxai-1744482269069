#!/bin/bash

# Disaster Recovery Script for Logistics Marketplace

set -e

BACKUP_DIR="/var/backups/logistics_marketplace"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_PATH="$BACKUP_DIR/backup_$TIMESTAMP"

echo "Starting disaster recovery backup..."

# Create backup directory if not exists
mkdir -p "$BACKUP_DIR"

# Backup database (assuming PostgreSQL, adjust as needed)
echo "Backing up PostgreSQL database..."
PGHOST=${PGHOST:-localhost}
PGPORT=${PGPORT:-5432}
PGUSER=${PGUSER:-postgres}
PGDATABASE=${PGDATABASE:-logistics_db}
PGPASSWORD=${PGPASSWORD:-}

export PGPASSWORD

pg_dump -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -F c -b -v -f "$BACKUP_PATH/db_backup.dump" "$PGDATABASE"

# Backup application files (adjust paths as needed)
echo "Backing up application files..."
tar czf "$BACKUP_PATH/app_files.tar.gz" ./cmd ./internal ./pkg ./frontend

# Verify backups
if [ -f "$BACKUP_PATH/db_backup.dump" ] && [ -f "$BACKUP_PATH/app_files.tar.gz" ]; then
  echo "Backup completed successfully at $BACKUP_PATH"
else
  echo "Backup failed!"
  exit 1
fi

# Optional: Upload backup to remote storage (e.g., AWS S3, Google Cloud Storage)
# echo "Uploading backup to remote storage..."
# aws s3 cp "$BACKUP_PATH" s3://your-bucket-name/ --recursive

echo "Disaster recovery backup process finished."
