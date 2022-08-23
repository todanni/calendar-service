Secrets used in this repo

GitHub Actions

Docker credentials 
```shell
username: ${{ secrets.DOCKER_USERNAME }}
password: ${{ secrets.DOCKER_PASSWORD }}
```


Discord webhook
We send notifications to our discord about deployments with the status of the deployment - whether it has succeeded or failed.
```shell
webhook_id: ${{ secrets.DISCORD_WEBHOOK_ID }}
webhook_token: ${{ secrets.DISCORD_WEBHOOK_TOKEN }}

```

Digital Ocean droplet details
To SSH to the droplet and run the docker container, we use the droplet IP and an SSH key. 

```shell
host: ${{ secrets.DROPLET_IP }}
key: ${{ secrets.ADMIN_SSH_KEY }}
```

Docker network name
??
```shell
--network {{ secrets.NETWORK }}
```


In code secrets

Vault address
Vault token 

Google Calendar JSON Credentials