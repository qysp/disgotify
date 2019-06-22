# disgotify
Discord bot for notifications written in Go.

## Getting started
Create a copy of `.env.example` and rename it to `.env`, add your credentials and preferences, compile the source and start the bot.

To get a list of all available commands use `(command prefix)help` (e. g. `+help`). For a more specific help message for a command use `(command prefix)help [command name]` (e. g. `+help remind`).
Note that in a DM channel with the bot, a command prefix is not needed.

## Adding more commands
In order to add your own commands, implement the functions of the `Command` interface and initialize it in the command index. You can use the Ping command as a template.