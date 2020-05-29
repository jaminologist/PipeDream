extends Node

class_name CenterMath

func _ready():
    pass

#Returns the x and y offset required to center a rectangle inside another rectangle.
#Add the returned values to the position of rectangle 1 to get the positions for rectangle 2 to be centered
func center_rectangle_position_offset(rect1width, rect1height, rect2width, rect2height)->Vector2:
    return Vector2((rect1width / 2 - rect2width / 2), (rect1height / 2 - rect2height / 2))