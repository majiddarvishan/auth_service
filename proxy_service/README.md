# Auth service

## Get an SSL Certificate

1. Use a Self-Signed Certificate (for testing)
2. Obtain a Certificate from Let’s Encrypt (for production)

### Generate a Self-Signed Certificate

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365
```

### Using Let’s Encrypt

1. Install Certbot

```bash
sudo apt install certbot
```

2. Get an SSL certificate

```bash
certbot certonly --standalone -d yourdomain.com
```
