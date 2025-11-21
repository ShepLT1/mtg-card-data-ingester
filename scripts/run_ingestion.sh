#!/bin/bash
# scripts/run_ingestion.sh

# Exit immediately if a command fails
set -e

# Optional: log timestamp
echo "Running ingestion at $(date -u)" >> /var/log/cron.log

# Call your ingestion service endpoint
# Assuming your ingestion service is reachable at http://data-ingester:8085/ingest
curl -s -X POST http://data-ingester:8085/ingest >> /var/log/cron.log 2>&1

echo "Ingestion finished at $(date -u)" >> /var/log/cron.log
