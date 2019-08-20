extends Control

var score = 0
var opponent_score = 0

var client = WebSocketClient.new()
var centerMath:CenterMath = load("res://math/center_math.gd").new()
var client_json_reader:ClientJsonReader = load("res://scenes/client_json_reader.gd").new()

onready var time_display:TimeLabel = $VBoxContainer/VBoxScoreTimeContainer/VBoxTimeContainer/HBoxContainer/VBoxContainer2/Time_Counter

func _ready():
    get_node("VictoryCenterContainer").hide()
    
    client_json_reader.time_label = time_display
    client_json_reader.grid = $Grid
    client_json_reader.player_score_label = $VBoxContainer/VBoxScoreTimeContainer/VBoxTimeContainer/HBoxContainer/VBoxContainer2/Score_Number_Label
    #Make opponenet explosive Giblets a bit smaller
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid/GibletFactory.width = 1.5
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid/GibletFactory.height = 1.5
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid/GibletFactory.maxspeed = 250
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid/GibletFactory.expiryTime = 2.5
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid/GibletFactory.maxFadeTime = 1.0
    pass 
    
func _process(delta):
    poll_client_and_update()
    if Input.is_action_just_pressed("ui_cancel"):
        get_tree().change_scene("res://scenes/main_menu.tscn")
        
func set_client(client: WebSocketClient):
    self.client = client
    client.connect("connection_failed", self, "_on_connection_error")
    
func disable():
    set_process(false)
    $Grid.set_process(false)
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid.set_process(false)
    
func enable():
    set_process(true)
    $Grid.set_process(true)
    $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid.set_process(true)
        
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
            
        if json.get("EnemyInformation", null) != null:
            var enemyJson = json.get("EnemyInformation")
            
            enemyJson as Dictionary
            var container = $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer
            var rivalGrid = $VBoxContainer/VBoxScoreTimeContainer/VRivalGridContainer/RivalGrid
            
            client_json_reader.use_json_from_server_for_grid(enemyJson, rivalGrid, container.rect_size)
            if enemyJson.get("Score", null) != null:
                set_enemy_score(enemyJson.get("Score"))

func update_time_counter_text(time_limit):
    time_display.convert_time_to_label_text_and_set_text(time_limit)
    
func open_score_screen():
    $Grid.set_process(false)
    $VictoryCenterContainer.show()
    $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryScoreLabel.text = str(score)
    
func set_score(score: int):
    $VBoxContainer/VBoxScoreTimeContainer/VBoxTimeContainer/HBoxContainer/VBoxContainer2/Score_Number_Label.set_score(score)
    self.score = str(score)
    
func set_enemy_score(score: int):
    $VBoxContainer/VBoxScoreTimeContainer/VBoxTimeContainer/HBoxContainer/VBoxContainer2/Rival_Score_Number_Label.set_score(score)
    self.opponent_score = str(score)

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
    client.get_peer(1).put_packet(JSON.print(inputDictionary).to_ascii())
