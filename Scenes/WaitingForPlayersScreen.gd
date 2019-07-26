extends Control

var waitTimer


var dotCount:int = 0

func _ready():
    var waitTimer = Timer.new()
    waitTimer.connect("timeout", self, "increment_charge_offline")
    pass # Replace with function body.

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
    pass
