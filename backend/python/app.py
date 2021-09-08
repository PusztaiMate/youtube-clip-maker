from ffmpeg import FfmpegCutCmd, FfmpegTimeStamp
import logging
import os

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from downloader import download_video
from utils import (
    find_video_based_on_url,
    get_video_title_from_filename_and_url,
    strip_list_from_video_url,
)
from models import YoutubeClip, validate_clip_fields

origins = ["*"]

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods="*",
    allow_headers="*",
)


CLIP_DIR = os.getenv(
    "CLIP_DIR", os.path.join(os.getenv("HOME"), os.path.join("clipper", "clips"))
)
SRC_DIR = os.getenv(
    "SRC_DIR", os.path.join(os.getenv("HOME"), os.path.join("clipper", "source_videos"))
)


@app.get("/")
def index():
    return {"hello": "world"}


@app.post("/clip")
def clip(ytv: YoutubeClip):
    if not ytv.url:
        return {"message": "url is missing"}, 400

    os.makedirs(SRC_DIR, exist_ok=True)
    download_video(strip_list_from_video_url(ytv.url), SRC_DIR)

    try:
        video_file = find_video_based_on_url(ytv.url, SRC_DIR)
    except FileNotFoundError as e:
        return {"message": "issue when downloading, please try again"}, 500

    video_title = get_video_title_from_filename_and_url(video_file, ytv.url).replace(
        " ", "_"
    )

    for clip in ytv.clips:
        if not validate_clip_fields(clip):
            logging.error(f"clip {clip} is not valid, skipping")
            continue

        player_clip_dir = os.path.join(CLIP_DIR, clip.player).replace(" ", "_")
        os.makedirs(player_clip_dir, exist_ok=True)

        try:
            ts = FfmpegTimeStamp.from_start_end(clip.start, clip.end)
        except ValueError as e:
            logging.error(f"error, when creating clip '{clip}': '{str(e)}'")
            continue
        out_name = f"{video_title}_{clip.start}_{clip.end}.mp4"
        out_path = os.path.join(player_clip_dir, out_name)
        c = FfmpegCutCmd(video_file, out_path, ts)
        c.do()

    return {"message": "success"}, 200
