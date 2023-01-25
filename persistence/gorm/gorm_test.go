package gorm

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/voc"
)

var (
	MockMetricID   = "MyMetric"
	MockMetricName = "MyMetricName"
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
					WithPostgres("", 0, "", "", "", ""),
				},
			},
			wantDialectorType: "",
			wantErr: func(tt assert.TestingT, err error, i ...any) bool {
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
			// Check if user has all necessary fields
			assert.NoError(t, userInput.Validate())
			err = s.Create(userInput)
			assert.NoError(t, err)

			// Get record via DB call
			userOutput := &auth.User{}
			err = s.Get(userOutput, "Username = ?", "SomeName")
			assert.NoError(t, err)
			assert.NoError(t, userOutput.Validate())
			assert.Equal(t, userInput, userOutput)
		})
	}
}

func Test_storage_Create(t *testing.T) {
	var (
		err error
		s   persistence.Storage
	)

<<<<<<< HEAD
=======
	metric = &assessment.Metric{
		Id:    MockMetricID,
		Name:  MockMetricName,
		Range: &assessment.Range{Range: &assessment.Range_MinMax{MinMax: &assessment.MinMax{Min: 1, Max: 2}}},
	}
	// Check if metric has all necessary fields
	assert.NoError(t, metric.Validate())

>>>>>>> main
	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Prepare evidence for creation
	value, err := voc.ToStruct(voc.Compute{
		Resource: &voc.Resource{
			ID:           "1",
			Name:         "SomeResource",
			CreationTime: 0,
			Type:         []string{"Compute", "VM"},
			GeoLocation: voc.GeoLocation{
				Region: "EU",
			},
		},
		NetworkInterface: nil,
	})
	assert.NoError(t, err)
	e := &evidence.Evidence{
		Id:        "1",
		Timestamp: timestamppb.Now(),
		ServiceId: "1",
		ToolId:    "1",
		Raw:       "Some:Raw:Stuff",
		Resource:  value,
	}

	// Create evidence
	err = s.Create(e)
	assert.NoError(t, err)

	// Create evidence a second time should fail since it is stored already
	err = s.Create(e)
	assert.Error(t, err)

	// Prepare and create metric
	metric := &assessment.Metric{Id: "Test"}
	err = s.Create(metric)
	assert.NoError(t, err)

	// Create metric a second time should fail since it is stored already
	err = s.Create(metric)
	assert.Error(t, err)
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
	// Check if user has all necessary fields
	assert.NoError(t, user.Validate())

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
	assert.NoError(t, gotUser.Validate())
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
	assert.NoError(t, gotUser3.Validate())
	assert.Equal(t, user, gotUser3)

	var metric = &assessment.Metric{
		Id:    MockMetricID,
		Name:  MockMetricName,
		Range: &assessment.Range{Range: &assessment.Range_MinMax{MinMax: &assessment.MinMax{Min: 1, Max: 2}}},
	}
	// Check if metric has all necessary fields
	assert.NoError(t, metric.Validate())

	// Create metric
	err = s.Create(metric)
	assert.NoError(t, err)

	// Get metric via Id
	gotMetric := &assessment.Metric{}
	err = s.Get(gotMetric, "id = ?", MockMetricID)
	assert.NoError(t, err)
	assert.NoError(t, gotMetric.Validate())
	assert.Equal(t, metric, gotMetric)

	var impl = &assessment.MetricImplementation{
		MetricId:  MockMetricID,
		Code:      "TestCode",
		UpdatedAt: timestamppb.New(time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)),
	}
	// Check if impl has all necessary fields
	assert.NoError(t, impl.Validate())

	// Create metric implementation
	err = s.Create(impl)
	assert.NoError(t, err)

	// Get metric implementation via Id
	gotImpl := &assessment.MetricImplementation{}
	err = s.Get(gotImpl, "metric_id = ?", MockMetricID)
	assert.NoError(t, err)
	assert.NoError(t, gotImpl.Validate())
	assert.Equal(t, impl, gotImpl)

	// Prepare evidence for creation
	value, err := voc.ToStruct(voc.Compute{
		Resource: &voc.Resource{
			ID:           "1",
			Name:         "SomeResource",
			CreationTime: 0,
			Type:         []string{"Compute", "VM"},
			GeoLocation: voc.GeoLocation{
				Region: "EU",
			},
		},
		NetworkInterface: nil,
	})
	assert.NoError(t, err)
	e := &evidence.Evidence{
		Id:        "1",
		Timestamp: timestamppb.Now(),
		ServiceId: "1",
		ToolId:    "1",
		Raw:       "Some:Raw:Stuff",
		Resource:  value,
	}

	// Create evidence
	err = s.Create(e)
	assert.NoError(t, err)

	// Get evidence via Id
	gotEvidence := &evidence.Evidence{}
	err = s.Get(gotEvidence, "id = ?", e.Id)
	assert.NoError(t, err)
	assert.Equal(t, e, gotEvidence)
}

func Test_storage_List(t *testing.T) {
	var (
		err   error
		s     persistence.Storage
		user1 *auth.User
		user2 *auth.User
		users []auth.User
	)

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Test user

	user1 = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}
	// Check if user has all necessary fields
	assert.NoError(t, user1.Validate())

	user2 = &auth.User{
		Username: "SomeName2",
		Password: "SomePassword2",
		Email:    "SomeMail2",
		FullName: "SomeFullName2",
	}
	// Check if user has all necessary fields
	assert.NoError(t, user2.Validate())

	// List should return empty list since no users are in DB yet
	err = s.List(&users, "", true, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Empty(t, users)

	// List should return list of 2 users (user1 and user2)
	err = s.Create(user1)
	assert.NoError(t, err)
	err = s.Create(user2)
	assert.NoError(t, err)
	err = s.List(&users, "", true, 0, -1)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, len(users), 2)
	// We only check one user and assume the others are also correct
	assert.NoError(t, users[0].Validate())

	// Test with certificates (associations included via states)
	var (
		certificate1 *orchestrator.Certificate
		certificate2 *orchestrator.Certificate
		certificates []*orchestrator.Certificate
	)

	// List should return empty list since no certificates are in DB yet
	err = s.List(&certificates, "", true, 0, 0)
	assert.ErrorIs(t, err, nil)
	assert.Empty(t, certificates)

	// Create two certificates
	certificate1 = orchestratortest.NewCertificate()
	certificate1.Id = "0"
	certificate2 = orchestratortest.NewCertificate()
	certificate2.Id = "1"
	err = s.Create(certificate1)
	assert.NoError(t, err)
	err = s.Create(certificate2)
	assert.NoError(t, err)

	// List should return list of 2 certificates with associated states
	err = s.List(&certificates, "id", false, 0, 0)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, len(certificates), 2)
	// Check ordering
	assert.Equal(t, certificate2.Id, certificates[0].Id)
	// We only check one certificate and assume the others are also correct
	assert.NoError(t, certificates[0].Validate())

	fmt.Println(certificates)

	// Check if certificate with id "1" (certificate2) is in the list and if states are included (association)
	for i := range certificates {
		if certificates[i].Id == certificate2.Id {
			fmt.Println("Certificate:", certificates[i])
			assert.NotEmpty(t, certificates[i].States)
			return
		}
	}
	// If not, let the test fail
	assert.FailNow(t, "%s is not listed but should be.", certificate1.Id)

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
	// Check if user has all necessary fields
	assert.NoError(t, user.Validate())

	user2 = &auth.User{
		Username: "SomeName2",
		Password: "SomePassword2",
		Email:    "SomeMail2",
		FullName: "SomeFullName2",
	}
	// Check if user has all necessary fields
	assert.NoError(t, user2.Validate())

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
	count, err = s.Count(&auth.User{}, "username = ?", "SomeName2")
	assert.NoError(t, err)
	assert.Equal(t, int(count), 1)

	// Calling s.Count() with unsupported record type should throw "unsupported" error
	_, err = s.Count(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported data type")
}

func Test_storage_Save(t *testing.T) {
	type MyTest struct {
		ID   int
		Name string
	}

	var (
		err     error
		s       persistence.Storage
		user    *auth.User
		newUser *auth.User
		gotUser *auth.User
		myVar   MyTest
	)
	user = &auth.User{
		Username: "SomeName",
		Password: "SomePassword",
		Email:    "SomeMail",
		FullName: "SomeFullName",
	}
	// Check if user has all necessary fields
	assert.NoError(t, user.Validate())

	// Create storage
	s, err = NewStorage(WithAdditionalAutoMigration(&MyTest{}))
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
	// Check if user has all necessary fields
	assert.NoError(t, newUser.Validate())

	err = s.Save(newUser, "username = ?", user.Username)
	assert.NoError(t, err)

	gotUser = &auth.User{}
	err = s.Get(gotUser, "username = ?", user.Username)
	assert.NoError(t, err)
	assert.NoError(t, gotUser.Validate())

	// UserName and FullName should be the same
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Username, gotUser.Username)
	// PW should be changed
	assert.Equal(t, newUser.Password, gotUser.Password)
	// Email should be zero
	assert.Equal(t, "", gotUser.Email)

	// Save MyTest
	myVar = MyTest{ID: 1, Name: "Test"}

	err = s.Save(&myVar)
	assert.NoError(t, err)
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
	// Check if user has all necessary fields
	assert.NoError(t, user.Validate())

	// Create storage
	s, err = NewStorage()
	assert.NoError(t, err)

	// Testing user
	// Create user
	err = s.Create(user)
	assert.NoError(t, err)

	err = s.Get(&auth.User{}, "username = ?", user.Username)
	assert.NoError(t, err)

	err = s.Update(&auth.User{FullName: "SomeNewName"}, "username = ?", "SomeOtherUser")
	assert.ErrorIs(t, err, persistence.ErrRecordNotFound)

	err = s.Update(&auth.User{FullName: "SomeNewName"}, "username = ?", user.Username)
	assert.NoError(t, err)

	gotUser := &auth.User{}
	err = s.Get(gotUser, "username = ?", user.Username)
	assert.NoError(t, err)
	assert.NoError(t, gotUser.Validate())

	// UserName should be changed
	assert.Equal(t, "SomeNewName", gotUser.FullName)

	// Other properties should stay the same
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Password, gotUser.Password)
	assert.Equal(t, user.Email, gotUser.Email)

	// Testing cloud service
	// Create cloud service
	cloudService := orchestrator.CloudService{
		Id:          orchestratortest.MockServiceID,
		Name:        "SomeName",
		Description: "SomeDescription",
		ConfiguredMetrics: []*assessment.Metric{
			{
				Id:    MockMetricID,
				Name:  MockMetricName,
				Range: &assessment.Range{Range: &assessment.Range_MinMax{MinMax: &assessment.MinMax{Min: 1, Max: 2}}},
			},
		},
	}
	// Check if cloud service has all necessary fields
	assert.NoError(t, cloudService.Validate())
	err = s.Create(&cloudService)
	assert.NoError(t, err)

	err = s.Get(&orchestrator.CloudService{}, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	err = s.Update(&orchestrator.CloudService{Name: "SomeNewName", Description: ""}, "Id = ?", cloudService.Id)
	assert.NoError(t, err)

	gotCloudService := &orchestrator.CloudService{}
	err = s.Get(gotCloudService, "Id = ?", cloudService.Id)
	assert.NoError(t, err)
	assert.NoError(t, gotCloudService.Validate())

	// Name should be changed
	assert.Equal(t, "SomeNewName", gotCloudService.Name)
	// Other properties should stay the same
	assert.Equal(t, cloudService.Id, gotCloudService.Id)
	assert.Equal(t, cloudService.Description, gotCloudService.Description)
	assert.Equal(t, len(cloudService.ConfiguredMetrics), len(gotCloudService.ConfiguredMetrics))
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
	// Check if user has all necessary fields
	assert.NoError(t, user.Validate())

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
