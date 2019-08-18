extends Label
class_name TimeLabel

func _ready():
    pass
    
    
func convert_time_to_label_text_and_set_text(time):
    var time_in_seconds = float(time) / 1000000000
    var minutes = time_in_seconds / 60
    var seconds = fmod(time_in_seconds, 60)
    var str_elapsed = "%2d:%02d" % [minutes, seconds]
    
    self.text = str_elapsed