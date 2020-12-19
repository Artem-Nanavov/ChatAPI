from rest_framework import generics
from .serializers import *
from rest_framework.permissions import IsAuthenticated


class MessageCreationAPIView(generics.ListCreateAPIView):
    queryset = Message.objects.all()
    serializer_class = MessageSerializer
    permission_classes = [IsAuthenticated, ]
