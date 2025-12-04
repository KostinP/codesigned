import os

# Metadata database URI (using the new superset-db service)
SQLALCHEMY_DATABASE_URI = (
    f"postgresql+psycopg2://{os.environ.get('SUPERSET_DB_USER', 'superset')}:"
    f"{os.environ.get('SUPERSET_DB_PASSWORD', 'superset_password')}@"
    f"{os.environ.get('SUPERSET_DB_HOST', 'superset-db')}:5432/"
    f"{os.environ.get('SUPERSET_DB_NAME', 'superset')}"
)

# Secret key from .env
SECRET_KEY = os.environ.get('SUPERSET_SECRET_KEY')

# Enable ClickHouse engine if not already (Superset should detect it after driver installation)
# Additional configs can be added here, e.g., for feature flags or caching