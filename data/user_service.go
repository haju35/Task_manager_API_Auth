package data

import (
    "errors"
    "sync"
    "time"

    "github.com/haju35/Task_manager_API_Auth/models"
    "golang.org/x/crypto/bcrypt"
)

var (
    users      = map[int]*models.User{}
    usersMutex sync.RWMutex
    nextUserID = 1
)

// CreateUser registers a new user
func CreateUser(username, plainPassword string) (*models.User, error) {
    usersMutex.Lock()
    defer usersMutex.Unlock()

    // unique username
    for _, u := range users {
        if u.Username == username {
            return nil, errors.New("username already exists")
        }
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    role := "user"
    if len(users) == 0 {
        role = "admin" // first user becomes admin
    }

    u := &models.User{
        ID:       nextUserID,
        Username: username,
        Password: string(hashed),
        Role:     role,
    }
    nextUserID++
    users[u.ID] = u

    _ = time.Now() // no-op placeholder

    return u, nil
}

// Authenticate checks username/password
func Authenticate(username, plainPassword string) (*models.User, error) {
    usersMutex.RLock()
    defer usersMutex.RUnlock()

    for _, u := range users {
        if u.Username == username {
            if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)); err != nil {
                return nil, errors.New("invalid credentials")
            }
            return u, nil
        }
    }
    return nil, errors.New("user not found")
}

// GetUserByID returns a copy of user
func GetUserByID(id int) (*models.User, error) {
    usersMutex.RLock()
    defer usersMutex.RUnlock()

    if u, ok := users[id]; ok {
        copy := *u
        return &copy, nil
    }
    return nil, errors.New("user not found")
}

// PromoteToAdmin sets user role to admin
func PromoteToAdmin(id int) (*models.User, error) {
    usersMutex.Lock()
    defer usersMutex.Unlock()

    if u, ok := users[id]; ok {
        u.Role = "admin"
        copy := *u
        return &copy, nil
    }
    return nil, errors.New("user not found")
}

// ListUsers returns all users
func ListUsers() []*models.User {
    usersMutex.RLock()
    defer usersMutex.RUnlock()

    out := make([]*models.User, 0, len(users))
    for _, u := range users {
        copy := *u
        out = append(out, &copy)
    }
    return out
}
