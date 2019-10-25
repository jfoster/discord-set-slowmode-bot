# discord-slowmode-bot

Simple discord bot to set channel slowmode to any value between 1 second and 6 hours.

## Why?

I wanted to set the discord slowmode setting to 1 second to not be too restrictive for people who\
talk\
like\
this\
but still try and thwart spambots.

Currently the GUI for setting slowmode is not particularly granular and does not allow for below 5 seconds.

![](https://i.imgur.com/5ki1rDd.png)

I then discovered it's possible to set slowmode to any integer between 0 to 21600 seconds via the discord api, hence the creation of this simple bot.

## How?

1. Create a discord bot app, [instructions here](https://github.com/andersfylling/disgord/wiki/Get-bot-token-and-add-it-to-a-server).
1. Download the [latest release](https://github.com/jfoster/discord-slowmode-bot/releases/latest) for your given platform.
1. Open a command line instance in the bot's directory and run ```./discord-slowmode-bot```, a warning should be printed ```client id is not specified, check cfg.yaml file``` and cfg.yaml will be created.
1. Copy your bot token into cfg.yaml replacing ```<your-bot-token-here>```.
1. Run ```./discord-slowmode-bot``` again, once connected, copy the discord invite link to your favourite browser and add the bot to a server. The bot should now be present in the desired server.  ![](https://transfer.sh/RkNk3/Screenshot-2019-04-02-at-18.09.52.png)
2. In discord, from the channel you would like to set slowmode, type ```@SlowModeBot <duration>``` e.g. ```@SlowModeBot#8558 1s``` for 1 second, ```@SlowModeBot#8558 2m``` for 2 minutes or ```@SlowModeBot#8558 2h``` for 2 hours.  ![](https://i.imgur.com/bSdpfMC.png)

## Credits

Many thanks to @andersfylling for creating [disgord](https://github.com/andersfylling/disgord), and for offering guidance! 🍻

## License

[MIT](LICENSE.txt)