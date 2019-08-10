extends Control

var client = WebSocketClient.new()

# Called when the node enters the scene tree for the first time.
func _ready():
    
    $TwoPlayerVersus.hide()
    $TwoPlayerVersus.set_process(false)
    $TwoPlayerVersus/Grid.set_process(false)
    
    print(Connections.VERSUS_PLAYER_WEBSOCKET_STRING)
    client.connect_to_url(Connections.VERSUS_PLAYER_WEBSOCKET_STRING)
    client.connect("connection_failed", self, "_on_connection_error")
    client.connect("connection_established", self, "_on_connection_established")
    pass # Replace with function body.

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
    
    if !$TwoPlayerVersus.is_processing():
        client.poll()
        
        if client.get_connection_status() == client.CONNECTION_DISCONNECTED:
            _on_connection_disconnected()
            
        var bytes = client.get_peer(1).get_packet()
        var json = parse_json(bytes.get_string_from_utf8())
        
                
        if json != null:
            json as Dictionary
            
            if json.has("IsStarted"):
                $TwoPlayerVersus.set_client(client)
                $WaitingForPlayersScreen.hide()
                $TwoPlayerVersus.show()
                $TwoPlayerVersus.set_process(true)
                $TwoPlayerVersus/Grid.set_process(true)
                $WaitingForPlayersScreen.set_process(true)
    
func _on_connection_established():
    print("Connection establised")
    
func _on_connection_error():
    print("Connection Error!")
    
func _on_connection_disconnected():
    get_tree().change_scene("res://scenes/main_menu.tscn")
