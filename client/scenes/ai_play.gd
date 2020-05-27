extends "res://scenes/play.gd"


func _ready():
    pass 
    
func setup():
    get_node("VictoryCenterContainer").hide()
    client.connect_to_url(Connections.AI_SINGLE_PLAYER_WEBSOCKET_STRING)
    client.connect("connection_failed", self, "_on_connection_error")  
    setup_client_json_reader()
    $Grid.set_touchable(false)
    grid = $Grid
    
func _on_Grid_pipe_touch(x:int, y:int):
    pass