from django.http import HttpResponse
from django.shortcuts import render
from django.utils.safestring import mark_safe
import json
from .engine.game import Game

def index(request):
    print("INDEX")
    return render(request, 'index.html', {})

def room(request, room_name):
    print("ROOM:", room_name)
    return render(request, 'room.html', {
        'room_name_json': mark_safe(json.dumps(room_name))
    })
