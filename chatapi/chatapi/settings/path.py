from pathlib import Path
import os
import sys


PROJECT_ROOT = os.path.join(os.path.dirname(__file__), "..")
sys.path.insert(0, os.path.join(PROJECT_ROOT, "apps"))

BASE_DIR = Path(__file__).resolve().parent.parent.parent

AUTH_USER_MODEL = 'accounts.User'

ROOT_URLCONF = 'chatapi.urls'
WSGI_APPLICATION = 'chatapi.wsgi.application'

MEDIA_ROOT = os.path.join(PROJECT_ROOT, 'media')
MEDIA_URL = '/media/'

STATIC_URL = '/static/'
STATIC_ROOT = os.path.join(PROJECT_ROOT, 'static')

