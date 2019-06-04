extends Label

var score_displayed_to_player: int
var score: int

func set_score(score: int):
    self.score = score

func _process(delta):
    if score_displayed_to_player < score:
        var difference = (score - score_displayed_to_player) / 49
        difference += randi() % 10
        if difference < 50:
            difference = 50
            
        score_displayed_to_player += difference
    else:
        score_displayed_to_player = score
        
    self.text = str(score_displayed_to_player)
