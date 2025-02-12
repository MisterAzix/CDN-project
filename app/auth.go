package app

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type User struct {
    ID       bson.ObjectID `bson:"_id,omitempty"`
    Email    string `bson:"email"`
    Password string `bson:"password"`
}

var userCollection *mongo.Collection

func InitAuth() {
    ConnectDB()
    userCollection = client.Database("file_manager").Collection("users")

    // Create an index on the email field
    indexModel := mongo.IndexModel{
        Keys:    bson.M{"email": 1}, // index in ascending order
        Options: options.Index().SetUnique(true),
    }
    _, err := userCollection.Indexes().CreateOne(context.TODO(), indexModel)
    if err != nil {
        log.Fatal(err)
    }
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    user.Password, err = hashPassword(user.Password)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    result, err := userCollection.InsertOne(context.TODO(), user)
    if err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    // Log the inserted user ID
    fmt.Printf("Inserted user with ID: %v\n", result.InsertedID)

    w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds User
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var user User
    err = userCollection.FindOne(context.TODO(), bson.D{{Key: "email", Value: creds.Email}}).Decode(&user)
    if err == mongo.ErrNoDocuments {
        fmt.Printf("User not found with email: %s\n", creds.Email) // Log the email
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }
    if err != nil {
        fmt.Printf("Error: %v\n", err) // Log the error
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    if !checkPasswordHash(creds.Password, user.Password) {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
    }

    token := "dummy-token"

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   token,
        Expires: time.Now().Add(24 * time.Hour),
    })

    w.WriteHeader(http.StatusOK)
}