FROM python:3.9.4-buster

ARG DEBIAN_FRONTEND=noninteractive

COPY ./requirements.txt /app/requirements.txt
RUN pip install --upgrade pip && \
    pip install -r /app/requirements.txt

COPY ./alembic.ini /app/alembic.ini

COPY ./alembic /app/alembic
COPY ./src /app/src

WORKDIR /app/src

ENV PYTHONPATH /app/src

EXPOSE 80

CMD ["uvicorn", "api:app", "--host", "0.0.0.0", "--port", "80"]
