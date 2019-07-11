extends Control

enum State {CONNECTING, PLAYING, COMPLETE}

var state = State.CONNECTING

var client = WebSocketClient.new()

var is_ai = true

var score:int = 0
var enemyscore: int = 0
var time_limit = 90


func _ready():
    
    if !is_ai:
        client.connect_to_url("ws://165.22.120.163:5080/connectToServer")
        client.connect("connection_error", self, "_on_connection_error")
        get_node("Grid").set_process(false)
        get_node("Grid").set_process_input(false)
    else:
        state = State.PLAYING
    
    get_node("VictoryCenterContainer").hide()
    
    #Centers Grid
    $Grid.position.x = (rect_size.x / 2 - (($Grid.column * $Grid.cell_size) / 2))
    $Grid.position.y = (rect_size.y / 2 - (($Grid.row * $Grid.cell_size) / 2)) + $Grid.cell_size * 2
    
    pass 

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
    
    if Input.is_action_just_pressed("ui_cancel"):
        client.disconnect_from_host()
        get_tree().change_scene("res://Scenes/MainMenu.tscn")
    
    if !is_ai:
        poll_client_and_update()

            
    pass

func poll_client_and_update():
    client.poll()
    var bytes = client.get_peer(1).get_packet()
    var dictionary = {}
    var json = bytes.get_string_from_utf8();
    var a = parse_json(json)
    
    if a != null:
        a as Dictionary
        if a.has("IsOver"):
            
            if state == State.CONNECTING:
                get_node("Grid").set_process(true)
                get_node("Grid").set_process_input(true)
                get_node("VBoxContainer/TitleContainer/Label").text = ""
            
            state = State.PLAYING
            time_limit = a["Time"]
            update_time_counter_text()
        elif a.has("score"):
            update_enemy_score_text(a["score"])
            enemyscore = a["score"]
    
    
func update_time_counter_text():
    var time_limit_in_seconds = float(time_limit) / 1000000000
    var minutes = time_limit_in_seconds / 60
    var seconds = fmod(time_limit_in_seconds, 60)
    var str_elapsed = "%2d:%02d" % [minutes, seconds]
    
    get_node("VBoxContainer/VBoxContainer3/VBoxTimeContainer/Time_Counter").text = str_elapsed
    
    if time_limit_in_seconds <= 0:
        open_score_screen()
    
func update_enemy_score_text(score):
    $VBoxContainer/VBoxContainer3/VBoxEnemyContainer/Enemy_Score_Number_Label.set_score(int(score))
    
    
func open_score_screen():
    $Grid.set_process(false)
    $VictoryCenterContainer.show()
    $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryScoreLabel.text = str(score)
    
    if score >= enemyscore:
        $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryTitle.text = "Victory!"
    else:
        $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryTitle.text = "Defeat!"
        
    
func _on_Grid_pipes_destroyed(number):
    score += (1000 * number) + (250 * number) 
    get_node("VBoxContainer/VBoxContainer3/VBoxScoreContainer/Score_Number_Label").set_score(score)
    
    var dictionaryToSendToOtherPlayer = {"score": score}
    #client.get_peer(1).put_packet(JSON.print(dictionaryToSendToOtherPlayer).to_ascii())
    

func _on_BlitzTimer_timeout():
    time_limit -= 1
    if time_limit <= 0:
        $BlitzTimer.stop()
        open_score_screen()
        time_limit = 0
    update_time_counter_text()


func _on_RetryButton_pressed():
    get_tree().reload_current_scene()

func _on_MainMenuButton_pressed():
    get_tree().change_scene("res://Scenes/MainMenu.tscn")


func _on_Grid_explosive_pipe_destroyed(power, time):
    $CameraShake2D.start_camera_shake(power, 0.25)



func _on_connection_error():
    get_tree().change_scene("res://Scenes/MainMenu.tscn")
    pass