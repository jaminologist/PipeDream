extends Control

func _ready():
    pass 

func _on_BlitzButton_pressed():
    get_tree().change_scene("res://Scenes/Play.tscn")


func _on_ClassicButton_pressed():
    get_tree().change_scene("res://Scenes/Classic.tscn")


func _on_VersusButton_pressed():
    get_tree().change_scene("res://Scenes/Versus.tscn")
