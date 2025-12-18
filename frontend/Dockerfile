
FROM node:20-alpine AS build-env
WORKDIR /app


COPY package.json pnpm-lock.yaml ./


RUN npm install -g pnpm


RUN pnpm install --frozen-lockfile


COPY . .


ARG VITE_API_BASE_URL=http://asto-lms.local
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
RUN pnpm run build


FROM nginx:alpine


COPY --from=build-env /app/build/client /usr/share/nginx/html


COPY nginx.conf /etc/nginx/conf.d/default.conf


EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]