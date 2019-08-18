extends Node

class_name ClientJsonReader

var grid:Grid
var time_label:TimeLabel
var player_score_label:ScoreLabel

var centerMath:CenterMath = load("res://math/center_math.gd").new()

# Called when the node enters the scene tree for the first time.
func _ready():
    pass 
    
func use_json_from_server_for_grid(json:Dictionary,grid:Grid, rect_size):
    if json != null:
        if json.get("BoardReports", null) != null:
            var boardReports = json.get("BoardReports", null) 
            if boardReports.size() > 0:
                grid.load_boardreports_into_grid(boardReports)
                
        if json.get("Board", null) != null:
            var firstload = grid.board == null
            grid.load_board_into_grid(json.get("Board"))
            
            if firstload:
                var pos:Vector2 = centerMath.center_rectangle_position_offset(rect_size.x, rect_size.y, grid.size.x, grid.size.y)
                grid.position.x = pos.x
                grid.position.y = pos.y + (grid.cell_size * 2)
    
func use_json_from_server(json, rect_size):
    
    if json != null:
        json as Dictionary
        
        #if json.get("IsOver", false):
            #open_score_screen()
        
        if json.get("Time", null) != null:
            time_label.convert_time_to_label_text_and_set_text(json.get("Time", 0))
            
        if json.get("Score", null) != null:
            player_score_label.set_score(json.get("Score", 0))

# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#    pass
