from django.http import HttpResponse
from engine.game import Game

def index(request):
    return HttpResponse("Future site of the PocketRockets poker app")
