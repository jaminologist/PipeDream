FROM ioribranford/godot-docker

RUN mkdir /pipedream-godot-client

#CREATE A EXPORTS FODLER AS IN GODOT 3.1.1, IF A FODLER DOES NOT EXIST YOU CAN NOT EXPORT TO IT
RUN mkdir /pipedream-godot-client/exports

WORKDIR /pipedream-godot-client

COPY . /pipedream-godot-client 

#EXPORT HTML5 GAME USING "HTML5" TEMPLATE
RUN godot --export "HTML5" "/pipedream-godot-client/exports/index.html"