package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.mau.fi/whatsmeow/types"
)

func newChatsClearCmd(flags *rootFlags) *cobra.Command {
	var jid string
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all messages in a chat (propagates to all devices)",
		Long: `Send a "clear chat" action via WhatsApp app-state sync.
This removes all messages from the specified chat on your phone
and all linked devices. It does NOT delete messages for other participants.

Requires an active WhatsApp connection (will connect temporarily).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if jid == "" {
				return fmt.Errorf("--jid is required")
			}

			target, err := types.ParseJID(jid)
			if err != nil {
				return fmt.Errorf("invalid JID %q: %w", jid, err)
			}

			ctx, cancel := withTimeout(context.Background(), flags)
			defer cancel()

			a, lk, err := newApp(ctx, flags, true, false)
			if err != nil {
				return err
			}
			defer closeApp(a, lk)

			if err := a.Connect(ctx, false, nil); err != nil {
				return fmt.Errorf("connect: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Clearing chat %s...\n", target.String())
			if err := a.WA().ClearChat(ctx, target); err != nil {
				return fmt.Errorf("clear chat: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "✅ Chat cleared successfully.")
			return nil
		},
	}
	cmd.Flags().StringVar(&jid, "jid", "", "chat JID to clear (e.g. 120363371857899816@g.us)")
	_ = cmd.MarkFlagRequired("jid")
	return cmd
}
