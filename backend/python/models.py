from typing import List
from pydantic import BaseModel


class YoutubeClipSection(BaseModel):
    start: str
    end: str
    player: str


class YoutubeClip(BaseModel):
    url: str
    clips: List[YoutubeClipSection]


def validate_clip_fields(c: YoutubeClipSection) -> bool:
    return c.start and c.end and c.player
