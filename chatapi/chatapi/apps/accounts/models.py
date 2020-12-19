from django.db import models
from django.contrib.auth.models import AbstractBaseUser, BaseUserManager, PermissionsMixin
from rest_framework.authtoken.models import Token


class UserManager(BaseUserManager):
    def create_user(self, email, username, password, **kwargs):
        if not email or not username:
            raise ValueError('Every field must be filled')
        email = self.normalize_email(email)
        user = self.model(email=email, username=username, **kwargs)
        user.set_password(password)
        user.save()
        return user

    def create_superuser(self, email, username, password, **kwargs):
        kwargs.setdefault('is_superuser', True)
        kwargs.setdefault('is_staff', True)
        if kwargs.get('is_superuser') is True and kwargs.get('is_staff') is True:
            return self.create_user(email=email, username=username, password=password, **kwargs)
        raise ValueError('Is not superuser')


class User(AbstractBaseUser, PermissionsMixin):
    email = models.EmailField(unique=True)
    username = models.CharField(max_length=255)
    password = models.CharField(max_length=255)

    created = models.DateTimeField(auto_now_add=True)
    last_login = models.DateTimeField(auto_now=True)
    is_active = models.BooleanField(default=True)
    is_staff = models.BooleanField(default=False)
    is_superuser = models.BooleanField(default=False)

    objects = UserManager()
    REQUIRED_FIELDS = ['username']
    USERNAME_FIELD = 'email'

    def __str__(self):
        return self.email

    def authenticate(self) -> str:
        token, _ = Token.objects.get_or_create(user=self)
        return token.key

    def logout(self):
        try:
            self.auth_token.delete()
        except:
            pass
