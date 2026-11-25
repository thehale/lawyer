// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

const batchSize = 256

type reading struct {
	licensable Licensable
	err        error
}

func readBatch(repository Repository, files []Path) []reading {
	var readings []reading
	for _, path := range files {
		licensable, isLicensable, err := licensableAt(repository, path)
		if err != nil || isLicensable {
			readings = append(readings, reading{licensable: licensable, err: err})
		}
	}
	return dateHistories(repository, readings)
}

func dateHistories(repository Repository, readings []reading) []reading {
	var files []Path
	for _, reading := range readings {
		if reading.err == nil {
			files = append(files, reading.licensable.path)
		}
	}
	repository.history = repository.history.Date(files)
	for index := range readings {
		readings[index].licensable.repository = repository
	}
	return readings
}

func yieldEach(readings []reading, yield func(Licensable, error) bool) bool {
	proceed := true
	for index := 0; proceed && index < len(readings); index++ {
		proceed = yield(readings[index].licensable, readings[index].err)
	}
	return proceed
}
