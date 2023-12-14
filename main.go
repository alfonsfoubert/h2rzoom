package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// Main struct representing the entire JSON object.
type MessageData struct {
	Messages []Message `json:"messages"`
}

// Message struct for each message in the array.
type Message struct {
	ID            string         `json:"id"`
	Timestamp     int64          `json:"timestamp"` // UNIX timestamp is best represented as int64
	Snippet       MessageSnippet `json:"snippet"`
	AuthorDetails Author         `json:"authorDetails"`
	Platform      PlatformDetail `json:"platform"`
}

// MessageSnippet struct for the snippet part of each message.
type MessageSnippet struct {
	DisplayMessage string `json:"displayMessage"`
}

// Author struct for author details.
type Author struct {
	DisplayName     string `json:"displayName"`
	ProfileImageUrl string `json:"profileImageUrl"`
}

// PlatformDetail struct for the platform details.
type PlatformDetail struct {
	Name    string `json:"name"`
	LogoUrl string `json:"logoUrl"`
}

func getProfileURL(sender string) string {
	switch sender {
	case "Alfons Foubert":
		return "https://alfonsfoubertcom.files.wordpress.com/2022/03/img_0405.jpg"
	default:
		return "https://alfonsfoubertcom.files.wordpress.com/2023/12/screenshot-2023-12-13-at-10.24.02.png"
	}
}

// Convert string in the format "HH:MM:SS" to Unix timestamp
func TimeStringToUnix(timeStr string) (int64, error) {
	// Assuming the current date
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	fullTimeStr := fmt.Sprintf("%d-%02d-%02d %s", now.Year(), now.Month(), now.Day(), timeStr)

	parsedTime, err := time.Parse(layout, fullTimeStr)
	if err != nil {
		return 0, err
	}

	return parsedTime.Unix(), nil
}

// ParseMessages parses messages from a reader.
func ParseMessages(reader *bufio.Reader, lastTime int64) ([]Message, error) {
	// Assumes the format "HH:MM:SS From [Sender's Name] to ..."
	re := regexp.MustCompile(`^\d{2}:\d{2}:\d{2} From (.+?) to`)

	// FindStringSubmatch will return a slice with the entire match in the first element
	// and the captured groups (in this case, the sender's name) in the subsequent elements.

	var messages []Message
	var currentMessage *Message
	var contentBuilder strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Add the last message if EOF is reached
				if currentMessage != nil && contentBuilder.Len() > 0 {
					currentMessage.Snippet.DisplayMessage = contentBuilder.String()
					messages = append(messages, *currentMessage)
				}
				break
			}
			return nil, err
		}

		if re.MatchString(line) {

			// Extract timestamp, sender, and start of content
			parts := strings.SplitN(line, " ", 6)
			if len(parts) < 6 {
				continue // Invalid line format
			}
			timestamp, err := TimeStringToUnix(parts[0])
			if err != nil {
				continue // Invalid timestamp
			}

			// Skip the already processed messages
			if timestamp <= lastTime {
				continue
			}

			// If we are starting a new message, save the previous one (if any)
			if currentMessage != nil && contentBuilder.Len() > 0 {
				currentMessage.Snippet.DisplayMessage = contentBuilder.String()
				messages = append(messages, *currentMessage)
				contentBuilder.Reset()
			}

			matches := re.FindStringSubmatch(line)
			var sender string
			if len(matches) < 2 {
				sender = "Anonymous"
			} else {
				sender = matches[1]
			}
			// Start a new message
			currentMessage = &Message{
				ID:        uuid.New().String(),
				Timestamp: timestamp,
				AuthorDetails: Author{
					DisplayName:     sender,
					ProfileImageUrl: getProfileURL(sender),
				},
				Platform: PlatformDetail{
					Name:    "Zoom",
					LogoUrl: "https://seeklogo.com/images/Z/zoom-app-logo-B3FD9D4973-seeklogo.com.png",
				},
			}
		} else if currentMessage != nil {
			// Append to the current message's content
			contentBuilder.WriteString(" " + strings.TrimSpace(line))
		}
	}

	return messages, nil
}

func sendMessageDataToH2R(message MessageData, url string) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// fmt.Println(string(jsonData))

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and print the response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url := os.Getenv("H2R_SOCIAL_URL")
	if url == "" {
		log.Fatal("H2R_SOCIAL_URL environment variable not set")
	}

	chatFile := os.Getenv("ZOOM_CHAT_FILE")
	if chatFile == "" {
		log.Fatal("ZOOM_CHAT_FILE environment variable not set")
	}

	lastTime := time.Now().AddDate(0, 0, -1).Unix()
	for {
		// Read and parse the file
		file, err := os.Open(chatFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
		}

		// Encode messages to JSON
		reader := bufio.NewReader(file)
		messages, err := ParseMessages(reader, lastTime)
		if err != nil {
			fmt.Println("Error parsing messages:", err)
		}

		file.Close()
		if len(messages) > 0 {
			msgData := MessageData{Messages: messages}
			sendMessageDataToH2R(msgData, url)
			lastTime = messages[len(messages)-1].Timestamp
		}

		time.Sleep(time.Second)
	}
}
