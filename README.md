# Concierge-Ingress

It allows us to create leases for incoming connections to Kubernetes Ingress Objects.

## Installation

1. Obtain Google client secret and client id for OAuth. 

2. Add `http://127.0.0.1:4180` in **Authorized JavaScript origins** and `http://127.0.0.1:4180/oauth2/callback` in **Authorized redirect URIs**.

3. Get the source code on your machine via git.

    ```shell
    git clone --branch ingress-concierge https://github.com/razorpay/concierge.git
    ```

4. Rename file `.env.example` to `.env` and change credentials.

    ```shell
    mv .env.example .env
    ```

5. Rename file `oauth2_proxy.example.cfg` to `oauth2_proxy.cfg` in `oauth2_proxy` dir and add the following values obtained from Google OAuth

    ```text
    client_id =
    client_secret =
    ```

6. Make sure to change the `seeding.go` file in `database` dir with your information else you won't be able to login.

7. Run docker-compose command to run the application.

    ```shell
    docker-compose up -d
    ```

8. Check for the running application.

    ```shell
    docker ps
    ```

9. Access the application on [127.0.0.1:4180](127.0.0.1:4180)
