FROM node:20-alpine AS builder
ARG VITE_LEAVE_API_BASE=/api/leave
ARG VITE_USER_API_BASE=/api/user
ENV VITE_LEAVE_API_BASE=$VITE_LEAVE_API_BASE
ENV VITE_USER_API_BASE=$VITE_USER_API_BASE
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf.template
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 3000
ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
