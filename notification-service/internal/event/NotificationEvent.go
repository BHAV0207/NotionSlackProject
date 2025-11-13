package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BHAV0207/notification-service/internal/handler"
	"github.com/BHAV0207/notification-service/internal/service"
)

func (c *Consumer) StartConsuming() {
	fmt.Printf("📩 [%s] Kafka consumer started on topic: %s\n", c.ServiceName, c.Kafka.Reader.Config().Topic)

	brokerURL := c.Kafka.Reader.Config().Brokers[0]
	dlqTopic := c.Kafka.Reader.Config().Topic + "-dlq"

	for {
		msg, err := c.Kafka.Reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("❌ Kafka read error:", err)
			continue
		}

		var event GenericEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Println("⚠️ Invalid Kafka event:", err)
			SendToDlq(brokerURL, dlqTopic, msg.Value)
			continue
		}

		// Retry logic
		const maxRetries = 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = c.ProcessEvent(event)
			if err == nil {
				break // ✅ Success
			}
			fmt.Printf("⚠️ Attempt %d/%d failed for event %+v: %v\n", attempt, maxRetries, event, err)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			fmt.Printf("💀 Event failed after %d attempts, sending to DLQ...\n", maxRetries)
			SendToDlq(brokerURL, dlqTopic, event)
		}
	}
}

func (c *Consumer) ProcessEvent(event GenericEvent) error {
	userMsg := buildMessage(event)
	if userMsg == "" {
		return fmt.Errorf("empty notification message for event: %+v", event)
	}

	// Try sending notification
	c.TrySendNotification(event, userMsg)

	// Store notification in DB
	query := `
		INSERT INTO notifications (user_id, type, message, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := c.Collection.Exec(context.Background(), query,
		event.UserID, event.EventType, userMsg, "SENT", time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("DB insert failed: %v", err)
	}
	return nil
}

// Helper for message content
func buildMessage(event GenericEvent) string {
	switch event.EventType {
	case "user-created":
		display := event.Name
		if display == "" {
			display = event.UserID
		}
		return fmt.Sprintf("👋 Welcome aboard, %s!", display)

	case "user-logged-in": // ✅ fixed name
		display := event.Name
		if display == "" {
			display = event.UserID
		}
		if display == "" {
			display = event.Email
		}
		if display == "" {
			display = "there"
		}
		return fmt.Sprintf("👋 Welcome back, %s!", display)

	case "payment-success":
		return fmt.Sprintf("💰 Payment for order #%s succeeded!", event.OrderID)

	case "payment-failed":
		return fmt.Sprintf("⚠️ Payment for order #%s failed. Please retry.", event.OrderID)

	default:
		return fmt.Sprintf("🔔 Update on your order #%s", event.OrderID)
	}
}

func (c *Consumer) TrySendNotification(event GenericEvent, message string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("💥 Panic recovered in sendNotification: %v\n", r)
		}
	}()

	c.sendNotification(event, message)
	return nil
}

func (c *Consumer) sendNotification(event GenericEvent, message string) {
	var (
		user *service.UserResponse
		err  error
	)

	if event.UserID != "" {
		user, err = service.GetUserByID(event.UserID)
		if err != nil {
			fmt.Printf("❌ Failed to fetch user %s: %v\n", event.UserID, err)
		}
	}

	email := event.Email
	name := event.Name

	if user != nil {
		if email == "" {
			email = user.Email
		}
		if name == "" {
			name = user.Name
		}
	}

	if email == "" {
		fmt.Printf("⚠️ Skipping notification: no email available for event %+v\n", event)
		return
	}

	if name == "" {
		name = "there"
	}

	subject := "Notification from E-com Website"
	body := fmt.Sprintf("Hey %s,<br><br>%s<br><br>– Team E-com", name, message)

	if err := handler.SendEmail(email, subject, body); err != nil {
		fmt.Printf("❌ Failed to send email to %s: %v\n", email, err)
		return
	}

	fmt.Printf("✅ Notification email sent to %s\n", email)
}
