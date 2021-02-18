# cara
I wrote this bot just for the purpose of automating some tasks on my private Discord server, but decided to upload the code because why not

Written entirely in Go, features some features like 
* A database of blacklisted strings which will prevent anyone from making a role with that string in it
* Snarky help and about embeds
* Vomit command
* Banned nickname fixer
* Self-service nickname changer (for the bot)
* Self-service status changer (also for bot)
* Say command
* Hidden privilage escalation command

Note that most of these are hardcoded to my Discord account ID, I do not recommend that you deploy this on your server unless you modify the code extensively

### License
Under BSD-3-Clause because why not

### Fun fact
My .gitignore was not working and I accidentally uploaded the bot token in it's startup script, and Discord automatically reset it, which led to me learning about their token-sniffing service
