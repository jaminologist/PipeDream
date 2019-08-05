extends Control

enum State {CONNECTING, PLAYING, COMPLETE}

var state = State.CONNECTING

var client = WebSocketClient.new()

var is_ai = true

var score:int = 0
var enemyscore: int = 0
var time_limit = 90

var charge:int = 0
var charge_increment = 1
var charge_increment_count = 0
var charge_increment_max_count = 25

var charge_increment_timer


var enemy_board_list = []
var selected_enemy_board

class ChargeHandler:
    var charge:int = 0
    var charge_increment = 1
    var charge_increment_count = 0
    var charge_increment_max_count = 25




func _ready():
    
    if !is_ai:
        client.connect_to_url("ws://165.22.120.163:5080/connectToServer")
        client.connect("connection_error", self, "_on_connection_error")
        get_node("Grid").set_process(false)
        get_node("Grid").set_process_input(false)
    else:
        state = State.PLAYING
        
        var charge_increment_timer = Timer.new()
        charge_increment_timer.connect("timeout", self, "increment_charge_offline")
        charge_increment_timer.set_wait_time(1)
        charge_increment_timer.start()
        self.add_child(charge_increment_timer)
    
    get_node("VictoryCenterContainer").hide()
    
    #Centers Grid
    $Grid.position.x = (rect_size.x / 2 - (($Grid.column * $Grid.cell_size) / 2))
    $Grid.position.y = (rect_size.y / 2 - (($Grid.row * $Grid.cell_size) / 2)) + $Grid.cell_size * 2
    
    enemy_board_list = [
        $VBoxContainer/EnemyBoardContainer/EnemyBoardCenterContainer,
        $VBoxContainer/EnemyBoardContainer/EnemyBoardCenterContainer2,
        $VBoxContainer/EnemyBoardContainer/EnemyBoardCenterContainer3,
        $VBoxContainer/EnemyBoardContainer/EnemyBoardCenterContainer4
    ]
    
    for i in range(0, enemy_board_list.size()):
        enemy_board_list[i].connect("enemy_board_clicked", self, "_on_enemy_board_clicked")
        enemy_board_list[i].set_id(i)
        

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
            
#func update():
    
    
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
    
    
func increment_charge_offline():
    
    charge += charge_increment
    charge_increment_count+=1
    
    if charge_increment_count >= charge_increment_max_count:
        charge_increment+=1
        charge_increment_count = 0
    
    if charge >= 100:
        charge = 100
        open_score_screen()
        
    $VBoxContainer/VBoxContainer3/VBoxScoreContainer/Charge_Number_Label.text = str(charge) + "%"
        
        
    
    
    
    
func open_score_screen():
    $Grid.set_process(false)
    $VictoryCenterContainer.show()
    $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryScoreLabel.text = str(score)
    
    if score >= enemyscore:
        $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryTitle.text = "Victory!"
    else:
        $VictoryCenterContainer/PanelContainer/VBoxContainer/VictoryTitle.text = "Defeat!"
        
    
func _on_Grid_pipes_destroyed(number):
    
    charge -= number
    
    if charge < 0:
        charge = 0
        
    $VBoxContainer/VBoxContainer3/VBoxScoreContainer/Charge_Number_Label.text = str(charge) + "%"
    
    #score += (1000 * number) + (250 * number) 
    #get_node("VBoxContainer/VBoxContainer3/VBoxScoreContainer/Score_Number_Label").set_score(score)
    
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


func _on_enemy_board_clicked(enemy_board):
    
    if selected_enemy_board != null:
        selected_enemy_board.unselect()
    
    selected_enemy_board = enemy_board
    
    enemy_board.select()
    
    pass


func _on_connection_error():
    get_tree().change_scene("res://Scenes/MainMenu.tscn")
    pass