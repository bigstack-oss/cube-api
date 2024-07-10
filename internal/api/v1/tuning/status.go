package tuning

import definition "github.com/bigstack-oss/cube-api/internal/definition/v1"

func setBatchPendingDeletion(tunings []definition.Tuning) {
	for i := range tunings {
		tunings[i].SetDesiredToDelete()
		tunings[i].SetCurrentToPending()
	}
}

func setBatchPendingUpdate(tunings []definition.Tuning) {
	for i := range tunings {
		tunings[i].SetDesiredToUpdate()
		tunings[i].SetCurrentToPending()
	}
}
