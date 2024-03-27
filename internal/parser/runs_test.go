package parser

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowRunsParse(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want *WorkflowRunsStatsSummary
	}{
		{
			name: "empty",
			args: args{
				file: "empty.json",
			},
			want: &WorkflowRunsStatsSummary{
				TotalRunsCount: 0,
				Conclusions: map[string]*WorkflowRunsConclusion{
					ConclusionSuccess: {
						RunsCount:    0,
						WorkflowRuns: []*WorkflowRun{},
					},
					ConclusionFailure: {
						RunsCount:    0,
						WorkflowRuns: []*WorkflowRun{},
					},
					ConclusionOthers: {
						RunsCount:    0,
						WorkflowRuns: []*WorkflowRun{},
					},
				},
			},
		},
		{
			name: "Success",
			args: args{
				file: "success.json",
			},
			want: &WorkflowRunsStatsSummary{
				TotalRunsCount: 2,
				Name:           "CI",
				Rate: Rate{
					SuccesRate:  1,
					FailureRate: 0,
					OthersRate:  0,
				},
				ExecutionDurationStats: executionDurationStats{
					Min: 20.0,
					Max: 40.0,
					Avg: 30.0,
					Std: 10,
					Med: 30.0,
				},
				Conclusions: map[string]*WorkflowRunsConclusion{
					ConclusionSuccess: {
						RunsCount: 2,
						WorkflowRuns: []*WorkflowRun{
							{
								ID:           10000,
								Status:       "completed",
								Conclusion:   "success",
								Actor:        "test-user",
								RunAttempt:   1,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10000/attempts/1",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10000/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10000/logs",
								RunStartedAt: timeParse("2024-01-01T00:00:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:00:20Z"),
								CreatedAt:    timeParse("2024-01-01T00:00:00Z"),
								Duration:     20.0,
							},
							{
								ID:           10001,
								Status:       "completed",
								Conclusion:   "success",
								Actor:        "test-user2",
								RunAttempt:   1,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10001/attempts/1",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10001/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10001/logs",
								RunStartedAt: timeParse("2024-01-01T00:01:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:01:40Z"),
								CreatedAt:    timeParse("2024-01-01T00:01:00Z"),
								Duration:     40.0,
							},
						},
					},
					ConclusionFailure: {
						RunsCount:    0,
						WorkflowRuns: []*WorkflowRun{},
					},
					ConclusionOthers: {
						RunsCount:    0,
						WorkflowRuns: []*WorkflowRun{},
					},
				},
			},
		},
		{
			name: "Failure and others",
			args: args{
				file: "failure-others.json",
			},
			want: &WorkflowRunsStatsSummary{
				TotalRunsCount: 2,
				Name:           "CI",
				Rate: Rate{
					SuccesRate:  0,
					FailureRate: 0.5,
					OthersRate:  0.5,
				},
				ExecutionDurationStats: executionDurationStats{
					Min: 0,
					Max: 0,
					Avg: 0,
					Std: 0,
					Med: 0,
				},
				Conclusions: map[string]*WorkflowRunsConclusion{
					ConclusionSuccess: {
						RunsCount:    0,
						WorkflowRuns: []*WorkflowRun{},
					},
					ConclusionFailure: {
						RunsCount: 1,
						WorkflowRuns: []*WorkflowRun{
							{
								ID:           10000,
								Status:       "completed",
								Conclusion:   "failure",
								Actor:        "test-user",
								RunAttempt:   1,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10000/attempts/1",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10000/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10000/logs",
								RunStartedAt: timeParse("2023-01-01T00:00:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:00:20Z"),
								CreatedAt:    timeParse("2023-01-01T00:00:00Z"),
								Duration:     3024000,
							},
						},
					},

					ConclusionOthers: {
						RunsCount: 1,
						WorkflowRuns: []*WorkflowRun{
							{
								ID:           10001,
								Status:       "completed",
								Conclusion:   "other",
								Actor:        "test-user",
								RunAttempt:   1,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10001/attempts/1",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10001/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10001/logs",
								RunStartedAt: timeParse("2024-01-01T00:00:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:00:20Z"),
								CreatedAt:    timeParse("2024-01-01T00:00:00Z"),
								Duration:     20.0,
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple conclusions",
			args: args{
				file: "multiple-conclusions.json",
			},
			want: &WorkflowRunsStatsSummary{
				TotalRunsCount: 3,
				Name:           "CI",
				Rate: Rate{
					SuccesRate:  0.3333333333333333,
					FailureRate: 0.3333333333333333,
					OthersRate:  0.3333333333333334,
				},
				ExecutionDurationStats: executionDurationStats{
					Min: 20.0,
					Max: 20.0,
					Avg: 20.0,
					Std: 0.0,
					Med: 20.0,
				},
				Conclusions: map[string]*WorkflowRunsConclusion{
					ConclusionSuccess: {
						RunsCount: 1,
						WorkflowRuns: []*WorkflowRun{
							{
								ID:           10000,
								Status:       "completed",
								Conclusion:   "success",
								Actor:        "test-user",
								RunAttempt:   2,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10000/attempts/2",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10000/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10000/logs",
								RunStartedAt: timeParse("2024-01-01T00:00:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:00:20Z"),
								CreatedAt:    timeParse("2024-01-01T00:00:00Z"),
								Duration:     20.0,
							},
						},
					},
					ConclusionFailure: {
						RunsCount: 1,
						WorkflowRuns: []*WorkflowRun{
							{
								ID:           10001,
								Status:       "completed",
								Conclusion:   "failure",
								Actor:        "test-user",
								RunAttempt:   1,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10001/attempts/1",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10001/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10001/logs",
								RunStartedAt: timeParse("2024-01-01T00:00:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:00:40Z"),
								CreatedAt:    timeParse("2024-01-01T00:00:00Z"),
								Duration:     40.0,
							},
						},
					},
					ConclusionOthers: {
						RunsCount: 1,
						WorkflowRuns: []*WorkflowRun{
							{
								ID:           10002,
								Status:       "completed",
								Conclusion:   "other",
								Actor:        "test-user",
								RunAttempt:   1,
								HTMLURL:      "https://github.com/owner/repos/actions/runs/10002/attempts/1",
								JobsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10002/jobs",
								LogsURL:      "https://api.github.com/repos/owner/repos/actions/runs/10002/logs",
								RunStartedAt: timeParse("2024-01-01T00:00:00Z"),
								UpdateAt:     timeParse("2024-01-01T00:00:30Z"),
								CreatedAt:    timeParse("2024-01-01T00:00:00Z"),
								Duration:     30.0,
							},
						},
					},
				},
			},
		},
	} // Add closing brace here

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := os.ReadFile("testdata/runs/" + tt.args.file)
			if err != nil {
				t.Fatal(err)
			}
			var wrs github.WorkflowRuns
			if err := json.Unmarshal(d, &wrs); err != nil {
				t.Fatal(err)
			}
			got := WorkflowRunsParse(wrs.WorkflowRuns)
			assert.Equal(t, tt.want.TotalRunsCount, got.TotalRunsCount)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Rate, got.Rate)
			assert.Equal(t, tt.want.ExecutionDurationStats, got.ExecutionDurationStats)

			for _, c := range []string{ConclusionSuccess, ConclusionFailure, ConclusionOthers} {
				assert.Equal(t, tt.want.Conclusions[c].RunsCount, got.Conclusions[c].RunsCount)
				for i, wr := range tt.want.Conclusions[c].WorkflowRuns {
					assert.Equal(t, wr.ID, got.Conclusions[c].WorkflowRuns[i].ID)
					assert.Equal(t, wr.Status, got.Conclusions[c].WorkflowRuns[i].Status)
					assert.Equal(t, wr.Conclusion, got.Conclusions[c].WorkflowRuns[i].Conclusion)
					assert.Equal(t, wr.Actor, got.Conclusions[c].WorkflowRuns[i].Actor)
					assert.Equal(t, wr.RunAttempt, got.Conclusions[c].WorkflowRuns[i].RunAttempt)
					assert.Equal(t, wr.HTMLURL, got.Conclusions[c].WorkflowRuns[i].HTMLURL)
					assert.Equal(t, wr.JobsURL, got.Conclusions[c].WorkflowRuns[i].JobsURL)
					assert.Equal(t, wr.LogsURL, got.Conclusions[c].WorkflowRuns[i].LogsURL)
					assert.True(t, wr.RunStartedAt.Equal(got.Conclusions[c].WorkflowRuns[i].RunStartedAt))
					assert.True(t, wr.UpdateAt.Equal(got.Conclusions[c].WorkflowRuns[i].UpdateAt))
					assert.True(t, wr.CreatedAt.Equal(got.Conclusions[c].WorkflowRuns[i].CreatedAt))
					assert.Equal(t, wr.Duration, got.Conclusions[c].WorkflowRuns[i].Duration)

				}
			}
		})
	}
}

func timeParse(s string) time.Time {
	t, _ := time.ParseInLocation(time.RFC3339, s, time.UTC)
	return t.Round(0)
}
