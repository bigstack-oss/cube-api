package tuning

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	definition "github.com/bigstack-oss/cube-api/internal/definition/v1"
	cubeHttp "github.com/bigstack-oss/cube-api/internal/helpers/http"
	"github.com/bigstack-oss/cube-api/internal/service"
	log "go-micro.dev/v5/logger"
)

func delegateTuningsReq(tunings []definition.Tuning) {
	for _, tuning := range tunings {
		if definition.DoseCurrentRoleShouldHandleTheTuning(tuning.Name, definition.CurrentRole) {
			delegateToCurrentNode(tuning)
		}

		delegateToOtherNodes(tuning)
	}
}

func delegateToCurrentNode(tuning definition.Tuning) {
	syncTuningRecord(tuning)
	reqQueue.Add(tuning)
}

func delegateToOtherNodes(tuning definition.Tuning) {
	roles, found := definition.GetRolesToHandleTuning(tuning.Name)
	if !found {
		log.Warnf("no roles to handle tuning(%s)", tuning.Name)
		return
	}

	for _, role := range roles {
		nodes, err := service.GetNodesByRole(role.Name)
		if err == nil {
			sendTuningToOtherNodes(tuning, nodes)
			continue
		}

		if errors.Is(err, definition.ErrServiceNotFound) {
			continue
		}
		log.Errorf(
			"Failed to get nodes by role(%s): %s",
			role,
			err.Error(),
		)
	}
}

func sendTuningToOtherNodes(tuning definition.Tuning, nodes []definition.Node) {
	for _, node := range nodes {
		tuningReq, err := genTuningReq(node, tuning)
		if err != nil {
			log.Errorf("failed to create request for tuning %s (%s)", tuning.Name, err.Error())
			continue
		}

		resp, code := cubeHttp.NewHelper().Send(tuningReq)
		if cubeHttp.Is2XXCode[code] {
			continue
		}

		log.Errorf(
			"Failed to send tuning %s to node %s: %d %s",
			tuning.Name,
			node.ID,
			code,
			string(resp),
		)
	}
}

func genTuningReq(node definition.Node, tuning definition.Tuning) (*http.Request, error) {
	u := url.URL{
		Scheme: node.Protocol,
		Host:   node.Address,
		Path:   fmt.Sprintf("/api/v1/tunings/%s", tuning.Name),
	}

	b, err := tuning.Bytes()
	if err != nil {
		log.Errorf("failed to encode tuning %s (%s)", tuning.Name, err.Error())
		return nil, err
	}

	return http.NewRequest(
		http.MethodPut,
		u.String(),
		bytes.NewReader(b),
	)
}
