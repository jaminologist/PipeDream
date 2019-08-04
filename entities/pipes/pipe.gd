extends Node2D

var pipe_end_texture = preload("res://entities/pipes/assets/pipe_end_0.png")
var pipe_end_explosion_2_texture = preload("res://entities/pipes/assets/pipe_explosion_0.png")
var pipe_end_explosion_3_texture = preload("res://entities/pipes/assets/pipe_explosion_big_0.png")
var pipe_l_texture = preload("res://entities/pipes/assets/pipe_l_0.png")
var pipe_line_texture = preload("res://entities/pipes/assets/pipe_line_0.png")


enum Direction {UP = 0, DOWN = 180, RIGHT = 90, LEFT = 270}
enum PipeColor {Color_0 = 0, Color_1 = 1, Color_2 = 2, Color_3 = 3}


export (int) var speed

export (Color) var pipe_color_0
export (Color) var pipe_color_1
export (Color) var pipe_color_2
export (Color) var pipe_color_3

signal pipe_moving
signal pipe_stop

var startPosition = Vector2()
var destination = Vector2()
var velocity = Vector2()
var direction = Direction.UP
var type = PipeType.NONE
var is_moving = false

var travel_time

func _ready():
    direction = Direction.UP
    
func init(type :int):
    self.type = type
    set_texture_using_type(type)
    
func get_type():
    return type
    
func get_direction():
    return direction

#Assume travel_time is passed down in nanoseconds
func move_to(startPosition:Vector2, destination: Vector2, travel_time:float):
    is_moving = true
    self.travel_time = travel_time / 1000000000
    emit_signal("pipe_moving")
    
    self.startPosition = startPosition
    self.destination = destination 
    
func set_texture_using_type(type: int):
    var texture 
    
    match type:
        PipeType.END: 
            texture = pipe_end_texture
        PipeType.END_EXPLOSION_2: 
            texture = pipe_end_explosion_2_texture
        PipeType.END_EXPLOSION_3:
            texture = pipe_end_explosion_3_texture
        PipeType.LINE:
            texture = pipe_line_texture
        PipeType.L_PIPE: 
            texture = pipe_l_texture
    
    self.type = type
    get_node("Sprite").set_texture(texture)
    
func set_pipeColor(pipeColor: int):
    match pipeColor:
        PipeColor.Color_0: 
             get_node("Sprite").modulate = pipe_color_0
        PipeColor.Color_1: 
             get_node("Sprite").modulate = pipe_color_1
        PipeColor.Color_2: 
             get_node("Sprite").modulate = pipe_color_2
        PipeColor.Color_3: 
             get_node("Sprite").modulate = pipe_color_3
    
    
func set_direction(direction):   
    self.direction = direction
    get_node("Sprite").rotation_degrees = direction
    
    
#Sets the size of the texture inside the sprite to the specified width and height
func set_size(width: float, height: float):
    var th = height 
    var tw = width
    
    var currentSize = get_node("Sprite").texture.get_size()
    
    var currentScale = self.scale
    var newScale =  Vector2(currentScale.x * (tw / currentSize.x) , currentScale.y * (th / currentSize.y))
    self.scale = newScale
    #$Sprite.scale = self.scale 

    
func _physics_process(delta):
    
    if is_moving:
        velocity = Vector2(0, (destination.y - startPosition.y) * (delta / travel_time))
        position += velocity 
        
        if (startPosition - position).length() > (startPosition - destination).length():
            is_moving = false
            emit_signal("pipe_stop")
            position = destination
            velocity = Vector2()
            
#Rotates the Pipe in a clockwise direciton 90 degrees
func rotate_pipe():
    match direction:
        Direction.UP:
            set_direction(Direction.RIGHT)
        Direction.RIGHT:
            set_direction(Direction.DOWN)
        Direction.DOWN:
            set_direction(Direction.LEFT)
        Direction.LEFT:
            set_direction(Direction.UP)
        