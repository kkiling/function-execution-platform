FROM python:3.9-alpine

ENV INPUT_DATA='{}'
ENV PARAMETERS='{}'
ENV EXECUTION_ID = ''
ENV BOOTSTRAP_SERVERS = ''

WORKDIR /app

COPY ./lib ./lib
RUN pip install ./lib/function_execute_platform/
RUN rm -r ./lib/

COPY ./calculator/requirements.txt ./requirements.txt
RUN pip3 install -r ./requirements.txt

COPY ./calculator/ .

ENTRYPOINT python3 main.py