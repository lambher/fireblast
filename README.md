# Dependencies

````sudo apt install libx11-dev
sudo apt install libxcursor-dev
sudo apt install libxrandr-dev
sudo apt-get install libxinerama1 libxinerama-dev
sudo apt install libxi-dev
sudo apt install mesa-common-dev
sudo apt install libglu1-mesa-dev freeglut3-dev
sudo apt install libglfw3-dev libgles2-mesa-dev
sudo apt install pkg-config
````

# How to run the game

Start the server: ```go run server/main.go```

Start the client: ``go run client/main.go``

If you want to host the game in your  own server you can change the conf.json file like this:
````
{
  "MaxX": 800,
  "MaxY": 800,
  "Address": "[IP ADDRESS]:27017"
}

````

# Video Demo

[![Watch the video](https://img.youtube.com/vi/29tvFfrUj2w/maxresdefault.jpg)](https://youtu.be/29tvFfrUj2w)
