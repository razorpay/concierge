## Running the app using docker-compose

Create an entry in `/etc/hosts`

```bash
$ echo '127.0.0.1       concierge.razorpay.dev' | sudo tee --append /etc/hosts > /dev/null
```

## Running it

If everything is in order, then running it would be

```bash
cp docker-compose.sample.yml docker-compose.dev.yml
# Fix the variables in the docker-compose.dev.yml file
docker-compose -f docker-compose.yml up --build
```

Now you should be able to access <http://concierge.razorpay.dev:28070> in your system. In order to get AWS parts working, you will need the AWS_ keys set correctly.
