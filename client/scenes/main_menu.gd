extends Control

func _ready():
    pass 

func _on_BlitzButton_pressed():
    get_tree().change_scene("res://scenes/play.tscn")

func _on_VersusButton_pressed():
    Connections.EDITABLE_PLAYER_WEBSOCKET_STRING = Connections.VERSUS_PLAYER_WEBSOCKET_STRING
    get_tree().change_scene("res://scenes/2player_versus_container.tscn")


func _on_AIBlitzButton_pressed():
    get_tree().change_scene("res://scenes/ai_play.tscn")

func _on_VersusAIBlitzButton_pressed():
    Connections.EDITABLE_PLAYER_WEBSOCKET_STRING = Connections.VERSUS_AI_WEBSOCKET_STRING
    get_tree().change_scene("res://scenes/2player_versus_container.tscn")


func _on_Tutorial_pressed():
    get_tree().change_scene("res://scenes/tutorial.tscn")
