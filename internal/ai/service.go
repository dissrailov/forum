package ai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"forum/internal/models"
	"forum/internal/repo"
	"log"
	"net/http"
	"os"
	"strings"
)

type Service struct {
	client *Client
	repo   repo.RepoI
	logger *log.Logger
}

func NewService(client *Client, repo repo.RepoI, logger *log.Logger) *Service {
	return &Service{
		client: client,
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) IsConfigured() bool {
	return s.client.IsConfigured()
}

func (s *Service) GenerateAndStore(postID int, title, content string, imageURL string) {
	if !s.client.IsConfigured() {
		return
	}

	prompt := fmt.Sprintf("Forum post title: %s\n\nPost content: %s\n\nProvide a helpful AI response to this fitness/health forum post.", title, content)

	var response string
	var err error

	if imageURL != "" {
		if imgData, imgErr := loadImageAsDataURL(imageURL); imgErr == nil {
			prompt += "\n\nThe user also attached an image. Describe what you see in it and incorporate that into your response."
			response, err = s.client.GenerateResponseWithImage(prompt, imgData)
		} else {
			s.logger.Printf("AI image load failed for post %d: %v", postID, imgErr)
			response, err = s.client.GenerateResponse(prompt)
		}
	} else {
		response, err = s.client.GenerateResponse(prompt)
	}
	if err != nil {
		s.logger.Printf("AI generation failed for post %d: %v", postID, err)
		return
	}

	similarPosts := s.findSimilarPosts(postID, title)
	similarJSON, err := json.Marshal(similarPosts)
	if err != nil {
		s.logger.Printf("AI similar posts marshal failed for post %d: %v", postID, err)
		similarJSON = []byte("[]")
	}

	if err := s.repo.CreateAIResponse(postID, response, string(similarJSON)); err != nil {
		s.logger.Printf("AI store failed for post %d: %v", postID, err)
	}
}

func (s *Service) findSimilarPosts(currentPostID int, title string) []models.SimilarPost {
	allPosts, err := s.repo.GetAllPosts()
	if err != nil {
		return nil
	}

	words := extractKeywords(title)
	if len(words) == 0 {
		return nil
	}

	type scored struct {
		id    int
		title string
		score int
	}

	var matches []scored
	for _, p := range allPosts {
		if p.ID == currentPostID {
			continue
		}
		score := 0
		lowerTitle := strings.ToLower(p.Title)
		lowerContent := strings.ToLower(p.Content)
		for _, w := range words {
			if strings.Contains(lowerTitle, w) {
				score += 2
			}
			if strings.Contains(lowerContent, w) {
				score++
			}
		}
		if score > 0 {
			matches = append(matches, scored{p.ID, p.Title, score})
		}
	}

	// Sort by score descending (simple bubble sort for small slice)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].score > matches[i].score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	limit := 3
	if len(matches) < limit {
		limit = len(matches)
	}

	result := make([]models.SimilarPost, limit)
	for i := 0; i < limit; i++ {
		result[i] = models.SimilarPost{ID: matches[i].id, Title: matches[i].title}
	}
	return result
}

// loadImageAsDataURL reads a local image file (from a /uploads/... path) and returns a data URL.
func loadImageAsDataURL(imageURL string) (string, error) {
	// imageURL is like "/uploads/abc.jpg" — convert to local path
	filePath := "." + imageURL
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read image: %w", err)
	}

	mimeType := http.DetectContentType(data)
	b64 := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, b64), nil
}

func extractKeywords(text string) []string {
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "can": true, "shall": true, "to": true,
		"of": true, "in": true, "for": true, "on": true, "with": true,
		"at": true, "by": true, "from": true, "as": true, "into": true,
		"about": true, "like": true, "through": true, "after": true, "over": true,
		"and": true, "but": true, "or": true, "not": true, "no": true,
		"it": true, "its": true, "my": true, "your": true, "his": true,
		"her": true, "our": true, "this": true, "that": true, "what": true,
		"how": true, "i": true, "me": true, "we": true, "you": true,
		"he": true, "she": true, "they": true, "them": true,
		"и": true, "в": true, "на": true, "с": true, "по": true,
		"для": true, "не": true, "что": true, "как": true, "я": true,
		"мы": true, "он": true, "она": true, "они": true, "это": true,
		"от": true, "из": true, "за": true, "к": true, "но": true,
		"или": true, "ли": true, "же": true, "бы": true, "мой": true,
	}

	words := strings.Fields(strings.ToLower(text))
	var keywords []string
	for _, w := range words {
		w = strings.Trim(w, ".,!?;:\"'()-")
		if len(w) < 3 || stopWords[w] {
			continue
		}
		keywords = append(keywords, w)
	}
	return keywords
}
