package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Структура для формирования JSON сообщения
type Payment struct {
	PayID    string `json:"payId"`
	ShopID   string `json:"shopId"`
	Amount   int    `json:"amount"`
	Checksum string `json:"checksum"`
}

const privateKey = "PrivateKeyHere" // Замените на ваш приватный ключ

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateChecksum(payID, shopID, amount, privateKey string) string {
	data := []byte(payID + "&" + shopID + "&" + amount + "&" + privateKey)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method!= http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Amount int    `json:"amount"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err!= nil {
		fmt.Fprint(w, "err2")
		//http.Error(w, r, http.StatusBadRequest)
		return
	}
	payID := generateRandomString(16) // Генерация random ID
	checksum := generateChecksum(payID, "ShopID12", fmt.Sprint(data.Amount), privateKey)

	jsonData := map[string]interface{}{
		"payId":     payID,
		"shopId":    "ShopID12",
		"amount":    data.Amount,
		"checksum":  checksum,
	}
	responseBody, _ := json.Marshal(jsonData)

	req, _ := http.NewRequest("POST", "http://10.10.0.6/token", strings.NewReader(string(responseBody)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err!= nil {
		fmt.Fprint(w, "err1")
		//http.Error(w, r, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	fmt.Fprint(w, checksum[:6]) // Возвращаем первую часть checksum
}
func testhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "working")
}
func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/pay", handler)
	http.HandleFunc("/test", testhandler)
	fmt.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err!= nil {
		panic(err)
	}
}
