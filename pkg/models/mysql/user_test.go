package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/jcorry/morellis/pkg/models"
)

func TestUserModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:        1,
				FirstName: "Alice",
				LastName:  "Jones",
				Email:     "alice@example.com",
				Phone:     "867-5309",
				Status:    "verified",
				Created:   time.Date(2019, 02, 24, 17, 25, 25, 0, time.UTC),
			},
			wantError: nil,
		},
	}

	db, teardown := newTestDB(t)
	defer teardown()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := UserModel{db}

			user, err := m.Get(tt.userID)

			// No way to generate this...has to come from DB
			tt.wantUser.UUID = user.UUID

			if err != tt.wantError {
				t.Errorf("want %v, got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v, got %v", tt.wantUser, user)
			}
		})
	}
}

func TestUserModel_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	result, err := m.Delete(1)
	if !result {
		t.Errorf("want true, got false")
	}

	if err != nil {
		t.Errorf("want nil err, got %s", err)
	}
}

func TestUserModel_List(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	m.Delete(1)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	users := []*models.User{
		{
			FirstName: "John",
			LastName:  "Corry",
			Email:     "jcorry@morellis.com",
			Phone:     "867-5309",
			Password:  string(hashedPassword),
		},
		{
			FirstName: "Garrett",
			LastName:  "Rap",
			Email:     "garrett@morellis.com",
			Phone:     "867-5309",
			Password:  string(hashedPassword),
		},
		{
			FirstName: "Brian",
			LastName:  "Morton",
			Email:     "brian@morellis.com",
			Phone:     "867-5309",
			Password:  string(hashedPassword),
		},
	}

	toD := []int64{}

	for _, u := range users {
		uid, err := uuid.NewRandom()
		if err != nil {
			t.Error("Failed to create UUID for user")
		}

		user, err := m.Insert(uid, u.FirstName, u.LastName, u.Email, u.Phone, u.Password)
		if err != nil {
			t.Fatal("Failed to insert new user for test")
		}
		user.UUID = uid
		toD = append(toD, user.ID)
		time.Sleep(time.Second)
	}

	tests := []struct {
		name          string
		order         string
		limit         int
		wantUserNames []string
		wantError     error
	}{
		{
			"Get all 3 users, no order",
			"created",
			3,
			[]string{"John", "Garrett", "Brian"},
			nil,
		},
		{
			"Get all 3 users, order by first name",
			"firstName",
			3,
			[]string{"Brian", "Garrett", "John"},
			nil,
		},
		{
			"Get 2 users, order by first name",
			"firstName",
			2,
			[]string{"Brian", "Garrett"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := m.List(tt.limit, 0, tt.order)
			if err != tt.wantError {
				t.Errorf("want nil, got %s", err)
			}

			if len(list) != tt.limit {
				t.Errorf("want list length %d; got %d", tt.limit, len(list))
			}

			for i, u := range list {
				if u.FirstName != tt.wantUserNames[i] {
					t.Errorf("want %s; got %s", tt.wantUserNames[i], u.FirstName)
				}
			}

		})
	}

}

func TestUserModel_GetByUUID(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	u, err := m.Get(1)
	if err != nil {
		t.Error(err)
	}

	u, err = m.GetByUUID(u.UUID)
	if err != nil {
		t.Error(err)
	}
}

func TestUserModel_GetByCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	tests := []struct {
		name      string
		email     string
		password  string
		wantError error
	}{
		{
			name:      "Found user",
			email:     "alice@example.com",
			password:  "password",
			wantError: nil,
		},
		{
			name:      "Wrong email",
			email:     "bob@example.com",
			password:  "password",
			wantError: models.ErrInvalidCredentials,
		},
		{
			name:      "Wrong password",
			email:     "alice@example.com",
			password:  "P@SSW0rd",
			wantError: models.ErrInvalidCredentials,
		},
	}

	db, teardown := newTestDB(t)
	defer teardown()

	m := UserModel{db}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := models.Credentials{
				Email:    tt.email,
				Password: tt.password,
			}

			_, err := m.GetByCredentials(c)
			if err != tt.wantError {
				t.Error(err)
			}
		})
	}
}
