FROM apache/superset:2.1.0

# Install ClickHouse drivers for Superset connectivity
USER root
RUN pip install --no-cache-dir \
    clickhouse-connect==0.7.15 \
    psycopg2-binary  # For PostgreSQL metadata DB support

# Switch back to superset user
USER superset

# Optional: Copy custom Superset config (create superset_config.py in ./superset/pythonpath)
COPY pythonpath /app/pythonpath