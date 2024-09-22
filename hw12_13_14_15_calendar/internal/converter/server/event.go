package server

import (
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	desc "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventFromReq(event *desc.EventInfo) (*model.Event, error) {
	notification := &model.Notification{
		EventID: event.GetNotification().GetEventId(),
		Title:   event.GetNotification().GetTitle(),
		Date:    event.GetNotification().GetDate().AsTime(),
		UserID:  event.GetNotification().GetUserId(),
	}

	return &model.Event{
		Title:        event.GetTitle(),
		StartTime:    event.GetStartTime().AsTime(),
		Duration:     event.GetDuration().AsDuration(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserId(),
		Notification: notification,
	}, nil
}

func EventToResp(e model.Event) *desc.Event {
	return &desc.Event{
		Id: e.ID.String(),
		Event: &desc.EventInfo{
			Title:        e.Title,
			StartTime:    timestamppb.New(e.StartTime),
			Duration:     durationpb.New(e.Duration),
			Description:  e.Description,
			UserId:       e.UserID,
			Notification: nil,
		},
	}
}

func EventsToResp(es []model.Event) *desc.GetResponse {
	resp := &desc.GetResponse{}
	for _, e := range es {
		resp.Event = append(resp.Event, EventToResp(e))
	}
	return resp
}
