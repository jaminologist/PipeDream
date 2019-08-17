extends Control

var waitTimer


var dotCount:int = 0

func _ready():
    var waitTimer = Timer.new()
    waitTimer.connect("timeout", self, "_on_waittimer_timeout")
    waitTimer.set_wait_time(0.25)
    waitTimer.start()
    self.add_child(waitTimer)
    pass

func _process(delta):
    pass

func _on_waittimer_timeout():
    
    dotCount += 1
    
    if dotCount > 4:
        dotCount = 0
        
    var dotString = ""
    
    for i in range(dotCount):
        dotString += "."
        if i < (dotCount - 1):
            dotString += " "
            
    $VBoxContainer/Dots.set_text(dotString)
        
        