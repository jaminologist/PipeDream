extends Node

class_name ClientJsonReader

var grid:Grid
var time_label:TimeLabel
var player_score_label:ScoreLabel

# Called when the node enters the scene tree for the first time.
func _ready():
    pass 
    
func use_json_from_server_for_grid(json:Dictionary,grid:Grid):
    if json != null:
        if json.get("BoardReports", null) != null:
            var firstload = grid.board == null
            var boardReports = json.get("BoardReports", null) 
            if boardReports.size() > 0:
                grid.load_boardreports_into_grid(boardReports)
        
    
func use_json_from_server(json):
    
    if json != null:
        json as Dictionary
        
        if json.get("Time", null) != null:
            time_label.convert_time_to_label_text_and_set_text(json.get("Time", 0))
            
        if json.get("Score", null) != null:
            player_score_label.set_score(json.get("Score", 0))

# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#    pass
