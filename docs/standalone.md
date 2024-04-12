# Development (Standalone)

Deploy Rasbora as standalone:

```bash
git clone https://github.com/openseawave/rasbora.git

cd rasbora/

# Modify the .env file to fit your system
# Also you will need to get your license key from our website
# [https://rasbora.openseawave.com/](https://rasbora.openseawave.com/license)

cp .env.example .env

# Run Docker Compose
docker-compose up -d
```

Note: Please ensure that you get your license key and Docker Compose,Docker and Git is installed on your system before running the above commands.
