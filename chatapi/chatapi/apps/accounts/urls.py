from django.urls import path
from .views import *

urlpatterns = [
    path("reg", registration),
    path("password/", change_password),
    path("login", login),
    path("logout/", logout),
    path("me/", MyProfileView.as_view()),
    path("users/", ALlUsersAPIView.as_view())
]
