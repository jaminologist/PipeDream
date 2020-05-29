extends Control

export (int) var numberOfMoves
export (int) var target

var currentNumberOfMoves

func _ready():
    get_node("VictoryCenterContainer").hide()
    currentNumberOfMoves = numberOfMoves
    update_moves_label()
    pass
    
func open_score_screen():
    $Grid.set_process(false)
    $VictoryCenterContainer.show()
    $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryScoreLabel.text = str(0)

func _on_Grid_connection_found():
    pass # Replace with function body.

func update_moves_label():
    var score_format = "%d/%d"
    $VBoxContainer/VBoxContainer3/VBoxMovesContainer/Moves_Counter.text = score_format % [currentNumberOfMoves, numberOfMoves]

func _on_Grid_pipe_touch():
    currentNumberOfMoves -= 1
    if currentNumberOfMoves <= 0:
        open_score_screen()
    update_moves_label()


func _on_Grid_pipes_destroyed(number):
    pass 
    
func _on_Grid_explosive_pipe_destroyed(power, time):
    $CameraShake2D.start_camera_shake(power, 0.25)
    
    
func _on_RetryButton_pressed():
    get_tree().reload_current_scene()

func _on_MainMenuButton_pressed():
    get_tree().change_scene("res://scenes/main_menu.tscn")