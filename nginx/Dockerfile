FROM nginx:stable-alpine

# Copy configuration
COPY default.conf /etc/nginx/conf.d/default.conf

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh

# Make entrypoint executable
RUN chmod +x /entrypoint.sh

EXPOSE 80

# Use our custom entrypoint
ENTRYPOINT ["/entrypoint.sh"]
