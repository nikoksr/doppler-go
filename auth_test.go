package doppler

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nikoksr/doppler-go/pointer"
)

func TestAuthRevokeOptions_MarshalJSON(t *testing.T) {
	t.Parallel()

	type fields struct {
		Tokens []AuthToken
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Tokens: []AuthToken{},
			},
			want:    []byte("[]"),
			wantErr: false,
		},
		{
			name: "one",
			fields: fields{
				Tokens: []AuthToken{
					{
						Token: pointer.To("token"),
					},
				},
			},
			want:    []byte("[{\"token\":\"token\"}]"),
			wantErr: false,
		},
		{
			name: "two",
			fields: fields{
				Tokens: []AuthToken{
					{
						Token: pointer.To("token1"),
					},
					{
						Token: pointer.To("token2"),
					},
				},
			},
			want:    []byte("[{\"token\":\"token1\"},{\"token\":\"token2\"}]"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := &AuthRevokeOptions{
				Tokens: tt.fields.Tokens,
			}
			got, err := opts.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("MarshalJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
