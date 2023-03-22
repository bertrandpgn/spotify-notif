package src

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

var FirestoreClient *firestore.Client

func CreateClient(ctx context.Context) {
	client, err := firestore.NewClient(ctx, EnvVars.FirebaseProjectID)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}

	FirestoreClient = client
}

// func test() {
// 	ctx := context.Background()
// 	client := CreateClient(ctx)

// 	defer client.Close()

// 	// Write
// 	_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
// 		"first":  "Alan",
// 		"middle": "Mathison",
// 		"last":   "Turing",
// 		"born":   1912,
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed adding aturing: %v", err)
// 	}

// 	// Read
// 	iter := client.Collection("users").Documents(ctx)
// 	for {
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("Failed to iterate: %v", err)
// 		}
// 		fmt.Println(doc.Data())
// 	}
// }
