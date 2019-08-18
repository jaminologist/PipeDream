extends Control

var score = 0
var time_limit = 0

var client = WebSocketClient.new()

var centerMath:CenterMath = load("res://math/center_math.gd").new()
var client_json_reader:ClientJsonReader = load("res://scenes/client_json_reader.gd").new()

onready var time_display:TimeLabel = $VBoxContainer/VBoxScoreTimeContainer/VBoxTimeContainer/Time_Counter

func _ready():
    get_node("VictoryCenterContainer").hide()
    client.connect_to_url(Connections.SINGLE_PLAYER_WEBSOCKET_STRING)
    client.connect("connection_failed", self, "_on_connection_error")  
    setup_client_json_reader()
    pass 
    
func setup_client_json_reader():
    client_json_reader.time_label = time_display
    client_json_reader.grid = $Grid
    client_json_reader.player_score_label = $VBoxContainer/VBoxScoreTimeContainer/VBoxScoreContainer/Score_Number_Label
    
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
        client_json_reader.use_json_from_server_for_grid(json, $Grid, rect_size)
        client_json_reader.use_json_from_server(json, rect_size)
        if json.get("IsOver", false):
            open_score_screen()
        

func update_time_counter_text(time_limit):
    time_display.convert_time_to_label_text_and_set_text(time_limit)
    
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
