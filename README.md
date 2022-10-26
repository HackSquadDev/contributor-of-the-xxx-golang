## COTO [Contributor of The Organization] Bot

> A bot that runs in a fixed time interval

- checks all the repositories in an organizaiton utilizing GitHub's GraphQL api
- finds the most amount of PR for a single contributor
- creates a banner for the person
- posts the banner to twitter and discord

## Roadmap

- [x] A basic HTTP server
- [x] Github GraphQL API to fetch the top contributor
- [ ] Design a template for badges
- [ ] Dynamically add the data of the contributor to generate the image
- [ ] Make a Gif (animated)
- [ ] Run a cron job every month for the above processes
