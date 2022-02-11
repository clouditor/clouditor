package gorm

import (
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/persistence"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageOptions(t *testing.T) {
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

func Test_storage_Get(t *testing.T) {
	var (
		err  error
		s    persistence.Storage
		user *auth.User
	)

	user = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}

	// Create storage
	s, err = NewStorage()
	assert.ErrorIs(t, err, nil)

	// Return error since no record in the DB yet
	err = s.Get(&auth.User{})
	assert.ErrorIs(t, err, persistence.ErrRecordNotFound)

	// Create user
	err = s.Create(user)
	assert.ErrorIs(t, err, nil)

	// Get user via passing entire record
	gotUser := &auth.User{}
	err = s.Get(gotUser)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, user, gotUser)

	// Get user via username
	gotUser2 := &auth.User{}
	err = s.Get(gotUser2, "username = ?", user.Username)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, user, gotUser2)

	// Get user via mail
	gotUser3 := &auth.User{}
	err = s.Get(gotUser3, "Email = ?", user.Email)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, user, gotUser3)

}

// TODO(lebogg): Add tests for List

func Test_storage_Count(t *testing.T) {
	var (
		count int64
		err   error
		s     persistence.Storage
		user  *auth.User
		user2 *auth.User
	)

	user = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}

	user2 = &auth.User{
		Username: "SomeName2",
		Password: "SomePassword2",
		Email:    "SomeMail2",
		FullName: "SomeFullName2",
	}

	// Create storage
	s, err = NewStorage()
	assert.ErrorIs(t, err, nil)

	// Since no records in DB yet, count of users should be 0
	count, err = s.Count(&auth.User{})
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, int(count), 0)

	// Create one user -> count of users should be 1
	assert.ErrorIs(t, s.Create(user), nil)
	count, err = s.Count(&auth.User{})
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, int(count), 1)

	// Add another one -> count of users should be 2
	assert.ErrorIs(t, s.Create(user2), nil)
	count, err = s.Count(&auth.User{})
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, int(count), 2)

	// Count of users with ID "SomeName2" should be 1
	count, err = s.Count(&auth.User{}, "Id = ?", "SomeName2")
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, int(count), 1)

	// Calling s.Count() with unsupported record type should throw "unsupported" error
	_, err = s.Count(nil)
	assert.NotNil(t, err)
	fmt.Println(err)
	assert.Contains(t, err.Error(), "unsupported data type")
}

// TODO(lebogg): Fix after I fixed the missing variadic ellipsis
func Test_storage_Update(t *testing.T) {
	var (
		err  error
		s    persistence.Storage
		user *auth.User
	)
	user = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}

	// Create storage
	s, err = NewStorage()
	assert.ErrorIs(t, err, nil)

	// Create user
	err = s.Create(user)
	assert.ErrorIs(t, err, nil)
	err = s.Get(&auth.User{}, "username = ?", user.Username)

	err = s.Update(&auth.User{FullName: "SomeNewName"}, "username = ?", user.Username)
	gotUser := &auth.User{}
	err = s.Get(gotUser, "username = ?", user.Username)
	assert.ErrorIs(t, err, nil)
	fmt.Println(gotUser)
	// UserName should be changed
	assert.Equal(t, "SomeNewName", gotUser.FullName)
	// Other properties should stay the same
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Password, gotUser.Password)
	assert.Equal(t, user.Email, gotUser.Email)

}

// TODO(lebogg): Add tests for delete
