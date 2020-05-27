extends Object
class_name BlitzGameResponse

var json:Dictionary
var BoardReports:Array
var Score:int
var IsOver:bool
var DestroyedPipes:Array
var timeLimit:TimeLimit

func _init(json):
    self.json = json
    return
    
func get_board_reports() -> Array:
    return json.get("BoardReports", []) 
    
func get_score() -> float:
    return json.get("Score", 0) 
    
func get_time_limit() -> TimeLimit:
    var timeLimit = json.get("TimeLimit") 
    if timeLimit == null:
        return null
    else:
        return TimeLimit.new(timeLimit.get("Time"))
        
func isOver() -> bool:
    return json.get("IsOver", false)


class BoardReport:
    var DestroyedPipes:Array = []
    var PipeMovementAnimations: Array = []
    var MaximumAnimationTime: float
    var IsNewBoard: bool
    var board:Board
    
    func _init(d:Dictionary):
        if d.get("DestroyedPipes") != null:
           self.DestroyedPipes = d.get("DestroyedPipes", [])
        
        if d.get("PipeMovementAnimations") != null:
           self.PipeMovementAnimations = d.get("PipeMovementAnimations", [])
        
        self.MaximumAnimationTime = d.get("MaximumAnimationTime", 0)
        self.IsNewBoard = d.get("IsNewBoard", false)
        self.board = Board.new(d.get("Board"))
        
    func get_destroyed_pipes() -> Array:
        return self.DestroyedPipes
        
    func get_pipe_movement_animations() -> Array:
        return self.PipeMovementAnimations
        
    func get_maximum_animation_time() -> float:
        return self.MaximumAnimationTime
        
    func get_board() -> Board:
        return self.board
        
    func is_new_board() -> bool:
        return self.IsNewBoard
    
class Board:
    var cells:Array
    var numberOfColumns:float
    var numberOfRows:float
    
    func _init(d:Dictionary):
        self.cells = d.get("Cells")
        self.numberOfColumns = d.get("NumberOfColumns")
        self.numberOfRows = d.get("NumberOfRows")
        return
        
class ResponsePipe:
    var type:float
    var direction:float
    var level:float
    var x:float
    var y:float
    
    func _init(d:Dictionary):
        self.type = d.get("Type")
        self.direction = d.get("Direction")
        self.level = d.get("Level")
        self.x = d.get("X")
        self.y = d.get("Y")
        pass
    
class DestroyedPipe:
    var type:float
    var x:float
    var y:float
    
    func _init(d:Dictionary):
        self.type = d.get("Type", PipeType.LINE)
        self.x = d.get("X", 0)
        self.y = d.get("Y", 0)
        pass
        
class PipeMovementAnimation:
    var x:float
    var startY:float
    var endY:float
    var travelTime:float
    
    func _init(d:Dictionary):
        self.x = d.get("X", 0)
        self.startY = d.get("StartY", 0)
        self.endY = d.get("EndY", 0)
        self.travelTime = d.get("TravelTime", 0)
        pass
        
        
class TimeLimit:
    var Time:float
    func _init(time:float):
        self.Time = time
          
    