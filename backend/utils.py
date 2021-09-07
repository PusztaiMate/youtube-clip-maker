import os


def find_video_based_on_url(url: str, src_dir: str) -> str:
    url = strip_list_from_video_url(url)
    vid_id = get_video_id_from_url(url)

    for f in os.listdir(src_dir):
        if vid_id in f:
            return os.path.join(src_dir, f)

    raise FileNotFoundError(f"file not found for url {url}")


def get_video_title_from_filename_and_url(file_path: str, url: str) -> str:
    # lets use just the video id for now
    video_id = get_video_id_from_url(url)
    return video_id


def get_video_id_from_url(url: str) -> str:
    url = url.split("&")[0]
    return url.split("?v=")[1]


def strip_list_from_video_url(url: str) -> str:
    return url.split("&list")[0]
