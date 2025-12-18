package handlers

import (
	"context"
	"encoding/json"

	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/entities"
	"github.com/paingphyoaungkhant/asto-microservice/services/course-service/internal/domain/repositories"
	"github.com/paingphyoaungkhant/asto-microservice/shared/events"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type ZoomMeetingCreatedHandler struct {
	moduleRepo repositories.SectionModuleRepository
	logger     *logger.Logger
}

func NewZoomMeetingCreatedHandler(
	moduleRepo repositories.SectionModuleRepository,
	logger *logger.Logger,
) *ZoomMeetingCreatedHandler {
	return &ZoomMeetingCreatedHandler{
		moduleRepo: moduleRepo,
		logger:     logger,
	}
}

func (h *ZoomMeetingCreatedHandler) Handle(body []byte) error {
	ctx := context.Background()
	
	var event events.ZoomMeetingCreatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("failed to unmarshal zoom meeting created event", zap.Error(err))
		return err
	}

	module, err := h.moduleRepo.FindByID(ctx, event.SectionModuleID)
	if err != nil {
		h.logger.Error("failed to find section module",
			zap.Error(err),
			zap.String("section_module_id", event.SectionModuleID),
		)
		return err
	}

	if module == nil {
		h.logger.Warn("section module not found",
			zap.String("section_module_id", event.SectionModuleID),
		)
		return nil
	}

	module.UpdateContent(&event.ZoomMeetingID, entities.ContentStatusCreated)
	if err := h.moduleRepo.Update(ctx, module); err != nil {
		h.logger.Error("failed to update section module with zoom meeting id",
			zap.Error(err),
			zap.String("section_module_id", event.SectionModuleID),
			zap.String("zoom_meeting_id", event.ZoomMeetingID),
		)
		return err
	}

	h.logger.Info("updated section module with zoom meeting id",
		zap.String("section_module_id", event.SectionModuleID),
		zap.String("zoom_meeting_id", event.ZoomMeetingID),
	)

	return nil
}

