package mailpit

import (
	"context"
	"fmt"
	"journey/internal/pgstore"
	"time"

	"github.com/gofrs/uuid"
	"github.com/wneessen/go-mail"
)


type store interface {
	GetTrip(context.Context, uuid.UUID) (pgstore.Trip, error)
}

type Mailpit struct {
	store store
}

func NewMailPit() Mailpit {
	return Mailpit{}
}

func (mp Mailpit) SendConfirmTripEmailToTripOwner(tripID uuid.UUID ) error {
	ctx := context.Background()
	trip, err := mp.store.GetTrip(ctx, tripID)

	if err != nil {
		return fmt.Errorf("mailpit: failed to get trip for SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg := mail.NewMsg()
	if err := msg.From("mailpit@journey.com"); err != nil {
		return fmt.Errorf("mailpit: failed to set From in email SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := msg.To(trip.OwnerEmail); err != nil {
		return fmt.Errorf("mailpit: failed to set To in email SendConfirmTripEmailToTripOwner: %w", err)
	}

	msg.Subject("Confirm your trip")

	msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
		Olá, %s!

		A sua viagem para o lugar %s que começa no dia %s precisa ser confirmada.
		Clique no botão abaixo para confirmar
	`, trip.OwnerName, trip.Destination, trip.StartsAt.Time.Format(time.DateOnly)),
	)

	client, err := mail.NewClient("localhost", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))

	if err != nil {
		return fmt.Errorf("mailpit: failed to create email client SendConfirmTripEmailToTripOwner: %w", err)
	}

	if err := client.DialAndSend(msg); err != nil {
		return fmt.Errorf("mailpit: failed to send email SendConfirmTripEmailToTripOwner: %w", err)
	}

	return nil
}