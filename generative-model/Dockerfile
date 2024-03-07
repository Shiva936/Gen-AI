FROM python:3.7-slim-bullseye
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt
CMD ["python", "./main.py"]

FROM debian:bullseye-slim
WORKDIR /root
COPY --from=0 /app/app ./
RUN apt-get update \
    && apt-get install -y --force-yes --no-install-recommends apt-transport-https curl ca-certificates \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/*
EXPOSE 3001
ENTRYPOINT ["python","./main.py"]