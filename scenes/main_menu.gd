extends Control

func _ready():
    pass 

func _on_BlitzButton_pressed():
    get_tree().change_scene("res://scenes/play.tscn")

func _on_VersusButton_pressed():
    get_tree().change_scene("res://scenes/2player_versus_container.tscn")


func _on_AIBlitzButton_pressed():
    get_tree().change_scene("res://scenes/ai_play.tscn")
