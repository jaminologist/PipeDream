extends Node2D

#Grid variables
export (int) var width
export (int) var height
export (int) var x_position
export (int) var y_position
export (float) var cell_size
export (int) var pipe_fall_speed

export (int) var minimum_connection_size
export (int) var medium_connection_size
export (int) var maximum_connection_size

export (Color) var minimum_connection_color
export (Color) var medium_connection_color
export (Color) var maximum_connection_color

signal pipe_touch
signal pipes_destroyed(number)
signal explosive_pipe_destroyed(power, time)
signal connection_found 

var pipe_l = preload("res://entities/pipes/pipe_l.tscn")
var pipe_line = preload("res://entities/pipes/Pipe.tscn")
var pine_end =  preload("res://entities/pipes/pipe_end.tscn")

var pine_end_explosion_2 =  preload("res://entities/pipes/pipe_end_explosion_2.tscn")
var pine_end_explosion_3 =  preload("res://entities/pipes/pipe_end_explosion_3.tscn")

var pipe_moving_count = 0

var possible_pieces = [
   pipe_l,
   pipe_line,
   pine_end,
   #pine_end_explosion_2
]

var corner_pieces = [
   pipe_l,
  # pipe_line,
   pine_end
]

var all_pieces = []

# Called when the node enters the scene tree for the first time.
func _ready():
    randomize()
    all_pieces = create_new_pipe_grid()
    var connections = pipe_pass(all_pieces)
    
    print(connections.size())
    
    #This code looks to make sure there are no connections when the game first begins.
    for i in range(0, connections.size()):
        var rootNode = connections[i] as PipeNode
        for node in rootNode.root_and_children():
            var chainBroken = false
            for j in range (0, 3):
                node.pipe.rotate_pipe()
                var visitedNodes = {}
                if !create_pipe_connection(visitedNodes, all_pieces, PipeNode.new(node.pipe), node.position.x, node.position.y):
                    chainBroken = true   
                    break
            if chainBroken:
                break
    
    connections = pipe_pass(all_pieces)
    
    if connections.size() > 0:
        print("Oh no, there was a conneciton, you should probably write a unit test about this")
    
    #full_grid_pipe_pass_sequence(all_pieces)

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta):
    on_mouse_click()
    pass
    
#Creates a new pipe and adds it into the grid at position x and y
func generate_new_pipe(x: int, y:int, pipe_grid):
    
    var rand
    var pipe
    
    if x == 0 || x == pipe_grid.size() - 1:
        rand = floor(rand_range(0, corner_pieces.size()))
        pipe = corner_pieces[rand].instance()
    else:
        rand = floor(rand_range(0, possible_pieces.size()))
        pipe = possible_pieces[rand].instance()
        
    add_child(pipe)
    add_pipe_to_grid(x, y, pipe, pipe_grid)
    return pipe
    
#Adds pipes to grid and connects the pipe singals to the 'Grid' nodes
func add_pipe_to_grid(x: int, y:int, pipe, pipe_grid):
    pipe.position = grid_to_pixel(x, y)
    pipe.set_size(cell_size,cell_size)
    pipe.randomize_direction()
    pipe.connect("pipe_moving", self, "_on_pipe_moving")
    pipe.connect("pipe_stop", self, "_on_pipe_stop")
    pipe.speed = pipe_fall_speed
    return pipe
    
func create_new_pipe_grid():
    var twoDArray = make_2d_array()
    
    for i in width:
        for j in height:
            twoDArray[i][j] = generate_new_pipe(i, j, twoDArray)
            
    return twoDArray
            

func make_2d_array(): 
    var array = []
    for i in width:
        array.append([])
        for j in height:
            array[i].append(null)
    return array
        

#Converts the given column and row into x and y pixel co-ordinates (the top left of a cell space)
func grid_to_pixel(column, row):
    var new_x = x_position + (column * cell_size) 
    var new_y = y_position + (row * cell_size)
    return Vector2(new_x, new_y)

#Converts the given x and y into positions on the grid   
func pixel_to_grid(x, y):
    var new_x = floor((x - x_position) / cell_size)
    var new_y = floor((y - y_position) / cell_size)
    return Vector2(new_x, new_y)

#Checks if the given column and row is contained in the grid
func contains(column, row, grid):
    if column < 0 || column > grid.size() -1:
        return false
    elif row < 0 || row > grid[column].size() -1:
        return false
    return true
    
func full_grid_pipe_pass_sequence(pipe_grid):
    var connections = pipe_pass(pipe_grid)
    remove_closed_connections_from_grid(connections, pipe_grid)
    add_new_pipes_based_on_closed_connections(connections, pipe_grid)
    emit_number_of_destroyed_pipes(pipe_grid)
    update_pipe_position_after_remove(pipe_grid)
    add_new_pipes(pipe_grid)
    return connections.size() > 0
        
func on_mouse_click():
    if Input.is_action_just_pressed("ui_touch"):
        var mouse_local_position = get_local_mouse_position()
        var mouse_grid_position = pixel_to_grid(mouse_local_position.x, mouse_local_position.y)
        
        if contains(mouse_grid_position.x, mouse_grid_position.y, all_pieces):   
            all_pieces[mouse_grid_position.x][mouse_grid_position.y].rotate_pipe()
            emit_signal("pipe_touch")
            if full_grid_pipe_pass_sequence(all_pieces):
                emit_signal("connection_found")
 

func point_is_empty(x: int, y: int, pipe_grid):
    return pipe_grid[x][y] == null
      
func add_new_pipes(pipe_grid: Array):
    for x in range(0, pipe_grid.size()):
        var depth := 1
        for y in range(pipe_grid[x].size() - 1, -1, -1):
            if point_is_empty(x, y, pipe_grid):
                var new_pipe = generate_new_pipe(x, y, pipe_grid)
                pipe_grid[x][y] = new_pipe
                new_pipe.position = grid_to_pixel(x, 0 - depth)
                new_pipe.move_to(grid_to_pixel(x, y))
                depth += 1
                
func emit_number_of_destroyed_pipes(pipe_grid: Array):
    
    var count = 0
    
    for x in range(0, pipe_grid.size()):
        for y in range(0, pipe_grid.size()):
            if point_is_empty(x, y, pipe_grid):
                count += 1
                
    if count > 0:
        emit_signal("pipes_destroyed", count)

#Loops over a pipe grid and removes any closed connections that have been found
#Also, adds in new 'explosive' pipes based on the size of the connections found
func pipe_pass(pipe_grid: Array):

    var hasConnection := false
    
    var visitedNodes = {}
    var closedConnections = []
    
    for i in range(0, pipe_grid.size()):
        for j in range(0, pipe_grid[i].size()):
            if !visitedNodes.has(Vector2(i,j)) && pipe_grid[i][j] != null:
                visitedNodes[Vector2(i,j)] = true
                var rootNode = PipeNode.new(pipe_grid[i][j])
                var closedConnection = create_pipe_connection(visitedNodes, pipe_grid, rootNode, i, j)
                
                if closedConnection:
                    closedConnections.append(rootNode)
                
                var connectionSize = rootNode.size()
                for node in rootNode.root_and_children():
                    if closedConnection:
                        hasConnection = true     
                    else:
                        node.pipe.get_node("Sprite").modulate = get_color_of_pipe_based_on_connection_size(connectionSize)
    return closedConnections
    
func add_new_pipes_based_on_closed_connections(closedConnections: Array, pipeGrid):
    
    var newly_created_pipes = []
    var newly_created_pipes_positions = []
    
    print(closedConnections.size())
    
    for i in range(0, closedConnections.size()):
        var rootNode = closedConnections[0] as PipeNode
        var connectionSize = rootNode.size()
        if connectionSize >= maximum_connection_size:
            newly_created_pipes.append(pine_end_explosion_3.instance())
            newly_created_pipes_positions.append(rootNode.position)
        elif connectionSize >= medium_connection_size:
            newly_created_pipes.append(pine_end_explosion_2.instance())
            newly_created_pipes_positions.append(rootNode.position)
        
    for i in range(0, newly_created_pipes.size()):
        var pos = newly_created_pipes_positions[i]
        add_child(newly_created_pipes[i])
        add_pipe_to_grid(pos.x, pos.y, newly_created_pipes[i], pipeGrid)
        pipeGrid[pos.x][pos.y] = newly_created_pipes[i]
    
    #                        remove(node.position.x, node.position.y, pipe_grid)      
func remove_closed_connections_from_grid(closedConnections: Array, pipeGrid):
    for i in range(0, closedConnections.size()):
        var rootNode = closedConnections[0] as PipeNode
        for node in rootNode.root_and_children():
            remove(node.position.x, node.position.y, pipeGrid)  
    
                        

func get_color_of_pipe_based_on_connection_size(connectionSize: int):  
    if connectionSize < minimum_connection_size:
        return Color.white
    elif connectionSize < medium_connection_size:
        return minimum_connection_color
    elif connectionSize < maximum_connection_size:
        return medium_connection_color
    elif connectionSize >= maximum_connection_size:
        return maximum_connection_color

func positionsInSquareRange(position: Vector2, size: int) -> Array:
    
    var positions = []
    
    if size <= 0:
        return []
        
    for x in range(0, size):
        for y in range(0, size):
            if x == 0 && y == 0:
                continue
            positions.append(Vector2(position.x + x, position.y + y))
            positions.append(Vector2(position.x + x, position.y - y))
            positions.append(Vector2(position.x - x, position.y + y))
            positions.append(Vector2(position.x - x, position.y - y))
    
    return positions

#Removes a pipe from the given x and y. Uses a 'GibletFactory' to create an explosion.
func remove(x: int, y: int, grid):
    
    if !contains(x,y, grid):
        return
    
    if grid[x][y] == null:
        return
    
    var pipe = grid[x][y]
    pipe.queue_free()
    grid[x][y] = null
    var pos = grid_to_pixel(x, y)
    pos.x += cell_size / 2
    pos.y += cell_size / 2
    
    if pipe.get_filename() == pine_end_explosion_3.get_path():
        $GibletFactory.numberOfGiblets = 24
        $GibletFactory.create_explosion(pos.x, pos.y)
        emit_signal("explosive_pipe_destroyed", 6, 2)
        for position in positionsInSquareRange(pixel_to_grid(pipe.position.x, pipe.position.y), 3):
            remove(position.x, position.y, grid)
    elif pipe.get_filename() == pine_end_explosion_2.get_path():
        $GibletFactory.numberOfGiblets = 12
        $GibletFactory.create_explosion(pos.x, pos.y)
        emit_signal("explosive_pipe_destroyed", 6, 2)
        for position in positionsInSquareRange(pixel_to_grid(pipe.position.x, pipe.position.y), 3):
            remove(position.x, position.y, grid)
    else:
        $GibletFactory.numberOfGiblets = 6
        $GibletFactory.create_explosion(pos.x, pos.y)
        
    

func is_empty(x: int, y: int, grid):
    return grid[x][y] == null
    

func update_pipe_position_after_remove(pipe_grid):
    
    for x in range(0, pipe_grid.size()):
        
        var resetPosition = pipe_grid[x].size()-1
        
        var y = pipe_grid[x].size()-1
        while(y >= 0):   
            if pipe_grid[x][y] != null:
                
                if y < resetPosition:
                    
                    var pipe = pipe_grid[x][y]
                    pipe_grid[x][y] = null
                    pipe_grid[x][resetPosition] = pipe
                    pipe.move_to(grid_to_pixel(x, resetPosition))
                    
                    y = resetPosition + 1
                     
                else:
                    resetPosition -= 1
            y -= 1
            
            
func _on_pipe_moving():
    self.set_process(false)  
    pipe_moving_count += 1
    
func _on_pipe_stop():
    pipe_moving_count -= 1
    if pipe_moving_count <= 0:
        pipe_moving_count = 0
        var hasConnection = full_grid_pipe_pass_sequence(all_pieces)
        if !hasConnection :
            self.set_process(true)  
        
    
    
                
                
class PipeNode:
    var parent = null
    var children = []
    var pipe
    var position
    
    func _init(pipe):
        self.pipe = pipe
    
    func addChild(pipe):
        pipe.parent = self
        children.append(pipe)
        
    func root() -> PipeNode:
        if self.parent == null:
            return self
        else:
            return self.parent.root()
            
    func size() -> int:
        var count = 1
        
        for i in children:
            count += i.size()
            
        return count
        
    func root_and_children() -> Array:
        
        var nodeArray = []
        nodeArray.append(self)
        for i in children:
            nodeArray = nodeArray + i.root_and_children()
            
        return nodeArray
    

func create_pipe_connection(visitedNodes: Dictionary, pipe_grid: Array, pipeNode: PipeNode, column : int, row : int) -> bool:

    pipeNode.position = Vector2(column, row)
    var new_points = pipeNode.pipe.points_to(column, row)
    
    var is_closed_connection = true
    
    for i in range(0, new_points.size()):
        var x = new_points[i].x
        var y = new_points[i].y
        if contains(x, y, pipe_grid) && pipe_grid[x][y] != null:            
            var child_pipe = pipe_grid[x][y]
            var child_new_points = child_pipe.points_to(x, y)
            
            if child_new_points.has(pipeNode.position):
                
                if  !visitedNodes.has(Vector2(x,y)):
                    var child = PipeNode.new(child_pipe)
                    visitedNodes[Vector2(x,y)] = true
                    pipeNode.addChild(child)
                    var child_connection = create_pipe_connection(visitedNodes, pipe_grid, child, x, y)
                    
                    if is_closed_connection == true:
                        is_closed_connection = child_connection
                
            else:
                is_closed_connection = false
        else:
            is_closed_connection = false

    return is_closed_connection
