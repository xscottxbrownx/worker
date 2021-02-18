package listeners

import (
	"context"
	"encoding/json"
	"github.com/TicketsBot/common/premium"
	"github.com/TicketsBot/common/sentry"
	"github.com/TicketsBot/worker"
	"github.com/TicketsBot/worker/bot/command"
	"github.com/TicketsBot/worker/bot/dbclient"
	"github.com/TicketsBot/worker/bot/errorcontext"
	"github.com/TicketsBot/worker/bot/logic"
	"github.com/TicketsBot/worker/bot/utils"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"golang.org/x/sync/errgroup"
	"strconv"
)

func OnPanelReact(worker *worker.Context, e *events.MessageReactionAdd) {
	errorContext := errorcontext.WorkerErrorContext{
		Guild:   e.GuildId,
		User:    e.UserId,
		Channel: e.ChannelId,
		Shard:   worker.ShardId,
	}

	// In DMs
	if e.GuildId == 0 {
		return
	}

	if e.UserId == worker.BotId || e.Member.User.Bot {
		return
	}

	// Get panel from DB
	panel, err := dbclient.Client.Panel.Get(e.MessageId)
	/// debug start
	extra := make(map[string]interface{})
	extra["panel"], _ = json.Marshal(panel)
	extra["err"] = err
	sentry.LogWithTags("react debug", extra, map[string]string{"guild_id": strconv.FormatUint(e.GuildId, 10)})
	/// debug end
	if err != nil {
		sentry.ErrorWithContext(err, errorContext)
		return
	}

	// Verify this is a panel
	if panel.MessageId != 0 {
		emoji := e.Emoji.Name // This is the actual unicode emoji (https://discordapp.com/developers/docs/resources/emoji#emoji-object-gateway-reaction-standard-emoji-example)

		// Check the right emoji ahs been used
		if panel.ReactionEmote != emoji && !(panel.ReactionEmote == "" && emoji == "📩") {
			return
		}

		// TODO: Check perms
		// Remove the reaction from the message
		if err := worker.DeleteUserReaction(e.ChannelId, e.MessageId, e.UserId, emoji); err != nil {
			sentry.LogWithContext(err, errorContext)
		}

		var blacklisted bool
		var premiumTier premium.PremiumTier

		group, _ := errgroup.WithContext(context.Background())

		// get blacklisted
		group.Go(func() (err error) {
			blacklisted, err = dbclient.Client.Blacklist.IsBlacklisted(e.GuildId, e.UserId)
			return
		})

		// get premium
		group.Go(func() error {
			premiumTier = utils.PremiumClient.GetTierByGuildId(e.GuildId, true, worker.Token, worker.RateLimiter)
			return nil
		})

		if err := group.Wait(); err != nil {
			sentry.ErrorWithContext(err, errorContext)
			return
		}

		if blacklisted {
			return
		}

		panelContext := command.NewPanelContext(worker, e.GuildId, e.ChannelId, e.UserId, premiumTier)

		go logic.OpenTicket(&panelContext, &panel, panel.Title)
	}
}
