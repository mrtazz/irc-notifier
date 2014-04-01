# irc-notifier [![Build Status](https://travis-ci.org/mrtazz/irc-notifier.svg?branch=master)](https://travis-ci.org/mrtazz/irc-notifier)

## Overview
This is a small program to fetch my IRC notifications from Redis according to
[the setup I wrote about][1]. It replaces the Python script I mentioned
because I got tired of having Python, pip, py-redis and gntp set up before
being able to receive notifications.

## Usage
```
panthor:irc-notifier:% irc-notifier --help
Usage of irc-notifier:
  -a="": the auth key if set
  -h="": the redis host
  -p=6379: the redis port
  
panthor:irc-notifier:% irc-notifier -h redis.example.com -a lolsecret
```

[1]: http://www.unwiredcouch.com/2012/11/03/irc-notifications-with-logstash.html
