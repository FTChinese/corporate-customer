package api

import (
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/brianvoe/gofakeit/v5"
	"testing"
)

func TestClient_ReaderSignup(t *testing.T) {

	faker.SeedGoFake()

	client := MockNewClient()

	type args struct {
		s input.SignupParams
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "User signup",
			args: args{
				s: input.SignupParams{
					Credentials: input.Credentials{
						Email:    gofakeit.Email(),
						Password: faker.SimplePassword(),
					},
				},
			},
			want:    200,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.ReaderSignup(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReaderSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("Response should not be empty")
				return
			}

			if got.StatusCode != 200 {
				t.Errorf("ReaderSignUp() response code got = %d, want %d", got.StatusCode, tt.want)
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ReaderSignup() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}
