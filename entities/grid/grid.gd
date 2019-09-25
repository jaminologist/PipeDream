extends Node2D
class_name Grid

#Grid variables
export (int) var x_position
export (int) var y_position
export (float) var cell_size
export (int) var pipe_fall_speed

signal pipe_touch
signal pipes_destroyed(number)
signal explosive_pipe_destroyed(power, time)

var pipe_preload = preload("res://entities/pipes/pipe.tscn")
var board
var boardReports:Array = []
var boardAnimationInProgress = false

var column
var row
var pipe_moving_count:int = 0
var isTouchable:bool = true

var size:Vector2

# Called when the node enters the scene tree for the first time.
func _ready():
    randomize()
    
# Called every frame. 'delta' is the elapsed time since the previous frame.
# Runs through the animations of each 'BoardReport' if one is availiable
func _process(_delta):
   
    if boardReports.size() > 0 && !boardAnimationInProgress:
        var boardReport = boardReports[0]
        
        if boardReport.get("DestroyedPipes", null) == null && boardReport.get("PipeMovementAnimations", null) == null:
            load_board_into_grid(boardReport.get("Board"))
            boardReports.pop_front()
            return
        if boardReport.get("DestroyedPipes", null) != null:
            load_destroyed_pipes(boardReport.get("DestroyedPipes", []))
            boardReport["DestroyedPipes"] =  null
            
        if boardReport.get("PipeMovementAnimations", null) != null:
            load_board_into_grid(boardReport.get("Board")) #Note a load the board here so the animations is correct
            #But there may not be a need to load the board again on the top. 
            load_pipe_movement_animation(boardReport.get("PipeMovementAnimations", []))
            boardReport["PipeMovementAnimations"] =  null
    else:
        boardReports = []
    
    if Input.is_action_just_pressed("ui_touch") && self.isTouchable:
        on_mouse_click() 
    
func set_touchable(isTouchable:bool):
    self.isTouchable = isTouchable

#Load the board report information
func load_boardreports_into_grid(boardReports: Array):
    if boardReports.size() <= 0:
        return
    self.boardReports = self.boardReports + boardReports

#Loads the board information and updates the grid GUI to display it
#Creates a new board if the a board goes not exist
func load_board_into_grid(board:Dictionary):
    var cells:Array = board.get("Cells")
    
    if self.board == null:
        self.column = board.get("NumberOfColumns", 0)
        self.row = board.get("NumberOfRows", 0)
        self.board = make_2d_array(self.column, self.row)
    
        self.size = Vector2(self.column * cell_size, self.row * cell_size)
        
        for x in column:
            for y in row:
                var pipe = pipe_preload.instance()
                add_child(pipe)
                pipe.connect("pipe_moving", self, "_on_pipe_moving")
                pipe.connect("pipe_stop", self, "_on_pipe_stop")
                pipe.position = grid_to_pixel(x, y)
                pipe.set_size(cell_size,cell_size)
                self.board[x][y] = pipe
                    
    for x in range(0, cells.size()):
        for y in range(0, cells[x].size()):
            var cell:Dictionary = cells[x][y]
            var pipe = self.board[x][y]
            pipe.set_texture_using_type(cell.get("Type", 2))
            pipe.set_direction(cell.get("Direction", 0))
            pipe.set_pipeColor(cell.get("Level", 0))

#Loads the destroyed pipes information. Displays an 'explosion' at the given x and y positions on the grid.
#Sends an explosive signal for different pipe types
func load_destroyed_pipes(destroyedPipes:Array):
     
    for i in range(0, destroyedPipes.size()):
        var type = destroyedPipes[i].get("Type", PipeType.LINE)
        
        var gridX = destroyedPipes[i].get("X", 0)
        var gridY = destroyedPipes[i].get("Y", 0)
        
        var pos = grid_to_pixel(gridX, gridY)
        var x = pos.x + cell_size / 2
        var y = pos.y + cell_size / 2
    
        if type == PipeType.END_EXPLOSION_3:
            $GibletFactory.numberOfGiblets = 24
            $GibletFactory.create_explosion(x, y)
            emit_signal("explosive_pipe_destroyed", 18, 3)
        elif  type == PipeType.END_EXPLOSION_2:
            $GibletFactory.numberOfGiblets = 12
            $GibletFactory.create_explosion(x, y)
            emit_signal("explosive_pipe_destroyed", 6, 2)
        else:
            $GibletFactory.numberOfGiblets = 6
            $GibletFactory.create_explosion(x, y)

#Loads the pipe movement animation instructions. The pipes must move from the given Start X and Y to the End X and Y
#in the given time frame
func load_pipe_movement_animation(pipeMovementAnimations:Array):
    
    if pipeMovementAnimations.size() > 0:
        boardAnimationInProgress = true
    
    for i in range(0, pipeMovementAnimations.size()):
        var pipeMovementAnimation = pipeMovementAnimations[i]
        var startX = pipeMovementAnimation.get("X", 0)
        var startY = pipeMovementAnimation.get("StartY", 0)
        var endY = pipeMovementAnimation.get("EndY", 0)
        var travel_time = pipeMovementAnimation.get("TravelTime", 0)
        var pipe = self.board[startX][endY]
        pipe.position = grid_to_pixel(startX, startY)
        pipe.move_to(grid_to_pixel(startX, startY), grid_to_pixel(startX, endY), travel_time)  
    pass

#Increases the count of pipes moving whenever a pipe is moved
#When this method is trigged input processing is turned off
func _on_pipe_moving():
    self.set_process(false)  
    pipe_moving_count += 1

#Decreases the count pipes moving whenever a pipe stops
#Once the count is back to zero, input processing is turned back on
func _on_pipe_stop():
    pipe_moving_count -= 1
    if pipe_moving_count <= 0:
        pipe_moving_count = 0
        boardAnimationInProgress = false
        self.set_process(true)  

func make_2d_array(column: int, row: int): 
    var array = []
    for i in column:
        array.append([])
        for j in row:
            array[i].append(null)
    return array
        

#Converts the given column and row into x and y pixel co-ordinates (the top left of a cell space)
#As the y co-ordinate is reversed when drawing. This method reverts so the pixel is correct
func grid_to_pixel(column, row):
    var new_x = x_position + (column * cell_size) 
    var new_y = y_position + ((self.row - 1) * cell_size) - (row * cell_size)
    return Vector2(new_x, new_y)

#Converts the given x and y into positions on the grid.
#As the y co-ordinate is reversed when drawing. This method reverts so the grid is correct
func pixel_to_grid(x, y):
    var new_x = floor((x - x_position) / cell_size)
    var new_y = (self.row - 1) - floor((y - y_position) / cell_size)
    return Vector2(new_x, new_y)
  
#Whenever the grid is touched it sends a signal of the x and y position of the grid where the touch occured   
func on_mouse_click():
        var mouse_local_position = get_local_mouse_position()
        var mouse_grid_position = pixel_to_grid(mouse_local_position.x, mouse_local_position.y)
        
        var gridX = mouse_grid_position.x
        var gridY = mouse_grid_position.y
        
        if board != null:
            if contains(gridX, gridY, board):   
                self.board[gridX][gridY].rotate_pipe()
                emit_signal("pipe_touch", gridX, gridY)

func get_new_pipe_instance(pipeType: int):
    var pipe = pipe_preload.instance()
    pipe.init(pipeType)
    return pipe

func contains(column, row, grid):
    if column < 0 || column > grid.size() -1:
        return false
    elif row < 0 || row > grid[column].size() -1:
        return false
    return true
    
                    
        