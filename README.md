<div align="center" id="top"> 
  <img src="./.github/app.gif" alt="LocationGrabber" />

&#xa0;

  <!-- <a href="https://locationgrabber.netlify.app">Demo</a> -->
</div>

<h1 align="center">LocationGrabber Telegram Bot</h1>

<p align="center">
  <img alt="Github top language" src="https://img.shields.io/github/languages/top/johnmerga/locationgrabber?color=56BEB8">

  <img alt="Github language count" src="https://img.shields.io/github/languages/count/johnmerga/locationgrabber?color=56BEB8">

  <img alt="Repository size" src="https://img.shields.io/github/repo-size/johnmerga/locationgrabber?color=56BEB8">

  <img alt="License" src="https://img.shields.io/github/license/johnmerga/locationgrabber?color=56BEB8">

  <!-- <img alt="Github issues" src="https://img.shields.io/github/issues/johnmerga/locationgrabber?color=56BEB8" /> -->

  <!-- <img alt="Github forks" src="https://img.shields.io/github/forks/johnmerga/locationgrabber?color=56BEB8" /> -->

  <!-- <img alt="Github stars" src="https://img.shields.io/github/stars/johnmerga/locationgrabber?color=56BEB8" /> -->
</p>

<!-- Status -->

<!-- <h4 align="center">
	🚧  LocationGrabber 🚀 Under construction...  🚧
</h4>

<hr> -->

<p align="center">
  <a href="#dart-about">About</a> &#xa0; | &#xa0; 
  <a href="#sparkles-features">Features</a> &#xa0; | &#xa0;
  <a href="#rocket-technologies">Technologies</a> &#xa0; | &#xa0;
  <a href="#white_check_mark-requirements">Requirements</a> &#xa0; | &#xa0;
  <a href="#checkered_flag-starting">Starting</a> &#xa0; | &#xa0;
  <a href="#memo-license">License</a> &#xa0; | &#xa0;
  <a href="https://github.com/johnmerga" target="_blank">Author</a>
</p>

<br>

## About

this bot takes a location from telegram group chat and once you reply to the location with any caption it will save latitude, longitude, username and name to a google sheet.
(basically this bot created to automate the process of saving locations to google sheet for siinqee bank)

## Features

- [x] save location to google sheet
- [x] save username and name to google sheet
- [x] save latitude and longitude to google sheet
- [x] save caption (Branch name) to google sheet

## Technologies

The following tools were used in this project:

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Google Sheets API](https://developers.google.com/sheets/api)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Google Cloud Platform](https://cloud.google.com/)
- [Google Cloud Platform API](https://cloud.google.com/apis)
- [Google Cloud Platform Service Account](https://cloud.google.com/iam/docs/creating-managing-service-accounts)
- [Google Cloud Platform Service Account Key](https://cloud.google.com/iam/docs/creating-managing-service-account-keys)

## Requirements

Before starting 🏁 , you need to have [Git](https://git-scm.com), [Go](https://golang.org/), [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) installed.

## 🏁 Starting

```bash
# Clone this project
$ git clone https://github.com//locationGrabber

# Access
$ cd locationgrabber

# Check if Docker and docker-compose are installed
$ docker -v
$ docker-compose -v
# create .env file
$ cp .env.example .env
# save your google cloud platform api key file as gkey.json
# start docker containers
$ docker-compose -f docker-compose.prod.yml up -d
# have fun

```

## License

This project is under license from MIT. For more details, see the [LICENSE](LICENSE.md) file.

Made with :heart: by <a href="https://github.com/johnmerga" target="_blank">John Merga</a>

&#xa0;

<a href="#top">Back to top</a>
