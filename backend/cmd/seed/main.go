package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"inventory-management-system/backend/internal/domain/entities"
	"inventory-management-system/backend/internal/infrastructure/config"
	"inventory-management-system/backend/internal/infrastructure/db"
	"inventory-management-system/backend/internal/infrastructure/security"
)

type seedUser struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx := context.Background()
	client, collections, err := db.Connect(ctx, cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
	}()

	users := []seedUser{
		{Name: "System Admin", Email: "admin@inventory.local", Password: "Admin@12345", Role: entities.RoleAdmin},
		{Name: "System Staff", Email: "staff@inventory.local", Password: "Staff@12345", Role: entities.RoleStaff},
	}

	for _, u := range users {
		if err = upsertUser(ctx, collections.Users, u); err != nil {
			log.Fatalf("seed user %s failed: %v", u.Email, err)
		}
		fmt.Printf("seeded user: %s (%s)\n", u.Email, u.Role)
	}

	fmt.Println("seed completed")
}

func upsertUser(ctx context.Context, usersColl interface {
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}, user seedUser) error {
	hash, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":         user.Name,
			"email":        user.Email,
			"passwordHash": hash,
			"role":         user.Role,
			"updatedAt":    now,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}

	_, err = usersColl.UpdateOne(
		ctx,
		bson.M{"email": user.Email},
		update,
		options.Update().SetUpsert(true),
	)
	return err
}
