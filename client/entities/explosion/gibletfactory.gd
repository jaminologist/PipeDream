extends Node2D

export (float) var numberOfGiblets
export (float) var minspeed
export (float) var maxspeed

export (float) var minFadeTime
export (float) var maxFadeTime

export (float) var width 
export (float) var height 

export (float) var expiryTime

var giblet =  preload("res://entities/explosion/giblet.tscn")

func _ready():
    
    get_local_mouse_position()
    
    pass
    
func create_explosion(x, y):
    for i in range(0, numberOfGiblets):
        var g = giblet.instance()
        g.angleInDegrees = randf() * 360
        g.speed = (randf() * (maxspeed - minspeed)) + minspeed
        g.fade = (randf() * (maxFadeTime - minFadeTime)) + minFadeTime
        g.set_size(width, height)
        g.position = Vector2(x, y)
        g.expiryTime = expiryTime
        add_child(g)
        
        var g2 = giblet.instance()
        g2.angleInDegrees = g.angleInDegrees + 180
        g2.speed = -g.speed
        g2.fade = g.fade
        g2.set_size(width, height)
        g2.position = Vector2(x, y)
        g2.expiryTime = expiryTime
        
func local_test():
    if Input.is_action_just_pressed("ui_touch"):
        var pos = get_local_mouse_position()
        create_explosion(pos.x, pos.y)
    
func _process(delta):
    #local_test()
    pass
