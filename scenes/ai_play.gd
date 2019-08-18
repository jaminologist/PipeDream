extends "res://scenes/play.gd"


func _ready():
    get_node("VictoryCenterContainer").hide()
    client.connect_to_url(Connections.SINGLE_PLAYER_WEBSOCKET_STRING)
    client.connect("connection_failed", self, "_on_connection_error")  
    setup_client_json_reader()
    $Grid.set_touchable(false)
    pass 
    
func _on_Grid_pipe_touch(x:int, y:int):
    pass