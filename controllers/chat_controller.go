package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"samsungvoicebe/config"
	"samsungvoicebe/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	config *config.Config
}

func NewChatController(cfg *config.Config) *ChatController {
	return &ChatController{
		config: cfg,
	}
}

func (cc *ChatController) Chat(c *gin.Context) {
	if !cc.config.IsValid() {
		c.JSON(http.StatusInternalServerError, models.ChatResponse{
			Error: "GEMINI_API_KEY not configured",
		})
		return
	}

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: "Invalid request format",
		})
		return
	}

	log.Printf("Chat request: %s", req.Message)

	screen, err := cc.determineScreen(req.Message)
	if err != nil {
		log.Printf("Error determining screen: %v", err)
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: err.Error(),
		})
		return
	}

	log.Printf("Determined screen: %s", screen)

	c.JSON(http.StatusOK, models.ChatResponse{
		Response: screen,
		Screen:   screen,
	})
}

func (cc *ChatController) determineScreen(message string) (string, error) {
	msg := strings.ToLower(message)

	if cc.isInappropriateContent(msg) {
		return "", fmt.Errorf("inappropriate content detected")
	}

	if len(strings.TrimSpace(msg)) < 2 {
		return "", fmt.Errorf("message too short or empty")
	}

	if cc.isGibberish(msg) {
		return "", fmt.Errorf("unable to understand the message")
	}

	playKeywords := []string{
		"main", "bermain", "permainan", "game", "play", "mulai",
		"beranda", "home", "utama", "awal", "start", "menu",
	}

	scanKeywords := []string{
		"scan", "pindai", "kamera", "camera", "foto", "gambar",
		"ambil", "capture", "barcode", "qr", "scanner",
	}

	lessonKeywords := []string{
		"belajar", "pelajaran", "lesson", "materi", "kursus",
		"tutorial", "panduan", "edukasi", "pembelajaran", "study",
	}

	analyzeKeywords := []string{
		"analisis", "analyze", "analisa", "periksa", "cek",
		"evaluasi", "tinjau", "review", "laporan", "data",
	}

	settingKeywords := []string{
		"pengaturan", "setting", "konfigurasi", "config", "atur",
		"preferensi", "opsi", "options", "setup", "setelan",
	}

	if cc.containsAny(msg, playKeywords) {
		return "play", nil
	}

	if cc.containsAny(msg, scanKeywords) {
		return "scan", nil
	}

	if cc.containsAny(msg, lessonKeywords) {
		return "lesson", nil
	}

	if cc.containsAny(msg, analyzeKeywords) {
		return "analyze", nil
	}

	if cc.containsAny(msg, settingKeywords) {
		return "setting", nil
	}

	return "", fmt.Errorf("unable to determine destination from your message")
}

func (cc *ChatController) isInappropriateContent(message string) bool {
	inappropriateWords := []string{
		"kontol", "memek", "anjing", "bangsat", "babi", "tai", "shit", "fuck",
		"bitch", "asshole", "damn", "hell",
	}

	for _, word := range inappropriateWords {
		if strings.Contains(message, word) {
			return true
		}
	}
	return false
}

func (cc *ChatController) isGibberish(message string) bool {
	cleaned := strings.ReplaceAll(message, " ", "")

	if len(cleaned) > 3 {
		repeated := true
		firstChar := cleaned[0]
		for _, char := range cleaned {
			if char != rune(firstChar) {
				repeated = false
				break
			}
		}
		if repeated {
			return true
		}
	}

	vowels := "aeiouAEIOU"
	hasVowel := false
	for _, char := range cleaned {
		if strings.ContainsRune(vowels, char) {
			hasVowel = true
			break
		}
	}

	// If no vowels and length > 2, probably gibberish
	if !hasVowel && len(cleaned) > 2 {
		return true
	}

	return false
}

func (cc *ChatController) containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (cc *ChatController) ChatWithAI(c *gin.Context) {
	if !cc.config.IsValid() {
		c.JSON(http.StatusInternalServerError, models.ChatResponse{
			Error: "GEMINI_API_KEY not configured",
		})
		return
	}

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ChatResponse{
			Error: "Invalid request format",
		})
		return
	}

	log.Printf("Chat request: %s", req.Message)

	enhancedPrompt := fmt.Sprintf(`
Analyze this user input and determine which screen they want to navigate to.
Available screens: play, scan, lesson, analyze, setting

User input: "%s"

Rules:
- If user wants to go to main menu, game, or start something: return "play"
- If user wants to scan, take photo, or use camera: return "scan"  
- If user wants to learn, study, or access lessons: return "lesson"
- If user wants to analyze, check, or review data: return "analyze"
- If user wants settings, configuration, or preferences: return "setting"

Respond with ONLY the screen name (play/scan/lesson/analyze/setting), no additional text.
`, req.Message)

	response, err := cc.callGeminiAPI(enhancedPrompt)
	if err != nil {
		log.Printf("Gemini API error, falling back to keyword matching: %v", err)
		// Fallback to keyword matching
		screen, fallbackErr := cc.determineScreen(req.Message)
		if fallbackErr != nil {
			log.Printf("Error in fallback determination: %v", fallbackErr)
			c.JSON(http.StatusBadRequest, models.ChatResponse{
				Error: fallbackErr.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, models.ChatResponse{
			Response: screen,
			Screen:   screen,
		})
		return
	}

	screen := strings.TrimSpace(strings.ToLower(response))
	validScreens := []string{"play", "scan", "lesson", "analyze", "setting"}

	isValid := false
	for _, validScreen := range validScreens {
		if screen == validScreen {
			isValid = true
			break
		}
	}

	if !isValid {
		fallbackScreen, err := cc.determineScreen(req.Message)
		if err != nil {
			log.Printf("Error in fallback determination: %v", err)
			c.JSON(http.StatusBadRequest, models.ChatResponse{
				Error: err.Error(),
			})
			return
		}
		screen = fallbackScreen
	}

	log.Printf("Determined screen: %s", screen)

	c.JSON(http.StatusOK, models.ChatResponse{
		Response: screen,
		Screen:   screen,
	})
}

func (cc *ChatController) callGeminiAPI(message string) (string, error) {
	geminiReq := models.GeminiRequest{
		Contents: []models.Content{
			{
				Parts: []models.Part{
					{Text: message},
				},
			},
		},
	}

	jsonData, err := json.Marshal(geminiReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", cc.config.GeminiAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp models.GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("gemini API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated by Gemini")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
