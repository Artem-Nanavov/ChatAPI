from django.contrib.auth import get_user_model
from django.db import models

User = get_user_model()


class Message(models.Model):
    text = models.TextField()
    owner = models.ForeignKey(User, on_delete=models.DO_NOTHING, related_name='owner')
    receiver = models.ForeignKey(User, on_delete=models.DO_NOTHING, related_name='receiver')

    def __str__(self):
        return self.id
