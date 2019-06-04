extends Camera2D

var shake_power = 0
var shake_time = 0

func _ready():
    randomize()
    pass # Replace with function body.
    
func start_camera_shake(power: float, time: float):  
    if shake_power < power:
        shake_power = power
        shake_time = time

func _process(delta):
    
    if shake_time > 0:
        shake_time -= delta
        var offsetX = ((randf()*2) - 1) * shake_power
        var offsetY = ((randf()*2) - 1) * shake_power
        
        offset.x = offsetX
        offset.y = offsetY
    else:
        shake_power = 0
        offset.x = 0
        offset.y = 0
        
