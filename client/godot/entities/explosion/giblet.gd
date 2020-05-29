extends Node2D


export (float) var speed
export (float) var angleInDegrees
export (float) var fade
export (float) var expiryTime

var velocity

func _ready():
    velocity = Vector2(speed * cos(deg2rad(angleInDegrees)), speed * sin(deg2rad(angleInDegrees)))
    $Tween.interpolate_property($GibletSprite, "modulate", Color(1, 1, 1, 1), Color(1, 1, 1, 0), fade, 
    Tween.TRANS_QUART, Tween.EASE_IN)
    $Tween.start()


func set_size(width: float, height: float):
    var th = height 
    var tw = width
    
    var currentScale = self.scale
    self.scale = Vector2((currentScale.x/(currentScale.x/tw))/50, (currentScale.y/(currentScale.y/th))/50)
    
    
func _physics_process(delta):
    position += velocity * delta
    expiryTime -= delta
    if expiryTime <= 0:
        self.queue_free()
    