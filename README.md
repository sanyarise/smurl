
![GitHub language count](https://img.shields.io/github/languages/count/sanyarise/smurl)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/sanyarise/smurl)
![GitHub repo file count](https://img.shields.io/github/directory-file-count/sanyarise/smurl)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sanyarise/smurl)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/sanyarise/smurl)
![GitHub contributors](https://img.shields.io/github/contributors/sanyarise/smurl)
![GitHub last commit](https://img.shields.io/github/last-commit/sanyarise/smurl)

<img align="right" width="50%" src="./static/images/gopher.png">

# smurl 
## Description
smurl (short for small url) is a service that allows the user to turn a long, cumbersome url address into a beautiful, convenient little link that can be easily used on social network pages, websites, etc. The service also allows you to track the statistics of clicks and clicks on short links for further analysis of the effectiveness of their use.
Chi was chosen as a router because of its idiomatic nature, speed, compliance of handlers with the standard library, sufficiency of tools, lack of hidden context, and a large number of standard middleware.

The API implements 4 main endpoints:
- GET / -home page
- GET /r/{small_url} -search for a small url, update statistics, redirect to the corresponding long address
- POST /create -creating a small url, creating an admin url, writing information about a small, admin and long url to the database
- GET /s/{admin_url} -get statistics on clicks on the received admin url

Postgresql database selected as storage

## HOWTO

- launch with `make run`

- To start, you need to enter the command: make run (there will be a check with testing, creating an executable file, starting the service in docker-compose)

- Then open your browser and type http://localhost:1234

## Video

<img src="./static/images/video.gif">

 
