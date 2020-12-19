from rest_framework import serializers
from .models import *


class UserSerializer(serializers.ModelSerializer):
    def save(self, **kwargs):
        user = User(**self.validated_data)
        user.set_password(self.validated_data['password'])
        user.save()
        return user

    class Meta:
        model = User
        read_only_fields = [
            'id', 'is_superuser',
        ]
        exclude = [
            'user_permissions',
            'groups',
            'is_staff',
            'is_active',
        ]

