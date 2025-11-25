package middleware

import (
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/haju35/Task_manager_API_Auth/data"
    "github.com/haju35/Task_manager_API_Auth/models"
)

// Claims contains custom JWT claims
type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func getSecret() []byte {
    if s := os.Getenv("JWT_SECRET"); s != "" {
        return []byte(s)
    }
    return []byte("replace_with_secure_secret")
}

// GenerateToken helper to create JWT tokens for a user
func GenerateToken(u *models.User, ttl time.Duration) (string, error) {
    claims := &Claims{
        UserID:   u.ID,
        Username: u.Username,
        Role:     u.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(getSecret())
}

// AuthMiddleware validates JWT and sets current user in context
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if auth == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }
        parts := strings.SplitN(auth, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
            return
        }

        tokenStr := parts[1]
        token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
            return getSecret(), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            return
        }

        // fetch user from data layer (optional, for safety)
        user, err := data.GetUserByID(claims.UserID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
            return
        }

        c.Set("currentUser", user)
        c.Next()
    }
}

// RequireRole middleware ensures the current user has one of allowed roles
func RequireRole(roles ...string) gin.HandlerFunc {
    allowed := map[string]bool{}
    for _, r := range roles {
        allowed[r] = true
    }
    return func(c *gin.Context) {
        cuI, exists := c.Get("currentUser")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not authenticated"})
            return
        }
        cu := cuI.(*models.User)
        if !allowed[cu.Role] {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient privileges"})
            return
        }
        c.Next()
    }
}

// GetCurrentUser helper
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
    cuI, ok := c.Get("currentUser")
    if !ok {
        return nil, false
    }
    cu, ok := cuI.(*models.User)
    return cu, ok
}

// TokenFromUser convenience used in controllers
func TokenFromUser(u *models.User) (string, error) {
    ttl := time.Hour * 24
    if v := os.Getenv("JWT_TTL_HOURS"); v != "" {
        if n, err := strconv.Atoi(v); err == nil && n > 0 {
            ttl = time.Duration(n) * time.Hour
        }
    }
    return GenerateToken(u, ttl)
}
