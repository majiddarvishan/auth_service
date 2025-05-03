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


Use mkcert (better!)

```bash
# Install mkcert (Linux example)
sudo apt install libnss3-tools
mkcert -install
mkcert localhost
```



swaggertodoc

https://mvnrepository.com/artifact/io.github.swagger2markup/swagger2markup-cli
wget https://repo1.maven.org/maven2/io/github/swagger2markup/swagger2markup-cli/1.3.3/swagger2markup-cli-1.3.3.jar
java -jar /path/to/swagger2markup-cli-1.3.3.jar convert -i <your_swagger_file.json_or_yaml> -d <output_directory>

swagger2markup convert -i ./docs/swagger.json -f output.adoc
pandoc output.adoc -o output.docx


pip install python-docx
python swagger2docx.py
