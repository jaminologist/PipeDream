extends Node

var WEBSOCKET_HEADER = "wss://server.wthpd.com"
var SINGLE_PLAYER_WEBSOCKET_STRING = WEBSOCKET_HEADER + "/singlePlayerBlitzGame"
var VERSUS_PLAYER_WEBSOCKET_STRING = WEBSOCKET_HEADER + "/versusBlitzGame"
var AI_SINGLE_PLAYER_WEBSOCKET_STRING = WEBSOCKET_HEADER + "/aiBlitzGame"
var VERSUS_AI_WEBSOCKET_STRING = WEBSOCKET_HEADER + "/versusAiBlitzGame"
var TUTORIAL_WEBSOCKET_STRING = WEBSOCKET_HEADER + "/tutorialGame"

var EDITABLE_PLAYER_WEBSOCKET_STRING = VERSUS_PLAYER_WEBSOCKET_STRING

func _ready():
    pass
