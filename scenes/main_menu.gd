extends Control

func _ready():
    pass 

func _on_BlitzButton_pressed():
    get_tree().change_scene("res://scenes/play.tscn")


func _on_ClassicButton_pressed():
    get_tree().change_scene("res://scenes/classic.tscn")


func _on_VersusButton_pressed():
    get_tree().change_scene("res://scenes/2player_versus_container.tscn")
