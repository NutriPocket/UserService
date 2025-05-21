package e2e_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/repository"
	"github.com/NutriPocket/UserService/service"
	"github.com/NutriPocket/UserService/test"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	jwtService, err := service.NewJWTService(nil)
	if err != nil {
		log.Fatalf("An error ocurred when creating the JWT service: %v\n", err)
	}

	testUser := model.User{
		Username: "test", Email: "test@test.com",
	}

	token, err := jwtService.Sign(testUser)

	if err != nil {
		log.Fatalf("An error ocurred when signing testUser: %v\n", err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	t.Run("It should retrieve an unauthorized status code if not auth header is provided", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Unauthorized user", data["title"])
		assert.Equal(t, "The user isn't authorized because no Authorization header is provided", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/users/", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve an empty array if the table is empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}
		assert.Empty(t, data, "If the users table is empty, it should return an empty array")
	})

	t.Run("It should retrieve all the users in the table", func(t *testing.T) {
		defer test.ClearUsers()
		w := httptest.NewRecorder()
		repository, err := repository.NewUserRepository(nil)
		if err != nil {
			t.Errorf("An error ocurred when creating the user repository: %v\n", err)
		}

		repository.CreateUser(&model.BaseUser{Username: "test1", Email: "test1@test.com", Password: "test1"})
		repository.CreateUser(&model.BaseUser{Username: "test2", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "test3", Email: "test3@test.com", Password: "test3"})

		req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}

		assert.NotEmpty(t, data, "If the users table is not empty, it should return a non-empty array")
		assert.Len(t, data, 3, "The length of the array should be 3")
		assert.Equal(t, "test3", data[0].Username)
		assert.Equal(t, "test2", data[1].Username)
		assert.Equal(t, "test1", data[2].Username)
		assert.Equal(t, "test3@test.com", data[0].Email)
		assert.Equal(t, "test2@test.com", data[1].Email)
		assert.Equal(t, "test1@test.com", data[2].Email)
	})

	t.Run("It should retrieve all the users in the table if searchUsername is an empty string", func(t *testing.T) {
		defer test.ClearUsers()
		w := httptest.NewRecorder()
		repository, err := repository.NewUserRepository(nil)
		if err != nil {
			t.Errorf("An error ocurred when creating the user repository: %v\n", err)
		}

		repository.CreateUser(&model.BaseUser{Username: "test1", Email: "test1@test.com", Password: "test1"})
		repository.CreateUser(&model.BaseUser{Username: "test2", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "test3", Email: "test3@test.com", Password: "test3"})

		req, _ := http.NewRequest(http.MethodGet, "/users/?searchUsername=", nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}

		assert.NotEmpty(t, data, "If the users table is not empty, it should return a non-empty array")
		assert.Len(t, data, 3, "The length of the array should be 3")
		assert.Equal(t, "test3", data[0].Username)
		assert.Equal(t, "test2", data[1].Username)
		assert.Equal(t, "test1", data[2].Username)
		assert.Equal(t, "test3@test.com", data[0].Email)
		assert.Equal(t, "test2@test.com", data[1].Email)
		assert.Equal(t, "test1@test.com", data[2].Email)
	})

	t.Run("It should retrieve only the users in the table that matchs searchUsername param partially", func(t *testing.T) {
		defer test.ClearUsers()
		w := httptest.NewRecorder()
		repository, err := repository.NewUserRepository(nil)
		if err != nil {
			t.Errorf("An error ocurred when creating the user repository: %v\n", err)
		}

		repository.CreateUser(&model.BaseUser{Username: "test1", Email: "test1@test.com", Password: "test1"})
		repository.CreateUser(&model.BaseUser{Username: "jorge", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "pedro", Email: "test3@test.com", Password: "test3"})

		req, _ := http.NewRequest(http.MethodGet, "/users/?searchUsername=o", nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}

		assert.NotEmpty(t, data, "If the users table is not empty, it should return a non-empty array")
		assert.Len(t, data, 2, "The length of the array should be 2")
		assert.Equal(t, "pedro", data[0].Username)
		assert.Equal(t, "jorge", data[1].Username)
		assert.Equal(t, "test3@test.com", data[0].Email)
		assert.Equal(t, "test2@test.com", data[1].Email)
	})

	t.Run("It should retrieve only the users in the table that matchs searchUsername param in the string order", func(t *testing.T) {
		defer test.ClearUsers()
		w := httptest.NewRecorder()
		repository, err := repository.NewUserRepository(nil)
		if err != nil {
			t.Errorf("An error ocurred when creating the user repository: %v\n", err)
		}

		repository.CreateUser(&model.BaseUser{Username: "egroj", Email: "test1@test.com", Password: "test1"})
		repository.CreateUser(&model.BaseUser{Username: "jorge", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "pedro", Email: "test3@test.com", Password: "test3"})

		req, _ := http.NewRequest(http.MethodGet, "/users/?searchUsername=jorge", nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}

		assert.NotEmpty(t, data, "If the users table is not empty, it should return a non-empty array")
		assert.Len(t, data, 1, "The length of the array should be 1")
		assert.Equal(t, "jorge", data[0].Username)
		assert.Equal(t, "test2@test.com", data[0].Email)
	})
}

func TestGetUser(t *testing.T) {
	jwtService, err := service.NewJWTService(nil)
	if err != nil {
		log.Fatalf("An error ocurred when creating the JWT service: %v\n", err)
	}

	testUser := model.User{
		Username: "test", Email: "test@test.com",
	}

	token, err := jwtService.Sign(testUser)

	if err != nil {
		log.Fatalf("An error ocurred when signing testUser: %v\n", err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	t.Run("It should retrieve an unauthorized status code if not auth header is provided", func(t *testing.T) {
		username := "test"
		req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Unauthorized user", data["title"])
		assert.Equal(t, "The user isn't authorized because no Authorization header is provided", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/users/"+username, data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a not found status if the user with username 'test' is not found", func(t *testing.T) {
		username := "test"

		req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "User not found", data["title"])
		assert.Equal(t, "The user with the username "+username+" was not found", data["detail"])
		assert.Equal(t, float64(404), data["status"])
		assert.Equal(t, "/users/"+username, data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve only the user with username 'test'", func(t *testing.T) {
		username := "test"
		w := httptest.NewRecorder()

		repository, err := repository.NewUserRepository(nil)
		if err != nil {
			t.Errorf("An error ocurred when creating the user repository: %v\n", err)
		}

		repository.CreateUser(&model.BaseUser{Username: testUser.Username, Email: testUser.Email, Password: "test"})
		repository.CreateUser(&model.BaseUser{Username: "test2", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "test3", Email: "test3@test.com", Password: "test3"})
		defer test.ClearUsers()

		req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data model.User
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a model.User parseable string, ", err)
		}

		assert.Equal(t, testUser.Username, data.Username)
		assert.Equal(t, testUser.Email, data.Email)
	})
}
