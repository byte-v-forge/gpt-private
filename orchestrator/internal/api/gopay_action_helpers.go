//go:build private_plugins

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"orchestrator/internal/jobstatus"
	"orchestrator/pb"
)

func structMap(data proto.Message) map[string]any {
	if data == nil {
		return map[string]any{}
	}
	if st, ok := data.(*structpb.Struct); ok {
		return st.AsMap()
	}
	raw, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(data)
	if err != nil {
		return map[string]any{}
	}
	out := map[string]any{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{}
	}
	return out
}

func structData(_ map[string]any) *pb.JobData {
	return nil
}

func mergeActionData(dst map[string]any, key string, value map[string]any) {
	if dst == nil || strings.TrimSpace(key) == "" || value == nil {
		return
	}
	dst[key] = value
}

func (s *Server) markActionFailed(ctx context.Context, jobID string, step string, status string, recoverable bool, retryable bool, err error, _ map[string]any) error {
	if err == nil {
		return nil
	}
	if status = strings.TrimSpace(status); status == "" {
		status = jobstatus.Failed(recoverable, retryable)
	}
	return s.activities.MarkJobFailedActivity(ctx, pb.JobFailureInput{
		JobId:        strings.TrimSpace(jobID),
		StepName:     strings.TrimSpace(step),
		Status:       status,
		Recoverable:  recoverable,
		Retryable:    retryable,
		ErrorMessage: err.Error(),
	})
}

func (req rawN8NActionRequest) DataString(key string) string {
	if req.Data == nil {
		return ""
	}
	value, ok := req.Data[strings.TrimSpace(key)]
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func firstNonZeroInt32(values ...int32) int32 {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}
