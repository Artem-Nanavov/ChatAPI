import os


EMAIL_HOST = os.environ.get("EMAIL_HOST", "somerandomhost")
EMAIL_PORT = int(os.environ.get("EMAIL_PORT", 1))
EMAIL_HOST_USER = os.environ.get("EMAIL_HOST_USER", "somerandomhost")
EMAIL_HOST_PASSWORD = os.environ.get("EMAIL_HOST_PASSWORD", "somerandompassword")
EMAIL_USE_TLS = os.environ.get("EMAIL_USE_TLS", "True") == "True"

SERVER_EMAIL = EMAIL_HOST_USER
DEFAULT_FROM_EMAIL = EMAIL_HOST_USER
