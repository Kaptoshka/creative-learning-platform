package grpc

import (
	"time"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	assignmentsv1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/assignments/v1"

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

func targetToProto(t domain.AssignmentTarget) *assignmentsv1.AssignmentTarget {
	if t.GroupID != "" {
		return &assignmentsv1.AssignmentTarget{
			Target: &assignmentsv1.AssignmentTarget_GroupId{GroupId: t.GroupID},
		}
	}
	return &assignmentsv1.AssignmentTarget{
		Target: &assignmentsv1.AssignmentTarget_StudentId{StudentId: t.StudentID},
	}
}

func targetsToProto(targets []domain.AssignmentTarget) []*assignmentsv1.AssignmentTarget {
	result := make([]*assignmentsv1.AssignmentTarget, len(targets))
	for i, t := range targets {
		result[i] = targetToProto(t)
	}
	return result
}

func protoToTarget(t *assignmentsv1.AssignmentTarget) domain.AssignmentTarget {
	if t == nil {
		return domain.AssignmentTarget{}
	}
	switch v := t.GetTarget().(type) {
	case *assignmentsv1.AssignmentTarget_GroupId:
		return domain.AssignmentTarget{GroupID: v.GroupId}
	case *assignmentsv1.AssignmentTarget_StudentId:
		return domain.AssignmentTarget{StudentID: v.StudentId}
	}
	return domain.AssignmentTarget{}
}

func protoToTargets(targets []*assignmentsv1.AssignmentTarget) []domain.AssignmentTarget {
	result := make([]domain.AssignmentTarget, len(targets))
	for i, t := range targets {
		result[i] = protoToTarget(t)
	}
	return result
}

// --- AssignmentTemplate ---

func protoToTemplate(t *assignmentsv1.AssignmentTemplate) domain.AssignmentTemplate {
	if t == nil {
		return domain.AssignmentTemplate{}
	}
	return domain.AssignmentTemplate{
		ID:           t.GetId(),
		CreatorID:    t.GetCreatorId(),
		Title:        t.GetTitle(),
		Description:  t.GetDescription(),
		WidgetID:     t.GetWidgetId(),
		WidgetConfig: protoToMap(t.GetWidgetConfig()),
		DueDate:      protoToTime(t.GetDueDate()),
		CreatedAt:    protoToTimeVal(t.GetCreatedAt()),
		UpdatedAt:    protoToTimeVal(t.GetUpdatedAt()),
	}
}

func protoToTemplateLight(t *assignmentsv1.AssignmentTemplateLight) domain.AssignmentTemplateLight {
	if t == nil {
		return domain.AssignmentTemplateLight{}
	}
	return domain.AssignmentTemplateLight{
		ID:         t.GetId(),
		Title:      t.GetTitle(),
		WidgetType: t.GetWidgetType(),
		DueDate:    protoToTime(t.GetDueDate()),
	}
}

// --- Submission ---

func protoToSubmission(s *assignmentsv1.Submission) domain.Submission {
	if s == nil {
		return domain.Submission{}
	}
	sub := domain.Submission{
		ID:          s.GetId(),
		TemplateID:  s.GetTemplateId(),
		StudentID:   s.GetStudentId(),
		Status:      domain.SubmissionStatus(s.GetStatus()),
		StartedAt:   protoToTime(s.GetStartedAt()),
		SubmittedAt: protoToTime(s.GetSubmittedAt()),
	}
	if s.GetLatestVersion() != nil {
		v := protoToVersionLight(s.GetLatestVersion())
		sub.LatestVersion = &v
	}
	return sub
}

func protoToSubmissions(items []*assignmentsv1.Submission) []domain.Submission {
	result := make([]domain.Submission, len(items))
	for i, s := range items {
		result[i] = protoToSubmission(s)
	}
	return result
}

// --- SubmissionVersion ---

func protoToVersion(v *assignmentsv1.SubmissionVersion) domain.SubmissionVersion {
	if v == nil {
		return domain.SubmissionVersion{}
	}
	return domain.SubmissionVersion{
		ID:            v.GetId(),
		VersionNumber: v.GetVersionNumber(),
		Payload:       protoToMap(v.GetPayload()),
		TimeSpent:     protoToDuration(v.GetTimeSpentSeconds()),
		IsAutosave:    v.GetIsAutosave(),
		CreatedAt:     protoToTimeVal(v.GetCreatedAt()),
		UpdatedAt:     protoToTimeVal(v.GetUpdatedAt()),
	}
}

func protoToVersionLight(v *assignmentsv1.SubmissionVersionLight) domain.SubmissionVersionLight {
	if v == nil {
		return domain.SubmissionVersionLight{}
	}
	return domain.SubmissionVersionLight{
		ID:            v.GetId(),
		VersionNumber: v.GetVersionNumber(),
		CreatedAt:     protoToTimeVal(v.GetCreatedAt()),
	}
}

func protoToVersions(items []*assignmentsv1.SubmissionVersion) []domain.SubmissionVersion {
	result := make([]domain.SubmissionVersion, len(items))
	for i, v := range items {
		result[i] = protoToVersion(v)
	}
	return result
}

// --- Feedback ---

func protoToFeedback(f *assignmentsv1.Feedback) domain.Feedback {
	if f == nil {
		return domain.Feedback{}
	}
	return domain.Feedback{
		ID:          f.GetId(),
		VersionID:   f.GetVersionId(),
		GraderID:    f.GetGraderId(),
		TextContent: f.GetTextContent(),
		Payload:     protoToMap(f.GetPayload()),
		IsPublished: f.GetIsPublished(),
		CreatedAt:   protoToTimeVal(f.GetCreatedAt()),
		UpdatedAt:   protoToTimeVal(f.GetUpdatedAt()),
	}
}

func protoToFeedbacks(items []*assignmentsv1.Feedback) []domain.Feedback {
	result := make([]domain.Feedback, len(items))
	for i, f := range items {
		result[i] = protoToFeedback(f)
	}
	return result
}

// --- SubmissionStatus ---

func statusToProto(s domain.SubmissionStatus) assignmentsv1.SubmissionStatus {
	return assignmentsv1.SubmissionStatus(s)
}
