# chat/routing.py
from django.conf.urls import url

from .front_end_consumer import FrontEndConsumer
from . import consumers

websocket_urlpatterns = [
    url(r'^ws/poker/room/(?P<room_name>[^/]+)/$', FrontEndConsumer),
    url(r'^ws/text/poker/room/(?P<room_name>[^/]+)/$', consumers.TextPokerConsumer),
]
