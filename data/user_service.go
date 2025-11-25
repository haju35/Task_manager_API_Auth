package data


import (
"errors"
"sync"
"Task_manager_API_Auth/models"
"time"


"golang.org/x/crypto/bcrypt"
)


var (
users = map[int]*models.User{}
usersMutex sync.RWMutex
nextUserID = 1
)


// CreateUser registers a new user with hashed password. If no users exist, first user becomes admin.
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
ID: nextUserID,
Username: username,
Password: string(hashed),
Role: role,
}
nextUserID++
users[u.ID] = u


// simple persistence simulation: (no-op) but could write timestamp/log
_ = time.Now()


return u, nil
}
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


// GetUserByID returns a copy of the user (or nil) by id.
func GetUserByID(id int) (*models.User, error) {
usersMutex.RLock()
defer usersMutex.RUnlock()
if u, ok := users[id]; ok {
// return a shallow copy
copy := *u
return &copy, nil
}
return nil, errors.New("user not found")
}


// PromoteToAdmin sets user's role to admin. Only called after authorization checks.
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


// ListUsers (convenience) for debugging.
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
users[u.ID] = u


// simple persistence simulation: (no-op) but could write timestamp/log
_ = time.Now()


return u, nil
}


// Authenticate checks username/password and returns user if OK.
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


// GetUserByID returns a copy of the user (or nil) by id.
func GetUserByID(id int) (*models.User, error) {
usersMutex.RLock()
defer usersMutex.RUnlock()
if u, ok := users[id]; ok {
// return a shallow copy
copy := *u
return &copy, nil
}
return nil, errors.New("user not found")
}


// PromoteToAdmin sets user's role to admin. Only called after authorization checks.
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


// ListUsers (convenience) for debugging.
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