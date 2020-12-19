from django.conf import settings
from django.core.mail import send_mail
from django.core import signing


def send_confirmation_key(data):
    conf_key = signing.dumps(data)
    url = data.get("url", "localhost:8000/accounts/registration/confirmation")
    link = "http://" + url + "/?conf_key=" + conf_key
    message = 'Для подтверждения регистрации перейдите по <a href="{}">ссылке</a>'.format(link)
    send_mail(
        "Подтверждение регистрации", "",
        settings.EMAIL_HOST_USER,
        [data.get("email"), ],
        fail_silently=True,
        html_message=message,
    )
    return conf_key
