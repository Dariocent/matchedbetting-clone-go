# Matched Betting Tool

This Go web server provides an **oddsmatcher**, a matched betting tool that compares odds between bookmakers' websites to offer users the best rating for their bets.
If you run it, it is available at http://localhost:8080/betting/

Note:
 - This project is just a way to get more familiar with the language Go. It's in no way a production tool.

## Features

- **Oddsmatcher:** The web server includes an Oddsmatcher tool that fetches and compares odds from various bookmakers.
- **Rating System:** Utilizes a rating system to help users identify the best odds for matched betting.
- **Dynamic Updates:** Users can refresh the Oddsmatcher tool to get the latest odds without refreshing the entire page, thanks to htmx.

## Challenges

- **Webscraping:** This tool uses GoQuery to scrape odds information on website that offer it in their HTML page. For other websites that use JSON as a medium to convey the information to the frontend, it reads those JSON files by sending a request and parsing the data.
- **Concurrency:** My aim is to optimize the code using Go's powerful concurrency.
- **Cloud Native** My aim is to have this application be able to scale easier by building it with Cloud Native in mind.
