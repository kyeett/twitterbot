
# Twitterbot

Simple twitter bot based on Francesc Compoys tutorial at #justforfunc [Youtube: justforfunc #14: a twitter bot and systemd](https://www.youtube.com/watch?v=SQeAKSJH4vw&t=1792s)

## Running the application

The following environment variables needs to be set to run the application

```bash
TWITTER_CONSUMER_KEY="<YOUR_CONSUMER_KEY>"
TWITTER_CONSUMER_SECRET="<YOUR_CONSUMER_SECRET>"
TWITTER_ACCESS_TOKEN="<YOUR_ACCESS_TOKEN>"
TWITTER_ACCESS_TOKEN_SECRET="<YOUR_ACCESS_TOKEN_SECRET>"
```

## Running the application as a systemd service

Create `twitterbot.env` and add the environment variables from above

### Build the application

```bash
GOOS=linux go build .
```

### ! (Before uploading)

You need to update the *twitterbot.service* to match your users, and the arguments keywords to monitor. In my case -k #golang -k #justforfunc

### Upload binary and configuration to google cloud

Scp to gcloud VM, **twitterbot** in my case.

```bash
gcloud compute scp ./twitterbot*  $USER@twitterbot:~/
```

On the *cloud VM instance*

```bash
sudo mv twitterbot.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl restart twitterbot
sudo systemctl status twitterbot
```