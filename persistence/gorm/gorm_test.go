package gorm

import (
	"fmt"
	"testing"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"

	"github.com/stretchr/testify/assert"
)

func TestStorageOptions(t *testing.T) {
	type args struct {
		opts []StorageOption
	}
	tests := []struct {
		name              string
		args              args
		wantDialectorType string
		wantErr           assert.ErrorAssertionFunc
	}{
		{
			name: "in memory with option",
			args: args{
				opts: []StorageOption{
					WithInMemory(),
				},
			},
			wantDialectorType: "*sqlite.Dialector",
			wantErr:           nil,
		},
		{
			name:              "in memory without option",
			wantDialectorType: "*sqlite.Dialector",
			wantErr:           nil,
		},
		{
			name: "postgres with option - invalid port",
			args: args{
				opts: []StorageOption{
					WithPostgres("", 0, "", "", ""),
				},
			},
			wantDialectorType: "",
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Contains(t, err.Error(), "invalid port")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewStorage(tt.args.opts...)
			if (err != nil) && tt.wantErr != nil {
				tt.wantErr(t, err, tt.args.opts)
				return
			}

			gorm, ok := s.(*storage)
			assert.True(t, ok)
			assert.NotNil(t, gorm)

			assert.Equal(t, tt.wantDialectorType, fmt.Sprintf("%T", gorm.dialector))

			// Test to create new user and get it again with respective 'Create' and 'Get'
			// Create record via DB call
			userInput := &auth.User{
				Username: "SomeName",
				Password: "SomePassword",
				Email:    "SomeMail",
				FullName: "SomeFullName",
			}
			err = s.Create(userInput)
			assert.NoError(t, err)

			// Get record via DB call
			userOutput := &auth.User{}
			err = s.Get(userOutput, "Username = ?", "SomeName")
			assert.NoError(t, err)
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
	assert.NoError(t, err)

	// Return error since no record in the DB yet
	err = s.Get(&auth.User{})
	assert.ErrorIs(t, err, persistence.ErrRecordNotFound)

	// Create user
	err = s.Create(user)
	assert.NoError(t, err)

	// Get user via passing entire record
	gotUser := &auth.User{}
	err = s.Get(gotUser)
	assert.NoError(t, err)
	assert.Equal(t, user, gotUser)

	// Get user via username
	gotUser2 := &auth.User{}
	err = s.Get(gotUser2, "username = ?", user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user, gotUser2)

	// Get user via mail
	gotUser3 := &auth.User{}
	err = s.Get(gotUser3, "Email = ?", user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user, gotUser3)

}

func Test_storage_List(t *testing.T) {
	var (
		err   error
		s     persistence.Storage
		user1 *auth.User
		user2 *auth.User
		users []auth.User
	)

	user1 = &auth.User{
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
	assert.NoError(t, err)

	// List should return empty list since no users are in DB yet
	err = s.List(&users, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Empty(t, users)

	// List should return list of 2 users (user1 and user2)
	err = s.Create(user1)
	assert.NoError(t, err)
	err = s.Create(user2)
	assert.NoError(t, err)
	err = s.List(&users, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, len(users), 2)

	// Check if user with name "SomeName" (user1) is in the list
	for i := range users {
		if users[i].Username == user1.Username {
			return
		}
	}
	// If not, let the test fail
	assert.FailNow(t, "user1 is not listed but should be.")

}

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
	assert.NoError(t, err)

	// Since no records in DB yet, count of users should be 0
	count, err = s.Count(&auth.User{})
	assert.NoError(t, err)
	assert.Equal(t, int(count), 0)

	// Create one user -> count of users should be 1
	assert.ErrorIs(t, s.Create(user), nil)
	count, err = s.Count(&auth.User{})
	assert.NoError(t, err)
	assert.Equal(t, int(count), 1)

	// Add another one -> count of users should be 2
	assert.ErrorIs(t, s.Create(user2), nil)
	count, err = s.Count(&auth.User{})
	assert.NoError(t, err)
	assert.Equal(t, int(count), 2)

	// Count of users with ID "SomeName2" should be 1
	count, err = s.Count(&auth.User{}, "Id = ?", "SomeName2")
	assert.NoError(t, err)
	assert.Equal(t, int(count), 1)

	// Calling s.Count() with unsupported record type should throw "unsupported" error
	_, err = s.Count(nil)
	assert.Error(t, err)
	fmt.Println(err)
	assert.Contains(t, err.Error(), "unsupported data type")
}

func Test_storage_Save(t *testing.T) {
	var (
		err     error
		s       persistence.Storage
		user    *auth.User
		newUser *auth.User
		gotUser *auth.User
	)
	user = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Create user
	err = s.Create(user)
	assert.NoError(t, err)

	err = s.Get(&auth.User{}, "username = ?", user.Username)
	assert.NoError(t, err)

	// Save new User: Change PW and delete email. Username and FullName remain unchanged
	newUser = &auth.User{
		Username: user.Username,
		Password: "SomeNewPassword",
		Email:    "",
		FullName: user.FullName,
	}
	err = s.Save(newUser, "username = ?", user.Username)
	assert.NoError(t, err)

	gotUser = &auth.User{}
	err = s.Get(gotUser, "username = ?", user.Username)
	assert.NoError(t, err)

	// UserName and FullName should be the same
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Username, gotUser.Username)
	// PW should be changed
	assert.Equal(t, newUser.Password, gotUser.Password)
	// Email should be zero
	assert.Equal(t, "", gotUser.Email)

}

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
	assert.NoError(t, err)

	// Create user
	err = s.Create(user)
	assert.NoError(t, err)

	err = s.Get(&auth.User{}, "username = ?", user.Username)
	assert.NoError(t, err)

	err = s.Update(&auth.User{FullName: "SomeNewName"}, "username = ?", user.Username)
	assert.NoError(t, err)

	gotUser := &auth.User{}
	err = s.Get(gotUser, "username = ?", user.Username)
	assert.NoError(t, err)

	// UserName should be changed
	assert.Equal(t, "SomeNewName", gotUser.FullName)

	// Other properties should stay the same
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Password, gotUser.Password)
	assert.Equal(t, user.Email, gotUser.Email)

	// Testing cloud service (A table test now would be better, probably)

	// Create user
	cloudService := orchestrator.CloudService{
		Id:          "SomeId",
		Name:        "SomeName",
		Description: "SomeDescription",
	}
	err = s.Create(&cloudService)
	assert.NoError(t, err)

	err = s.Get(&orchestrator.CloudService{}, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	err = s.Update(&orchestrator.CloudService{Name: "SomeNewName", Description: ""}, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	gotCloudService := &orchestrator.CloudService{}
	err = s.Get(gotCloudService, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	// Name should be changed
	assert.Equal(t, "SomeNewName", gotCloudService.Name)
	// Other properties should stay the same
	assert.Equal(t, cloudService.Id, gotCloudService.Id)
	assert.Equal(t, cloudService.Description, gotCloudService.Description)
}

func Test_storage_Delete(t *testing.T) {
	var (
		err  error
		s    persistence.Storage
		user *auth.User
		//gotUser *auth.User
	)
	user = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Create user
	err = s.Create(user)
	assert.NoError(t, err)

	// Should return ErrRecordNotFound since there is no user "FakeUserName" in DB
	assert.ErrorIs(t, s.Delete(&auth.User{}, "username = ?", "FakeUserName"), persistence.ErrRecordNotFound)

	// Successful deletion
	assert.Nil(t, s.Delete(&auth.User{}, "username = ?", user.Username))
	// Check with s.Get that user is not in DB anymore
	assert.ErrorIs(t, s.Get(&auth.User{}, "username = ?", user.Username), persistence.ErrRecordNotFound)

	// Should return DB error since a non-supported type is passed (just a string instead of, e.g., &auth.User{})
	assert.Contains(t, s.Delete("Unsupported Type").Error(), "unsupported data type")

}
