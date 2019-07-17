extends CenterContainer

var id:int

signal enemy_board_clicked

func _ready():
    pass
    
func _process(delta):
    if Input.is_action_just_pressed("ui_touch"):
        
        var rect = Rect2(get_global_transform_with_canvas().origin, get_rect().size)
        if rect.has_point(get_global_mouse_position()):
            emit_signal("enemy_board_clicked", self)
            print("emit!: " + str(id))
    pass 
    
func set_id(id:int):
    self.id = id
    
    
func select():
    $Selected.show() 
    
func unselect():
    $Selected.hide() 