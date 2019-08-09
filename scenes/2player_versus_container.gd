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
	poll_client_and_update()
	pass

func poll_client_and_update():
    client.poll()


func _on_connection_established():
	print("Connection establised")
