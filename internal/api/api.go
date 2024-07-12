package api

import (
	"context"
	"errors"
	"journey/internal/api/spec"
	"journey/internal/pgstore"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type store interface {
	GetParticipant(ctx context.Context, participantID uuid.UUID) (pgstore.Participant, error)
	ConfirmParticipant(ctx context.Context, participantID uuid.UUID) error
}

type API struct {
	store store
	logger *zap.Logger
}

// Confirms a participant on a trip.
// (PATCH /participants/{participantId}/confirm)
func (api API) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter,  r *http.Request, participantID string) *spec.Response {
	id, err := uuid.Parse(participantID)
	if err != nil {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "uuid invalido"})
	}

	participant, err := api.store.GetParticipant(r.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "participante não encontrado"})
		}
		api.logger.Error("failed to get participant", zap.Error(err), zap.String("participant_id", participantID))
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "something went wrong, try again"})
	}

	if participant.IsConfirmed {
		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(
			spec.Error{Message: "participante já confirmado"},
		)
	}

	if err := api.store.ConfirmParticipant(r.Context(), id); err != nil {
		api.logger.Error("failed to confirm participant", zap.Error(err), zap.String("participant_id", participantID))

		return spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "something went wrong, try again"})
	}

	return spec.PatchParticipantsParticipantIDConfirmJSON204Response(nil)
}

// Creates a new trip.
// (POST /trips)
func (api API) PostTrips(w http.ResponseWriter, r *http.Request) *spec.Response {
	panic("not implemented") //TODO: Implement
}

// Get a trip details.
// (GET /trips/{tripId})
func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripId string) *spec.Response {
	panic("not implemented") //TODO: Implement
}

// Update a trip
// (PUT /trips/{tripId})
func (api API) PutTripsTripID(w http.ResponseWriter, r *http.Request, tripId string) *spec.Response {
	panic("not implemented") //TODO: Implement
}

// Get a trip activities.
// (GET /trips/{tripId}/activities)
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripdID string) *spec.Response {
	panic("not implemented") //TODO: Implement
}

// Create a new trip activity.
// (POST /trips/{tripId}/activities)
func (api API) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	panic("not implemented") //TODO: Implement
}