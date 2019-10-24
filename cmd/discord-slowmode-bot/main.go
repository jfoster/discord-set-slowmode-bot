package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/std"
	"github.com/smallfish/simpleyaml"
	"github.com/urfave/cli"

	"github.com/jfoster/discord-slowmode-bot/internal/log"
)

const (
	cfgfilename     = "cfg.yaml"
	cfgfiletemplate = `token: <your-bot-token-here>`
)

var (
	logr    = log.New()
	version = "# filled by go build #"
)

func main() {
	if err := cliApp(); err != nil {
		logr.Fatal(err)
	}
}

func cliApp() error {
	app := cli.NewApp()
	app.Name = "discord-set-slowmode-bot"
	app.Usage = ""
	app.Version = version
	app.Authors = []cli.Author{
		{
			Name:  "Jacob Foster",
			Email: "jacobfoster@pm.me",
		},
	}
	app.Compiled = time.Now()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "Specify bot token.",
		},
	}
	app.Action = func(c *cli.Context) error {
		token := c.String("token")
		// get token from cfg if not specified with flag
		if token == "" {
			cfg, err := getCfg()
			if err != nil {
				return err
			}

			token, err = cfg.Get("token").String()
			if err != nil {
				return err
			}
		}

		if token == "" || token == "<your-bot-token-here>" {
			return fmt.Errorf("client token is not specified, check %s file or specify with -t flag", cfgfilename)
		}

		if err := runBot(token); err != nil {
			return err
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}

func getCfg() (*simpleyaml.Yaml, error) {
	yaml, err := simpleyaml.NewYaml([]byte(cfgfiletemplate))
	if err != nil {
		return yaml, err
	}

	if _, err := os.Stat(cfgfilename); os.IsNotExist(err) {
		logr.Warnf("%s file does not exist, creating %s file from template...", cfgfilename, cfgfilename)
		err = ioutil.WriteFile(cfgfilename, []byte(cfgfiletemplate), 0644)
		if err != nil {
			return yaml, err
		}
	}

	file, err := ioutil.ReadFile(cfgfilename)
	if err != nil {
		return yaml, err
	}

	yaml, err = simpleyaml.NewYaml(file)
	if err != nil {
		return yaml, err
	}
	return yaml, nil
}

func runBot(token string) error {
	logr.Info("Creating Discord session")

	config := &disgord.Config{
		BotToken: token,
		ProjectName: "jfoster/discord-slowmode-bot",
	}
	if logr.Debug {
		config.Logger = logr.Logger
	}

	bot, err := disgord.NewClient(config)
	if err != nil {
		return err
	}

	bot.AddPermission(disgord.PermissionManageChannels)

	bot.On(event.Ready, onReady)

	filter, err := std.NewMsgFilter(bot)
	if err != nil {
		return err
	}
	bot.On(event.MessageCreate, filter.HasBotMentionPrefix, onMessageCreate)

	logr.Info("Connecting")
	start := time.Now()
	if err = bot.Connect(); err != nil {
		return err
	}
	logr.Info("Connection took ", time.Since(start))

	logr.Info("Press Ctrl+C to exit")
	bot.DisconnectOnInterrupt()

	return nil
}

func onReady(session disgord.Session, evt *disgord.Ready) {
	guilds, err := session.GetGuilds(&disgord.GetCurrentUserGuildsParams{
		Limit: 999,
	})
	if err != nil {
		logr.Error(err)
	} else {
		guildString := "Connected guilds:\n"
		for i, guild := range guilds {
			guildString += fmt.Sprintf("%d: %s (%s)\n", i+1, guild.Name, guild.ID)
		}
		logr.Info(guildString)
	}

	user, err := session.GetCurrentUser()
	if err != nil {
		logr.Error(err)
	} else {
		url := fmt.Sprintf("https://discordapp.com/oauth2/authorize?scope=bot&client_id=%s&permissions=%d", user.ID, session.GetPermissions())
		logr.Infof("Link to add bot to your guild:\n%s", url)
	}
}

func onMessageCreate(session disgord.Session, evt *disgord.MessageCreate) {
	message := evt.Message

	channel, err := session.GetChannel(message.ChannelID)
	if err != nil {
		logr.Error(err)
		return
	}

	guild, err := session.GetGuild(channel.GuildID)
	if err != nil {
		logr.Error(err)
		return
	}

	member, err := session.GetMember(guild.ID, message.Author.ID)
	if err != nil {
		logr.Error(err)
		return
	}

	isPermitted := false
	for _, roleID := range member.Roles {
		role, err := guild.Role(roleID)
		if err != nil {
			logr.Error(err)
		} else if (role.Permissions & (disgord.PermissionAdministrator | disgord.PermissionManageChannels)) != 0 {
			isPermitted = true
			break
		}
	}
	if guild.OwnerID == member.User.ID {
		isPermitted = true
	}
	if !isPermitted {
		logr.Infof("%s#%s (%s) is not permitted!", member.User.Username, member.User.Discriminator.String(), member.Nick)
		return
	}

	contents := strings.Split(message.Content, " ")
	if len(contents) < 2 {
		return
	}

	ratelimit := contents[1]
	if ratelimit == "?" {
		duration := time.Duration(channel.RateLimitPerUser * uint(time.Second))
		message.Reply(session, fmt.Sprintf("slowmode is set to %s", duration.String()))
	} else {
		duration, err := time.ParseDuration(ratelimit)
		if err != nil {
			logr.Error(err)
			return
		}

		secs := uint(duration.Seconds())
		if secs < 0 || secs > 21600 {
			logr.Error(fmt.Errorf("Specified duration is not between 0-21600 seconds"))
			return
		}

		ucb := session.UpdateChannel(channel.ID).SetRateLimitPerUser(secs)
		_, err = ucb.Execute()
		if err != nil {
			logr.Error(err)
			return
		}

		message.Reply(session, fmt.Sprintf("slowmode set to %s", duration.String()))
	}
}
