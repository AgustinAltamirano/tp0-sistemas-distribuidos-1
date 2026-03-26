FROM python:3.14-slim
COPY server /
RUN python -m unittest tests/test_common.py
ENTRYPOINT ["/bin/sh"]
