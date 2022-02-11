package gorm

import (
	"clouditor.io/clouditor/api/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {
	type args struct {
		opts []StorageOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "in memory with option",
			args: args{
				opts: []StorageOption{
					WithInMemory(),
				},
			},
			wantErr: false,
		},
		{
			name:    "in memory without option",
			wantErr: false,
		},
		{
			name: "postgres with option - invalid port",
			args: args{
				opts: []StorageOption{
					WithPostgres("", 0),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStorage(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				assert.Contains(t, err.Error(), "invalid port")
				return
			}

			// Test to create new user and get it again with respective 'Create' and 'Get'
			// Create record via DB call
			userInput := &auth.User{
				Username: "SomeName",
				Password: "SomePassword",
				Email:    "SomeMail",
				FullName: "SomeFullName",
			}
			err = s.Create(userInput)
			assert.ErrorIs(t, err, nil)

			// Get record via DB call
			userOutput := &auth.User{}
			err = s.Get(userOutput, "Id = ?", "SomeName")
			assert.ErrorIs(t, err, nil)
			assert.Equal(t, userInput, userOutput)

		})
	}
}
