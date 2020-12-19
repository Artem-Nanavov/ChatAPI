from django.urls import path

from .views import *

urlpatterns = [
    path('messages/', MessageCreationAPIView.as_view()),
]
