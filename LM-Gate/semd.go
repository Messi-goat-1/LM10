package lmgate

import (
	"encoding/json"
)

func SerializeReport(report *AnalysisReport) ([]byte, error) {
	return json.Marshal(report)
}
