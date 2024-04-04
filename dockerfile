FROM python:3.12.2-slim

# Configure Poetry
ENV POETRY_VERSION=1.8.2
ENV POETRY_HOME=/opt/poetry

# Install poetry
RUN pip3 install poetry==${POETRY_VERSION}

# Add `poetry` to PATH
ENV PATH="${PATH}:${POETRY_HOME}/bin"

WORKDIR /async-user-service

# Install dependencies
RUN poetry config virtualenvs.create false
COPY user_service/poetry.lock user_service/pyproject.toml ./
RUN poetry install

COPY user_service/. /async-user-service
EXPOSE 8080
EXPOSE 5432
RUN chmod 755  docker-entrypoint.sh
ENV PYTHONPATH /async-user-service

ENTRYPOINT ["./docker-entrypoint.sh"]
