package grpc

import (
	"time"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/domain"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- time.Time <-> timestamppb ---

func timeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func protoToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func protoToTimeVal(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// --- time.Duration <-> durationpb ---

func durationToProto(d time.Duration) *durationpb.Duration {
	return durationpb.New(d)
}

func protoToDuration(d *durationpb.Duration) time.Duration {
	if d == nil {
		return 0
	}
	return d.AsDuration()
}

// --- map[string]any <-> structpb.Struct ---

func mapToProto(m map[string]any) *structpb.Struct {
	if m == nil {
		return nil
	}
	s, _ := structpb.NewStruct(m)
	return s
}

func protoToMap(s *structpb.Struct) map[string]any {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// --- AssignmentTarget ---

func targetToProto(t domain.AssignmentTarget) *tasksv1.AssignmentTarget {
	if t.GroupID != "" {
		return &tasksv1.AssignmentTarget{
			Target: &tasksv1.AssignmentTarget_GroupId{GroupId: t.GroupID},
		}
	}
	return &tasksv1.AssignmentTarget{
		Target: &tasksv1.AssignmentTarget_StudentId{StudentId: t.StudentID},
	}
}

func targetsToProto(targets []domain.AssignmentTarget) []*tasksv1.AssignmentTarget {
	result := make([]*tasksv1.AssignmentTarget, len(targets))
	for i, t := range targets {
		result[i] = targetToProto(t)
	}
	return result
}

func protoToTarget(t *tasksv1.AssignmentTarget) domain.AssignmentTarget {
	if t == nil {
		return domain.AssignmentTarget{}
	}
	switch v := t.Target.(type) {
	case *tasksv1.AssignmentTarget_GroupId:
		return domain.AssignmentTarget{GroupID: v.GroupId}
	case *tasksv1.AssignmentTarget_StudentId:
		return domain.AssignmentTarget{StudentID: v.StudentId}
	}
	return domain.AssignmentTarget{}
}

func protoToTargets(targets []*tasksv1.AssignmentTarget) []domain.AssignmentTarget {
	result := make([]domain.AssignmentTarget, len(targets))
	for i, t := range targets {
		result[i] = protoToTarget(t)
	}
	return result
}

// --- AssignmentTemplate ---

func protoToTemplate(t *tasksv1.AssignmentTemplate) domain.AssignmentTemplate {
	if t == nil {
		return domain.AssignmentTemplate{}
	}
	return domain.AssignmentTemplate{
		ID:           t.Id,
		CreatorID:    t.CreatorId,
		Title:        t.Title,
		Description:  t.Description,
		WidgetID:     t.WidgetId,
		WidgetConfig: protoToMap(t.WidgetConfig),
		DueDate:      protoToTime(t.DueDate),
		CreatedAt:    protoToTimeVal(t.CreatedAt),
		UpdatedAt:    protoToTimeVal(t.UpdatedAt),
	}
}

func protoToTemplateLight(t *tasksv1.AssignmentTemplateLight) domain.AssignmentTemplateLight {
	if t == nil {
		return domain.AssignmentTemplateLight{}
	}
	return domain.AssignmentTemplateLight{
		ID:         t.Id,
		Title:      t.Title,
		WidgetType: t.WidgetType,
		DueDate:    protoToTime(t.DueDate),
	}
}

// --- Submission ---

func protoToSubmission(s *tasksv1.Submission) domain.Submission {
	if s == nil {
		return domain.Submission{}
	}
	sub := domain.Submission{
		ID:          s.Id,
		TemplateID:  s.TemplateId,
		StudentID:   s.StudentId,
		Status:      domain.SubmissionStatus(s.Status),
		StartedAt:   protoToTime(s.StartedAt),
		SubmittedAt: protoToTime(s.SubmittedAt),
	}
	if s.LatestVersion != nil {
		v := protoToVersion(s.LatestVersion)
		sub.LatestVersion = &v
	}
	return sub
}

func protoToSubmissions(items []*tasksv1.Submission) []domain.Submission {
	result := make([]domain.Submission, len(items))
	for i, s := range items {
		result[i] = protoToSubmission(s)
	}
	return result
}

// --- SubmissionVersion ---

func protoToVersion(v *tasksv1.SubmissionVersion) domain.SubmissionVersion {
	if v == nil {
		return domain.SubmissionVersion{}
	}
	return domain.SubmissionVersion{
		ID:            v.Id,
		VersionNumber: v.VersionNumber,
		Payload:       protoToMap(v.Payload),
		TimeSpent:     protoToDuration(v.TimeSpentSeconds),
		IsAutosave:    v.IsAutosave,
		CreatedAt:     protoToTimeVal(v.CreatedAt),
		UpdatedAt:     protoToTimeVal(v.UpdatedAt),
	}
}

func protoToVersions(items []*tasksv1.SubmissionVersion) []domain.SubmissionVersion {
	result := make([]domain.SubmissionVersion, len(items))
	for i, v := range items {
		result[i] = protoToVersion(v)
	}
	return result
}

// --- Feedback ---

func protoToFeedback(f *tasksv1.Feedback) domain.Feedback {
	if f == nil {
		return domain.Feedback{}
	}
	return domain.Feedback{
		ID:          f.Id,
		VersionID:   f.VersionId,
		GraderID:    f.GraderId,
		TextContent: f.TextContent,
		Payload:     protoToMap(f.Payload),
		IsPublished: f.IsPublished,
		CreatedAt:   protoToTimeVal(f.CreatedAt),
		UpdatedAt:   protoToTimeVal(f.UpdatedAt),
	}
}

func protoToFeedbacks(items []*tasksv1.Feedback) []domain.Feedback {
	result := make([]domain.Feedback, len(items))
	for i, f := range items {
		result[i] = protoToFeedback(f)
	}
	return result
}

// --- SubmissionStatus ---

func statusToProto(s domain.SubmissionStatus) tasksv1.SubmissionStatus {
	return tasksv1.SubmissionStatus(s)
}
