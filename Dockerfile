FROM python:3.10-alpine AS base

# Configure Python environment
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# Start a new stage for building the application dependencies
FROM base AS depbuilder

# Install required building dependencies
RUN apk add --no-cache build-base

# Install and build required Python dependencies
COPY requirements.txt requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Start a new stage for the final application image
FROM base AS final

# Configure image labels
LABEL org.opencontainers.image.source https://github.com/hhromic/traefik-oidcintro-fwdauth

# Configure default entrypoint and exposed port of the application
ENTRYPOINT ["python", "-m", "oidcintro_fwdauth"]
EXPOSE 4181

# Copy application dependency artifacts
COPY --from=depbuilder /usr/local /usr/local

# Configure runtime user/group/home for the application
ARG APP_USER=app
ARG APP_GROUP=app
ARG APP_HOME=/app
RUN addgroup ${APP_GROUP} \
    && adduser -D -h ${APP_HOME} -G ${APP_GROUP} ${APP_USER}
WORKDIR ${APP_HOME}

# Copy application artifacts
COPY oidcintro_fwdauth oidcintro_fwdauth

# Configure default application user and group
USER ${APP_USER}:${APP_GROUP}
