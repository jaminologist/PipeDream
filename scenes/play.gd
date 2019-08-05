extends Control

var score = 0
var time_limit = 0

var client = WebSocketClient.new()

func _ready():
    get_node("VictoryCenterContainer").hide()
    #update_time_counter_text(90)
    
    client.connect_to_url("ws://165.22.120.163:5080/singlePlayerBlitzGame")
   # client.connect_to_url("ws://192.168.99.100:5080/singlePlayerBlitzGame")
    print(client.get_connection_status())
    client.connect("connection_failed", self, "_on_connection_error")
    pass 
    
func _process(delta):
    poll_client_and_update()
    if Input.is_action_just_pressed("ui_cancel"):
        get_tree().change_scene("res://scenes/main_menu.tscn")
        
        
func poll_client_and_update():
    client.poll()
    
    if client.get_connection_status() == client.CONNECTION_DISCONNECTED:
        _on_connection_error()
    
    var bytes = client.get_peer(1).get_packet()
    var json = parse_json(bytes.get_string_from_utf8())
    
    if json != null:
        json as Dictionary
        
        if json.get("IsOver", false):
            open_score_screen()
        
        if json.get("BoardReports", null) != null:
            var boardReports = json.get("BoardReports", null) 
            if boardReports.size() > 0:
                $Grid.load_boardreports_into_grid(boardReports)
        
        if json.get("Time", null) != null:
            update_time_counter_text(json.get("Time"))
            
        if json.get("Score", null) != null:
            set_score(json.get("Score", 0))
        if json.get("Board", null) != null:
            var firstload = $Grid.board == null
            $Grid.load_board_into_grid(json.get("Board"))
            
            if firstload:
                $Grid.position.x = (rect_size.x / 2 - (($Grid.column * $Grid.cell_size) / 2))
                $Grid.position.y = (rect_size.y / 2 - (($Grid.row * $Grid.cell_size) / 2)) + $Grid.cell_size * 2
                
        if json.get("DestroyedPipes", null) != null:
            $Grid.load_destroyed_pipes(json.get("DestroyedPipes", []))
        

func update_time_counter_text(time_limit):
    var time_limit_in_seconds = float(time_limit) / 1000000000
    var minutes = time_limit_in_seconds / 60
    var seconds = fmod(time_limit_in_seconds, 60)
    var str_elapsed = "%2d:%02d" % [minutes, seconds]
    
    get_node("VBoxContainer/VBoxScoreTimeContainer/VBoxTimeContainer/Time_Counter").text = str_elapsed
    
func open_score_screen():
    $Grid.set_process(false)
    $VictoryCenterContainer.show()
    $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryScoreLabel.text = str(score)
    
func set_score(score: int):
    get_node("VBoxContainer/VBoxScoreTimeContainer/VBoxScoreContainer/Score_Number_Label").set_score(score)
    self.score = str(score)

func _on_BlitzTimer_timeout():
    time_limit -= 1
    if time_limit <= 0:
        $BlitzTimer.stop()
        open_score_screen()
        time_limit = 0
    #update_time_counter_text(90)


func _on_RetryButton_pressed():
    get_tree().reload_current_scene()

func _on_MainMenuButton_pressed():
    get_tree().change_scene("res://scenes/main_menu.tscn")


func _on_Grid_explosive_pipe_destroyed(power, time):
    $CameraShake2D.start_camera_shake(power, 0.25)
    
func _on_connection_error():
    get_tree().change_scene("res://scenes/main_menu.tscn")
    pass


func _on_Grid_pipe_touch(x:int, y:int):
    var inputDictionary = {"x": x, "y": y}
    var start = OS.get_ticks_usec()
    client.get_peer(1).put_packet(JSON.print(inputDictionary).to_ascii())
    var elapsed = OS.get_ticks_usec() - start
