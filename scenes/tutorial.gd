extends "res://scenes/play.gd"



func setup():
    get_node("VictoryCenterContainer").hide()
    client.connect_to_url(Connections.TUTORIAL_WEBSOCKET_STRING)
    client.connect("connection_failed", self, "_on_connection_error")  
    setup_client_json_reader()
    $Grid.set_touchable(true)
    grid = $Grid

func _ready():
    pass