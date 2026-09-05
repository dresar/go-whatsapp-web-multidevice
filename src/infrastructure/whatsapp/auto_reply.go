package whatsapp

import (
	"context"
	"strings"

	domainChatStorage "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/chatstorage"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func handleAutoReply(ctx context.Context, evt *events.Message, chatStorageRepo domainChatStorage.IChatStorageRepository, client *whatsmeow.Client) {
	if client == nil {
		return
	}

	// Skip groups, broadcasts, and self messages
	if utils.IsGroupJID(evt.Info.Chat.String()) || evt.Info.IsIncomingBroadcast() || evt.Info.IsFromMe {
		return
	}

	// Only reply to direct 1:1 chats (e.g., *@s.whatsapp.net)
	if evt.Info.Chat.Server != types.DefaultUserServer {
		return
	}

	// Extra safety: skip any broadcast/status contexts
	source := evt.Info.SourceString()
	if strings.Contains(source, "broadcast") ||
		strings.HasSuffix(evt.Info.Chat.String(), "@broadcast") ||
		strings.HasPrefix(evt.Info.Chat.String(), "status@") {
		return
	}

	// Require actual typed text (not captions or synthetic labels)
	hasText := false
	var content string

	// Unwrap FutureProof wrappers to access the inner message content first
	innerMsg := evt.Message
	for i := 0; i < 3; i++ { // safeguard against excessively nested wrappers
		if vm := innerMsg.GetViewOnceMessage(); vm != nil && vm.GetMessage() != nil {
			innerMsg = vm.GetMessage()
			continue
		}
		if em := innerMsg.GetEphemeralMessage(); em != nil && em.GetMessage() != nil {
			innerMsg = em.GetMessage()
			continue
		}
		if vm2 := innerMsg.GetViewOnceMessageV2(); vm2 != nil && vm2.GetMessage() != nil {
			innerMsg = vm2.GetMessage()
			continue
		}
		if vm2e := innerMsg.GetViewOnceMessageV2Extension(); vm2e != nil && vm2e.GetMessage() != nil {
			innerMsg = vm2e.GetMessage()
			continue
		}
		break
	}

	// Check for genuine typed text on the unwrapped content
	if conv := innerMsg.GetConversation(); conv != "" {
		hasText = true
		content = conv
	} else if ext := innerMsg.GetExtendedTextMessage(); ext != nil && ext.GetText() != "" {
		hasText = true
		content = ext.GetText()
	} else if protoMsg := innerMsg.GetProtocolMessage(); protoMsg != nil {
		if edited := protoMsg.GetEditedMessage(); edited != nil {
			if ext := edited.GetExtendedTextMessage(); ext != nil && ext.GetText() != "" {
				hasText = true
				content = ext.GetText()
			} else if conv := edited.GetConversation(); conv != "" {
				hasText = true
				content = conv
			}
		}
	}
	if !hasText {
		return
	}

	if chatStorageRepo == nil {
		return
	}

	// Fetch auto replies from DB
	replies, err := chatStorageRepo.GetAutoReplies()
	if err != nil {
		log.Errorf("Failed to fetch auto replies: %v", err)
		return
	}

	if len(replies) == 0 {
		return
	}

	messageContent := strings.TrimSpace(content)
	var responseText string

	// Match logic
	for _, reply := range replies {
		if reply.MatchType == "Exact" {
			if strings.EqualFold(messageContent, reply.Keyword) {
				responseText = reply.Response
				break
			}
		} else if reply.MatchType == "Contains" {
			if strings.Contains(strings.ToLower(messageContent), strings.ToLower(reply.Keyword)) {
				responseText = reply.Response
				break
			}
		}
	}

	if responseText == "" {
		return
	}

	// Format recipient JID
	recipientJID := utils.FormatJID(evt.Info.Sender.String())

	// Send the auto-reply message
	response, err := client.SendMessage(
		ctx,
		recipientJID,
		&waE2E.Message{Conversation: proto.String(responseText)},
	)

	if err != nil {
		log.Errorf("Failed to send auto-reply message: %v", err)
		return
	}

	// Store the auto-reply message in chat storage if send was successful
	// Get our own JID as sender
	senderJID := ""
	if client.Store.ID != nil {
		senderJID = client.Store.ID.String()
	}

	// Store the sent auto-reply message
	if err := chatStorageRepo.StoreSentMessageWithContext(
		ctx,
		response.ID,           // Message ID from WhatsApp response
		senderJID,             // Our JID as sender
		recipientJID.String(), // Recipient JID
		responseText,          // Auto-reply content
		response.Timestamp,    // Timestamp from response
	); err != nil {
		// Log storage error but don't fail the auto-reply
		log.Errorf("Failed to store auto-reply message in chat storage: %v", err)
	} else {
		log.Debugf("Auto-reply message %s stored successfully in chat storage", response.ID)
	}
}
