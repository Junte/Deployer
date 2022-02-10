package core

import (
	"fmt"
	"github.com/slack-go/slack"
	"log"
)

func notifySlack(component string, component_config *ComponentConfig, failed bool, stdout, stderr string) {
	if Config.Notification.Slack.ApiToken == "" {
		return
	}

	channel := getSlackChannel(component_config)
	if channel == "" {
		return
	}

	slackMessage := buildSlackMessage(
		component,
		failed,
		stdout,
		stderr,
	)

	client := slack.New(Config.Notification.Slack.ApiToken)
	_, _, err := client.PostMessage(channel, slackMessage)
	if err != nil {
		log.Printf("error on slack message send: %v", err)
		return
	}
}

func buildSlackMessage(component string, failed bool, stdout, stderr string) slack.MsgOption {
	message := ""
	if failed {
		message = buildFailMessage(component)
	} else {
		message = buildSuccessMessage(component)
	}

	attachments := []slack.Attachment{}
	attachments = append(attachments, slack.Attachment{
		Title:   ":memo: stdout",
		Pretext: message,
		Color:   "#36a64f",
		Text:    stdout,
	})
	if stderr != "" {
		attachments = append(attachments, slack.Attachment{
			Title: ":fire: strerr",
			Color: "#eb343a",
			Text:  stderr,
		})
	}
	return slack.MsgOptionAttachments(attachments...)
}

func buildFailMessage(component string) string {
	return fmt.Sprintf(":x: Failed component \"%s\" deployment to environment \"%s\"",
		component,
		Config.Environment,
	)
}

func buildSuccessMessage(component string) string {
	return fmt.Sprintf(":white_check_mark: Component \"%s\" was deployed to environment \"%s\"",
		component,
		Config.Environment,
	)
}

func getSlackChannel(component_config *ComponentConfig) string {
	channel := Config.Notification.Slack.Channel
	if component_config.Notification.Slack.Channel != "" {
		channel = component_config.Notification.Slack.Channel
	}

	return channel
}
