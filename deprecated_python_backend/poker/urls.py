from django.conf.urls import url

from . import views

urlpatterns = [
    url(r'^$', views.index, name='index'),
    url(r'^poker/room/(?P<room_name>[^/]+)/$', views.room, name='room')
]
