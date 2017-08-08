## Running the app using docker-compose

To run locally, you need to have a `~/.aws` folder. The steps should be followed as mentioned in the docs here [https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html)

You just need to add a small change in your `docker-compose.yml` here

```yaml
version: '2'
services:
  concierge-web:
    build:
      context: .
      dockerfile: Dockerfile
    links:
      - db-concierge
    depends_on: 
      - db-concierge
    ports: 
      - 80:80
    environment:
      DB_HOST: db-concierge
      DB_DATABASE: concierge
      DB_USERNAME: concierge
      DB_PASSWORD: concierge
    volumes:
      - ~/.aws:/var/lib/nginx/.aws

  db-concierge:
    image: mysql:5.6
    ports:
      - 3306:3306
    volumes:
      - ./docker/data:/docker-entrypoint-initdb.d
    environment:
      MYSQL_ROOT_PASSWORD: concierge
      MYSQL_DATABASE: concierge
      MYSQL_USER: concierge
      MYSQL_PASSWORD: concierge
```

Create an entry in `/etc/hosts` 

```bash
$ echo '127.0.0.1       aws-manage.dev' | sudo tee --append /etc/hosts > /dev/null
```

## Running it

If everything is in order, then running it would be 

```bash
$ docker-compose -f docker-compose.yml up --build
```

