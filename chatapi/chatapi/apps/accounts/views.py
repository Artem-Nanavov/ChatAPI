from rest_framework.response import Response
from rest_framework.decorators import api_view, permission_classes
from rest_framework import status, generics
from rest_framework.permissions import IsAuthenticated
from rest_framework.views import APIView

from .serializers import *
from .models import *


@api_view(['POST'])
def registration(request):
    serializer = UserSerializer(data=request.data)
    if serializer.is_valid():
        user = serializer.save()
        token = user.authenticate()
        return Response({"token": token, "id": user.id, "username": user.username}, status=status.HTTP_201_CREATED)
    return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)


@api_view(['POST'])
def login(request):
    email = request.data.get("email")
    try:
        user = User.objects.get(email=email)
    except User.DoesNotExist:
        return Response({"error": "there are no user with email " + email}, status=status.HTTP_400_BAD_REQUEST)

    if not user.check_password(request.data.get("password")):
        return Response({"error": "wrong email or password"}, status=status.HTTP_400_BAD_REQUEST)

    token = user.authenticate()
    return Response({"token": token, "id": user.id, "username": user.username}, status=status.HTTP_201_CREATED)


@api_view(['POST'])
@permission_classes([IsAuthenticated, ])
def logout(request):
    request.user.logout()
    return Response({
        "success": "logged out successfully",
    }, status=status.HTTP_200_OK)


class MyProfileView(APIView):
    permission_classes = [IsAuthenticated, ]

    def get(self, request):
        serializer = UserSerializer(request.user)
        return Response(serializer.data, status=status.HTTP_200_OK)

    def put(self, request):
        serializer = UserSerializer(instance=request.user, data=request.data, )
        if serializer.is_valid():
            serializer.save()
            return Response(serializer.data, status=status.HTTP_200_OK)
        return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)


@api_view(['PUT'])
@permission_classes([IsAuthenticated, ])
def change_password(request):
    if not request.user.check_password(request.data.get("password")):
        return Response({
            "error": "wrong password",
        }, status=status.HTTP_400_BAD_REQUEST)
    new_password = request.data.get("new_password", "")
    if new_password == "":
        return Response({
            "error": "new password must be provided",
        }, status=status.HTTP_400_BAD_REQUEST)
    request.user.set_password(new_password)
    request.user.save()
    return Response({
        "success": "password has changed"
    }, status=status.HTTP_200_OK)


class ALlUsersAPIView(generics.ListAPIView):
    queryset = User.objects.all()
    serializer_class = UserSerializer
    permission_classes = [IsAuthenticated, ]
