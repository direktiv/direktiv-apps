FROM tiangolo/uvicorn-gunicorn-fastapi:python3.7

RUN pip install requests

ENV HOST=0.0.0.0
ENV PORT=8080

COPY ./app /app
