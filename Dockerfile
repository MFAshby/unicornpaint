FROM python:latest
RUN pip install unicornhathd Flask Flask-Sockets
RUN pip install numpy
COPY build/ .
COPY server.py .
ENTRYPOINT python server.py
