package controllers

import (
	"context"
	"fmt"
	"go-rest-app/pkg/config"
	"go-rest-app/pkg/models"
	"go-rest-app/pkg/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cardsCollection *mongo.Collection = config.GetCollection("cardsCollections")
)

func GetCard(c *fiber.Ctx) error {
	fmt.Print("Get card endpoint hit")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	card := models.Card{}

	err := cardsCollection.FindOne(ctx, bson.M{"id": id}).Decode(&card)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.CardResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responses.CardResponse{Status: http.StatusOK, Message: "success", Data: card})
}
func CreateCard(c *fiber.Ctx) error {
	fmt.Print("Create card endpoint hit")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	card := models.Card{}

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.CardResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: err.Error()})
	}

	newCard := models.Card{
		ID:     primitive.NewObjectID(),
		Name:   card.Name,
		Number: card.Number,
	}

	_, err := cardsCollection.InsertOne(ctx, newCard)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.CardResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.CardResponse{Status: http.StatusCreated, Message: "success", Data: newCard})
}

func GetCards(c *fiber.Ctx) error {
	fmt.Print("Get cards endpoint hit")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cards := []models.Card{}

	cursor, err := cardsCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.CardResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: err.Error()})
	}
	if err = cursor.All(ctx, &cards); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.CardResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(responses.CardResponse{Status: http.StatusOK, Message: "success", Data: cards})
}

func DeleteCard(c *fiber.Ctx) error {
	fmt.Print("Delete card endpoint hit")

	ctx, Cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer Cancel()

	id, _ := primitive.ObjectIDFromHex(c.Params("id"))

	result, err := cardsCollection.DeleteOne(ctx, bson.M{"id": id})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.CardResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.CardResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(responses.CardResponse{Status: http.StatusCreated, Message: "success", Data: "Card is Successfully deleted"})

}

func UpdateCard(c *fiber.Ctx) error {
	fmt.Println("Update endpoint get hit")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	card := models.Card{}

	defer cancel()

	cardId, _ := primitive.ObjectIDFromHex(c.Params("id"))

	fmt.Println("Update endpoint get hit 1", cardId)

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.CardResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	updatedCard := bson.M{"name": card.Name, "number": card.Number}

	result, err := cardsCollection.UpdateOne(ctx, bson.M{"id": cardId}, bson.M{"$set": updatedCard})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.CardResponse{Status: http.StatusNotFound, Message: "err", Data: err.Error()})
	}

	if result.MatchedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.CardResponse{Status: http.StatusAccepted, Message: "Car not found", Data: result})
	}
	card.ID = cardId

	return c.Status(fiber.StatusCreated).JSON(responses.CardResponse{Status: http.StatusAccepted, Message: "Updated", Data: card})

}
