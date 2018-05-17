FROM python
RUN pip install unicornhathd Flask Flask-Sockets
COPY build/ .
COPY server.py .
ENTRYPOINT python server.py