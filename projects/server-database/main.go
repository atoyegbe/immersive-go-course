package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	// "log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	pgx "github.com/jackc/pgx/v5"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}


func isValidURL(inputURL string) bool {
	_, err := url.ParseRequestURI(inputURL)
	return err == nil
}

func fetchesImages(conn *pgx.Conn) ([]Image, error) {
	var images []Image
	rows, err := conn.Query(context.Background(), "select title, alt_text, url from images")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	for rows.Next() {
		var image Image
		if err := rows.Scan(&image.Title, &image.AltText, &image.URL); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		images = append(images, image)
	}
	return images, nil
}

func uploadImages(conn *pgx.Conn, img Image) (Image, error) {
	_, err := conn.Exec(context.Background(), "insert into public.images(title, url, alt_text) values($1, $2, $3)", img.Title, img.URL, img.AltText)
	if err != nil {
		return Image{}, err
	}
	return img, nil
}

func main() {

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		DATABASE_URL = "postgresql://postgres:test@localhost:5432/go-server-database"
	}
	con_bkg := context.Background()
	db, err := pgx.Connect(con_bkg, DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	}

	defer db.Close(con_bkg)

	// uncomment the code below to create the database
	// statements := []string{
	// 	`CREATE TABLE public.images (
	// 		id serial NOT NULL,
	// 		title text NOT NULL,
	// 		url text NOT NULL,
	// 		alt_text text,
	// 		PRIMARY KEY (id)
	// 	);`,

	// 	`ALTER TABLE IF EXISTS public.images OWNER to postgres;`,

	// 	`INSERT INTO public.images(title, url, alt_text)
	// 		VALUES (
	// 			'Sunset',
	// 			'https://images.unsplash.com/photo-1506815444479-bfdb1e96c566?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80',
	// 			'Clouds at sunset')`,
	// 	`INSERT INTO public.images(title, url, alt_text)
	// 	VALUES (
	// 		'Mountain',
	// 		'https://images.unsplash.com/photo-1540979388789-6cee28a1cdc9?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1000&q=80',
	// 		'A mountain at sunset')`,
	// }

	// for _, stmt := range statements {
	// 	_, err := db.Exec(context.Background(), stmt)
	// 	if err != nil {
	// 		log.Fatalf("Error executing SQL statement: %v\nSQL: %s", err, stmt)
	// 	}
	// 	fmt.Printf("Executed: %s\n", stmt)
	// }
	// fmt.Println("Table created successfully!")

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		indentValue := queryValues.Get("indent")
		res, err := strconv.Atoi(indentValue)
		if err != nil {
			w.Write([]byte("Invalid indent value"))
		}

		images, _ := fetchesImages(db)
		resp, err := json.MarshalIndent(images, "", strings.Repeat(" ", res))
		if err != nil {
			w.Write([]byte("Error while fetch json data for images"))
		}
		w.Write([]byte(resp))
	})

	http.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error", 500)
			}
			defer r.Body.Close()

			var image Image
			error := json.Unmarshal(body, &image)
			if error != nil {
				http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
				return
			}
			if isValidURL(image.URL) {
				img, err := uploadImages(db, image)
				if err != nil {
					http.Error(w, "Error uploading image", http.StatusInternalServerError)
					return
				}
				// Write the image URL as the response
				w.Write([]byte(img.URL))
			} else {
				http.Error(w, "Invalid image URL", http.StatusBadRequest)
				return
			}
		}
	})

	http.ListenAndServe(":8000", nil)
}
