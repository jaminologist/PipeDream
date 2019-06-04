extends Node2D

enum Direction {UP = 0, DOWN = 180, RIGHT = 90, LEFT = 270}

var direction = Direction.UP

export (int) var speed

signal pipe_moving
signal pipe_stop

var destination = Vector2()
var velocity = Vector2()
var is_moving = false

func _ready():
    direction = Direction.UP
    

func move_to(destination: Vector2):
    is_moving = true
    emit_signal("pipe_moving")
    self.destination = destination 
    
func set_direction(direction):
    self.direction = direction
    get_node("Sprite").rotation_degrees = direction
    
func randomize_direction():
    var directions = [Direction.UP, Direction.DOWN, Direction.LEFT, Direction.RIGHT]
    set_direction(directions[randi() % directions.size()])
    
func set_size(width: float, height: float):
    var th = height 
    var tw = width
    
    var currentSize = get_node("Sprite").texture.get_size()
    
    var currentScale = self.scale
    var newScale =  Vector2(currentScale.x * (tw / currentSize.x) , currentScale.y * (th / currentSize.y))
    self.scale = newScale
    #$Sprite.scale = self.scale 
    
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

#Returns which column and row this pipe points to from the give column and row
func points_to(column: int, row: int) -> Array:
    
    match direction:
        Direction.UP, Direction.DOWN:
            return [Vector2(column, row + 1), Vector2(column, row - 1)]
        Direction.RIGHT, Direction.LEFT:
            return [Vector2(column + 1, row), Vector2(column - 1, row)]
            
    return []
    
func _physics_process(delta):
    
    if is_moving:
        velocity = (destination - position).normalized() * speed * delta
        #print("Vectors!", destination, position, destination - position, velocity)
        if (destination - position).length() > 10:
            position += velocity 
        else:
            is_moving = false
            emit_signal("pipe_stop")
            position = destination
            velocity = Vector2()
        