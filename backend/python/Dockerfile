FROM python:3.9-buster

WORKDIR /app

RUN apt-get update && apt-get upgrade -y
RUN apt-get install ffmpeg -y

COPY requirements.txt ./
RUN pip install -r requirements.txt

COPY *.py ./

CMD [ "gunicorn", "--bind", "0.0.0.0:8000", "app:app", "-k", "uvicorn.workers.UvicornWorker" ]