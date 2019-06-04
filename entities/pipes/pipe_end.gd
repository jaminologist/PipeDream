extends "res://entities/pipes/Pipe.gd"
    
func points_to(column: int, row: int) -> Array:
    
    match direction:
        Direction.UP:
            return [Vector2(column, row - 1)]
        Direction.DOWN:
            return [Vector2(column, row + 1)]
        Direction.LEFT:
            return [Vector2(column - 1, row)]
        Direction.RIGHT:
            return [Vector2(column + 1, row)]
            
    return []