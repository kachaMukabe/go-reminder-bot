
# Whatsapp Reminder Bot

This is a whatsapp chatbot built with Golang for the Linode hackathon on Hasura.
This chatbot will recieve a message from a user store it as a reminder and send it back
to the user at a later time.


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`PORT`

`FACEBOOK_TOKEN`

`VERIFY_TOKEN`



## Installation

Clone this github repository

```bash
  git clone https://github.com/kachaMukabe/go-reminder-bot
  cd go-reminder-bot
  go mod tidy
  go run .
```
    
## Tech Stack

**Server:** Golang, Gin

**Other:** Facebook cloud api

## Deployment

Read my article [here](https://dev.kachamukabe.com/hosting-a-golang-whatsapp-bot-on-linode) on how to deploy this to a Linode server.


## Acknowledgements

 - [Linode Hackathon](https://www.linode.com/?utm_source=hashnode&utm_medium=article&utm_campaign=hackathon_announcement)
 - [Hashnode](https://hashnode.com/)



## License

[MIT](https://choosealicense.com/licenses/mit/)

