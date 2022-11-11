package doppler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToParameters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		v       interface{}
		want    parameters
		wantErr bool
	}{
		{
			name:    "nil",
			v:       nil,
			want:    make(parameters),
			wantErr: false,
		},
		{
			name:    "empty",
			v:       struct{}{},
			want:    make(parameters),
			wantErr: false,
		},
		{
			name: "project.ListOptions",
			v: &ProjectListOptions{
				ListOptions: ListOptions{
					Page:    1,
					PerPage: 2,
				},
			},
			want: parameters{
				"page":     []string{"1"},
				"per_page": []string{"2"},
			},
			wantErr: false,
		},
		{
			name: "project.ProjectGetOptions",
			v: &ProjectGetOptions{
				Name: "123",
			},
			want: parameters{
				"project": []string{"123"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractQueryParameters(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractQueryParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("extractQueryParameters() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
