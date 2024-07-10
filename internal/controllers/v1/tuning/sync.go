package tuning

import (
	"fmt"

	"github.com/bigstack-oss/cube-api/internal/cubeos"
	definition "github.com/bigstack-oss/cube-api/internal/definition/v1"
	log "go-micro.dev/v5/logger"
)

func (c *Controller) syncByDesiredAction(tuning definition.Tuning) error {
	switch tuning.Status.Desired {
	case definition.Create, definition.Update:
		return c.applyTuning(tuning)
	case definition.Delete:
		return c.deleteTuning(tuning)
	}

	return fmt.Errorf(
		"unknown desired action(%s) for tuning(%s)",
		tuning.Status.Desired,
		tuning.Name,
	)
}

func (c *Controller) deleteTuning(tuning definition.Tuning) error {
	etcTunings, err := cubeos.GetEtcPoliciesTunings()
	if err != nil {
		log.Errorf("Failed to get all tunings: %s", err.Error())
		return err
	}

	etcTunings.RemoveTuning(tuning.Name)
	err = cubeos.HexTuningConfigure(etcTunings.Tunings)
	if err != nil {
		log.Errorf("Failed to delete tunings: %s", err.Error())
		return err
	}

	err = cubeos.IsHexTuningDeleted(tuning)
	if err != nil {
		log.Errorf("Failed to check if tuning %s is deleted: %s", tuning.Name, err.Error())
		return err
	}

	return nil
}

func (c *Controller) applyTuning(tuning definition.Tuning) error {
	etcTunings, err := cubeos.GetEtcPoliciesTunings()
	if err != nil {
		return err
	}

	etcTunings.AppendTunings([]definition.Tuning{tuning})
	err = cubeos.HexTuningConfigure(etcTunings.Tunings)
	if err != nil {
		log.Errorf("Failed to apply tuning %s: %s", tuning.Name, err.Error())
		return err
	}

	err = cubeos.IsHexTuningApplied(tuning)
	if err != nil {
		log.Errorf("Failed to check if tuning %s is applied: %s", tuning.Name, err.Error())
		return err
	}

	return nil
}
