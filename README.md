# discord-scam-detection
Experimental detection, and proof-of-concept, of Discord Nitro phishing/scam links via analyzing images.

**This repository is not considered production-ready. Use with caution.**

## Setup
Set environment variable `BOT_TOKEN` to your bot's OAuth 2 token and run.

## Background

### Initial thoughts
A few days ago (January 2022), I had a thought of detecting Nitro phishing links
by analyzing the image in the message embed. Most Nitro scam websites use the same,
or similar, images, and all of them are Discord's images (some images even link to Discord's website
directly, and aren't even self-hosted.) I spent about 2 hours total making this small
project to test out my theory, and after testing against two scam websites has proven effective.

### Method of detection
This works by analyzing the image of a URL, if the URL has an embedded image.
Malicious image hashes are stored in `image_analyzer.go`. The hashes are
hex-encoded SHA-256 hashes of the image content. If the URL is `discord.com`,
the message is obviously ignored for detection since the image(s) will be legitimate.

### Why this method instead of pre-existing methods?
This method of detection is much more effective than checking against a list of
domains, such as [this repo](https://github.com/nikolaischunk/discord-phishing-links),
since scammers can easily register a new domain for only a few dollars (or even less)
and host a new website for free. In this case, any manually-created list of domains
becomes ineffective until manually updated again; which could take several hours, or even days,
before the new domain becomes widespread.

Obviously, this method can also be circumvented by using different images that are not
listed in the list of hashes. This method can also be circumvented by changing a single pixel
in the image, since the entire SHA-256 hash would change. However, this is much more effective
than the aforementioned method of checking against a list of domains, and is future-proof against new
domains (assuming new websites use images in the list of malicious hashes.)
