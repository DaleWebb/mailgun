# Mailgun

Mailgun driver for itsabot/abot/shared/interface/email to send and receive emails

## Abot Setup

Sign up for [Mailgun](https://mailgun.com/signup). For development, you can use your [sandbox domain and sandbox API key](https://mailgun.com/app/domains), given to you when you signup. You'll want to set the following environment variables in your `~/.bash_profile` or `~/.bashrc`:
```bash
export MAILGUN_DOMAIN="REPLACE"
export MAILGUN_API_KEY="REPLACE"
```

We'll also want to update our environment variables on Heroku, so run the following in your terminal:
```bash
heroku config:set MAILGUN_DOMAIN=REPLACE \
MAILGUN_API_KEY=REPLACE
```

Now we'll add the Mailgun driver. Since we've written a plugin for Mailgun like any other, you can simply add it to your plugins.json like so:
```json
{
    "Name": "abot",
    "Version": "0.1.0",
    "Dependencies": {
        "github.com/itsabot/plugin_onboard": "*",
        "github.com/DaleWebb/mailgun": "*"
    }
}
```

Then from your terminal run:
```bash
abot plugin install
```

## Mailgun Configuration

The Mailgun driver will automatically add a special POST /mailgun route to Abot.

Let's also make sure that Mailgun knows how to communicate with our plugin.

For Mailgun to send your emails to Abot, your Abot will need to be hosted online.

Add a route to your Mailgun account that will POST to our Abot's route every time Mailgun receives an email. There are two ways to do that, either:

1. **Via the [Mailgun dashboard](https://mailgun.com/cp/routes#edit)** and add the filter expression: `match_recipient(".*@MAILGUN_DOMAIN")` and add the action: `forward("http://MAILGUN_DOMAIN/mailgun/")``.

2. **Via the Terminal**, you can execute this command.
```bash
curl -s --user 'api:MAILGUN_API_KEY' \
    https://api.mailgun.net/v3/routes \
    -F priority=0 \
    -F description='Sample route' \
    -F expression='match_recipient(".*@MAILGUN_DOMAIN")' \
    -F action='forward("http://MAILGUN_DOMAIN/mailgun/")' \
    -F action='stop()'
```
##Testing it out

To try it out, let's first create an Abot account through the web interface. Go to your Abot domain (heroku open) and click on Sign Up in the header at the top right.

Once you've signed up, send Abot an email to any-name@MAILGUN_DOMAIN. Sometimes responses via Mailgun take a few seconds, but you should get a reply back shortly.

You can also simulate a Mailgun callback via cURL, if you do not have a online Abot.

```bash
curl -H "Content-Type: application/json" -X POST -d '{"stripped-text":"hi", "recipient": "abot@$MAILGUN_DOMAIN", "sender": "you@example.com", "subject": "test"}' http://localhost:4200/mailgun
```
