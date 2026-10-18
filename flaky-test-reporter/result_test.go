package main

import (
	"reflect"
	"testing"

	"knative.dev/toolbox/pkg/junit"
	"knative.dev/toolbox/pkg/prow"
)

func Test_includeFinishedBuild(t *testing.T) {
	startTime := int64(1)
	finishTime := int64(2)

	tests := []struct {
		name  string
		build prow.Build
		want  bool
	}{
		{
			name:  "unfinished build",
			build: prow.Build{},
			want:  false,
		},
		{
			name: "finished build without start time",
			build: prow.Build{
				StoragePath: "logs/job/123",
				FinishTime:  &finishTime,
			},
			want: false,
		},
		{
			name: "finished build with start time",
			build: prow.Build{
				StoragePath: "logs/job/124",
				StartTime:   &startTime,
				FinishTime:  &finishTime,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := includeFinishedBuild(&tt.build); got != tt.want {
				t.Fatalf("includeFinishedBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterOutParentTests(t *testing.T) {
	tests := []struct {
		name          string
		originalCases []junit.TestCase
		want          []junit.TestCase
	}{
		{
			name: "no parent cases",
			originalCases: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "ccc",
			}},
			want: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "ccc",
			}},
		},
		{
			name: "one parent case",
			originalCases: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "bbb/ccc",
			}},
			want: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb/ccc",
			}},
		},
		{
			name: "two nested cases",
			originalCases: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb",
			}, {
				Name: "bbb/ccc",
			}, {
				Name: "bbb/ddd",
			}, {
				Name: "bbb/ddd/fff",
			}},
			want: []junit.TestCase{{
				Name: "aaa",
			}, {
				Name: "bbb/ccc",
			}, {
				Name: "bbb/ddd/fff",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterOutParentTests(tt.originalCases); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterOutParentTests() = %v, want %v", got, tt.want)
			}
		})
	}
}
