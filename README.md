# boss

boss is a tool to assign a reviewer to Pull Request as an assignee when a new Pull Request is created.

# Usage

## build

```
git clone git@github.com:hiroakis/boss.git
cd boss
make build
 -> Build as a linux/386 binary. If you'd like to build for the other platform, edit Makefile or set GOOS environment variable before you run go build.
```

## run

```
boss
```

* options

```
-c: The listen IP. Default: 0.0.0.0 
-h: The listen port. Default: 9000
-p: The configration file. Default: ./config.yml
```

## Setting up a GitHub WebHook

1. Go to Settings page of your repository, and click on `Webhooks` .
2. Enter the URL that you have run `boss` to the `Payload URL` .
3. Choose the `application/json` content type.
4. Check the `Send me everything.` box.
5. Click on `Add webhook`

## Configuration File Format

```
repos:
  hiroakis/boss: # repository full name
    token: abcdefghij1234567890abcdefghij1234567890 # github token(require write privilege to the repository)
    members: # candidates for assignee
      - hiroakis
      - zawinul
      - jarrett
      - corea
      - hancock
    labels: # labels that you'd like to add new pull request
      - Enhancement
      - waiting for review
  hiroakis/twitter-streaming:
    token: abcdefghij1234567890abcdefghij1234567890
    members:
      - hiroakis
      - jaco
      - mingus
      - collins
      - flea
      - marcus
    labels:
      - waiting for review
```

# Specification

* If hiroakis creates a new pull request, hiroakis will not be assigned.
* If an assignee is added when a new pull request is created. Add labels, but doesn't assign.

# Lisence

MIT
