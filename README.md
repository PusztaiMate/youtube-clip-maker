# Youtube-clip creator utility

## About

The repository contains the source code for an app, that can be used to create clips from youtube videos. When started it provides a web UI, that lets the user input a youtube URL (or I guess any other URL that can be handled by the [youtube-dl](https://youtube-dl.org) utility), and zero or more clips. As the apps porpuse right now is to create small clips from our training games, the clips are currently made up of for information:

1. Start time of the clip in "mm:ss" format
2. End time of the clip in "mm:ss" format
3. Who is the subject of the clip (named player or the team)?
4. What happens (e.g.: 'nice teamplay', 'good tackle', etc.)?

## Internals

The app can be spun up by using the docker-compose file with the following command executed in the git root dir:

```bash
docker-compose up --build -d
```

This fires up two containers: the backend and the frontend.

### The frontend

The frontend is JS, using the VUE(3) framework. It is as minimalistic as one could make it with the limited knowledge I have. The docker building process is uses an build image to build (with npm) the app, and then uses an nginx base to serve it. It forwards the requests to the backend using either http/json ("""REST""") in case of python backend, or gRPC in case of go backend.

### The backend container(s)

For the backend - and the clip making - to work, the backend needs to be able to call the **youtube-dl** and the **ffmpeg** utilities, so these must be installed and in the PATH.

#### Python

Simple API that expects a json request and tries its best to create the clips using the aforementioned utilities. The API itself is using the FastAPI framework. Nothing too interesting.

#### Go

Created only to do something in golang. Uses the gRPC. To test it first lets spin up the server:

```bash
# lets go into the directory
cd backend/go

# run the server
CLIPS_DIR=/home/pusztai/Videók/clipper/clips SRC_DIR=/home/pusztai/Videók/clipper/sources PORT=5000 go run cmd/main.go
```

Lets test the API with the [**grpcurl**](https://github.com/fullstorydev/grpcurl) tool.

```bash
# in the git repo root
grpcurl -plaintext -d '{"url":"<url>", "clips":[{"player":"<PlayerName>", "start_time": "<StartTime>", "end_time": "<EndTime>"}]}' -proto clips/clips.proto localhost:5000 clips.Clips.NewClip
```

To generate the go stubs from the proto file, run the _generate_go_protoc.sh_ command in the git root.
