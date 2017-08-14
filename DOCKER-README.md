## Running the app using docker-compose

Create an entry in `/etc/hosts`

```bash
$ echo '127.0.0.1       concierge.razorpay.dev' | sudo tee --append /etc/hosts > /dev/null
```

## Running it

If everything is in order, then running it would be

```bash
$ docker-compose -f docker-compose.yml up --build
```
