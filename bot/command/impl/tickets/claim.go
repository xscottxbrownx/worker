package tickets

import (
	"fmt"

	"github.com/TicketsBot-cloud/common/permission"
	"github.com/TicketsBot-cloud/worker/bot/command"
	"github.com/TicketsBot-cloud/worker/bot/command/registry"
	"github.com/TicketsBot-cloud/worker/bot/constants"
	"github.com/TicketsBot-cloud/worker/bot/customisation"
	"github.com/TicketsBot-cloud/worker/bot/dbclient"
	"github.com/TicketsBot-cloud/worker/bot/logic"
	"github.com/TicketsBot-cloud/worker/i18n"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/interaction"
)

type ClaimCommand struct {
}

func (ClaimCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "claim",
		Description:     i18n.HelpClaim,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Support,
		Category:        command.Tickets,
		Timeout:         constants.TimeoutOpenTicket,
	}
}

func (c ClaimCommand) GetExecutor() interface{} {
	return c.Execute
}

func (ClaimCommand) Execute(ctx registry.CommandContext) {
	// Get ticket struct
	ticket, err := dbclient.Client.Tickets.GetByChannelAndGuild(ctx, ctx.ChannelId(), ctx.GuildId())
	if err != nil {
		ctx.HandleError(err)
		return
	}

	// Verify this is a ticket channel
	if ticket.UserId == 0 {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNotATicketChannel)
		return
	}

	// Check if thread
	ch, err := ctx.Worker().GetChannel(ctx.ChannelId())
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if ch.Type == channel.ChannelTypeGuildPrivateThread {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageClaimThread)
		return
	}

	if err := logic.ClaimTicket(ctx, ctx, ticket, ctx.UserId()); err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ReplyPermanent(customisation.Green, i18n.TitleClaimed, i18n.MessageClaimed, fmt.Sprintf("<@%d>", ctx.UserId()))
}
