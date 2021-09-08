from youtube_dl import YoutubeDL


def download_video(video_url: str, video_dir: str):
    opts = {
        "format": "bestaudio/best",
        "download_archive": "./archive",
        "outtmpl": f"{video_dir}/%(id)s.%(ext)s",
    }

    with YoutubeDL(opts) as ydl:
        ydl.download([video_url])
