# Build context is the repo root (see docker-compose.prod.yml) so this stage can also
# pull in the shared nginx config from ../nginx.
FROM node:22-alpine AS build
WORKDIR /src

COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

COPY frontend/ .
# Same-origin '/api' prefix — nginx strips it before proxying to the backend (see
# nginx/conf.d/app.conf). Keeps prod free of CORS entirely; dev still talks cross-origin
# straight to localhost:8080 via the default in src/api/http.ts.
ENV VITE_API_URL=/api
RUN npm run build

# --- runtime stage ---
FROM nginx:1.27-alpine
COPY --from=build /src/dist /usr/share/nginx/html
COPY nginx/conf.d/ /etc/nginx/conf.d/

EXPOSE 80 443
