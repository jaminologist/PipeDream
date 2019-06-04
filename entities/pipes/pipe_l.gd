extends "res://entities/pipes/Pipe.gd"

var pipe_texture = preload("res://entities/pipes/assets/pipe_l_0.png")

func _ready():
    get_node("Sprite").set_texture(pipe_texture)
    
func points_to(column: int, row: int) -> Array:
    
    match direction:
        Direction.UP:
            return [Vector2(column + 1, row), Vector2(column, row - 1)]
        Direction.DOWN:
            return [Vector2(column - 1, row), Vector2(column, row + 1)]
        Direction.LEFT:
            return [Vector2(column - 1, row), Vector2(column, row - 1)]
        Direction.RIGHT:
            return [Vector2(column + 1, row), Vector2(column, row + 1)]
            
    return []