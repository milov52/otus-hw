package server

import (
	"github.com/milov52/hw12_13_14_15_calendar/internal/model"
	desc "github.com/milov52/hw12_13_14_15_calendar/pkg/api/event/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EventFromReq(event *desc.EventInfo) (*model.Event, error) {
	return &model.Event{
		Title:        event.GetTitle(),
		StartTime:    event.GetStartTime().AsTime(),
		Duration:     event.GetDuration().AsDuration(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserId(),
		NotifyBefore: event.GetNotifyBefore().AsDuration(),
		Sent:         event.GetSent(),
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
			NotifyBefore: durationpb.New(e.NotifyBefore),
			Sent:         e.Sent,
		},
	}
}

func EventsToResp(es []model.Event) *desc.GetResponse {
	resp := &desc.GetResponse{}
	for _, e := range es {
		resp.Events = append(resp.Events, EventToResp(e))
	}
	return resp
}
