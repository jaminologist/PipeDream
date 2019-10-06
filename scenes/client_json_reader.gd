extends Object
class_name ClientJsonReader

var grid:Grid
var time_label:TimeLabel
var player_score_label:ScoreLabel

# Called when the node enters the scene tree for the first time.
func _ready():
    pass 
    
func use_json_from_server_for_grid(response:BlitzGameResponse, grid:Grid):
    var boardReports = response.get_board_reports()
    if boardReports != null && boardReports.size() > 0:
        grid.load_boardreports_into_grid(boardReports)
        
    
func use_json_from_server(response:BlitzGameResponse):
    
    if response.get_time_limit() != null:
        time_label.convert_time_to_label_text_and_set_text(response.get_time_limit().Time)
        
    var score = response.get_score()
    if score != null && score != 0:
        player_score_label.set_score(score)

# Called every frame. 'delta' is the elapsed time since the previous frame.
#func _process(delta):
#    pass
