from collections import namedtuple
from subprocess import check_output as co
from datetime import timedelta
import logging
from typing import Tuple


class FffmpeCmdException(RuntimeError):
    pass


class FfmpegTimeStamp:
    def __init__(self, ts: str, length: str):
        self.ts = ts
        self.length = length

    @classmethod
    def from_start_end(cls, start: str, end: str) -> "FfmpegTimeStamp":
        cls._validate_format(start)
        cls._validate_format(end)

        start_min, start_sec = cls._get_mins_and_secs(start)
        ts1 = timedelta(minutes=start_min, seconds=start_sec)

        end_min, end_sec = cls._get_mins_and_secs(end)
        ts2 = timedelta(minutes=end_min, seconds=end_sec)

        return cls(str(ts1), str(ts2 - ts1))

    @staticmethod
    def _validate_format(s: str):
        sep = ":"
        s = s.replace(".", ":")
        mins, secs = s.split(sep)
        try:
            int(mins)
            int(secs)
        except ValueError:
            ValueError(f"'{s}' is not in expected format 'mm:ss' or 'mm.ss'")

    def _get_mins_and_secs(s: str) -> Tuple[int, int]:
        sep = ":" if ":" in s else "."
        mins, secs = s.split(sep)
        return int(mins), int(secs)


class FfmpegCutCmd:
    def __init__(self, inp: str, out: str, ts: FfmpegTimeStamp):
        if not inp or not out:
            raise FffmpeCmdException(
                f"output ('{out}') or input ('{inp}') is not specified"
            )

        self.inp = inp
        self.out = out
        self.command = "ffmpeg"
        self.ts = ts

    def do(self):
        self.command += self._get_timestap_command(self.ts.ts, self.inp, self.ts.length)
        # self.command += " -c copy"
        self.command += f" -y {self.out}"

        logging.info(f"executing command {self.command}")
        logging.info(co(self.command, shell=True))

    @staticmethod
    def _get_timestap_command(ts: str, inp: str, length: str) -> str:
        return f" -ss {ts} -i '{inp}' -t {length}"


# for testing
if __name__ == "__main__":
    cmd = FfmpegCutCmd(
        "inp_5_secs.mp4", "out.mp4", FfmpegTimeStamp("00:00:01", "00:00:02")
    )
    cmd.do()
