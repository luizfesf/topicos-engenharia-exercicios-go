package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type PlayerRequest struct {
	Nickname string
	Life     int

	Attack   int
}

type PlayerResponse struct {
	Message string `json:"message"`
}

var players []PlayerRequest

type Enemy struct {
	Nickname string `json:"nickname"`
	Life     int    `json:"life"`
	Attack   int    `json:"attack"`
}

type EnemyResponse struct {
	Message string `json:"message"`
}

var enemies []Enemy

type BattleRequest struct {
	ID        string
	Enemy     string
	Player    string
	DiceThrow int
}

type BattleResponse struct {
	Id        string
	DiceThrow int
	Player    PlayerRequest
	Enemy     Enemy
}

type Response struct {
	Message string
}

var battles []BattleRequest

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if playerRequest.Nickname == "" || playerRequest.Life == 0 || playerRequest.Attack == 0 {
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname, life and attack is required"})
		return
	}

	if playerRequest.Attack > 10 || playerRequest.Attack <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player attack must be between 1 and 10"})
		return
	}

	if playerRequest.Life > 100 || playerRequest.Life <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player life must be between 1 and 100"})
		return
	}

	for _, player := range players {
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname already exits"})
			return
		}
	}

	player := PlayerRequest{
		Nickname: playerRequest.Nickname,
		Life:     playerRequest.Life,
		Attack:   playerRequest.Attack}
	players = append(players, player)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(player)
}

func LoadPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.PathValue("nickname")

	for i, player := range players {
		if player.Nickname == nickname {
			players = append(players[:i], players[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{
		Message: "Player nickname not found",
	})
}

func LoadPlayerByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.PathValue("nickname")

	for _, player := range players {
		if player.Nickname == nickname {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(player)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{
		Message: "Player nickname not found",
	})
}

func SavePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nickname := r.PathValue("nickname")

	var playerRequest PlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&playerRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if playerRequest.Nickname == "" {
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname is required"})
		return
	}

	indexPlayer := -1
	for i, player := range players {
		if player.Nickname == nickname {
			indexPlayer = i
		}
		if player.Nickname == playerRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Player nickname already exits"})
			return
		}
	}

	if indexPlayer != -1 {
		players[indexPlayer].Nickname = playerRequest.Nickname
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(players[indexPlayer])
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(PlayerResponse{
		Message: "Player nickname not found",
	})
}

func AddEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var enemyRequest Enemy
	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(EnemyResponse{Message: "Internal Server Error"})
		return
	}

	if enemyRequest.Nickname == "" {
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname is required"})
		return
	}

	for _, enemy := range enemies {
		if enemy.Nickname == enemyRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(EnemyResponse{Message: "Enemy nickname already exists"})
			return
		}
	}

	rand.Seed(time.Now().UnixNano())
	enemy := Enemy{
		Nickname: enemyRequest.Nickname,
		Life:     rand.Intn(10) + 1,
		Attack:   rand.Intn(10) + 1,
	}

	enemies = append(enemies, enemy)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enemy)
}

func LoadEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(enemies)
}

func LoadEnemyByNickname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	for _, enemy := range enemies {
		if enemy.Nickname == nickname {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(enemy)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func SaveEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	var enemyRequest Enemy
	if err := json.NewDecoder(r.Body).Decode(&enemyRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Internal Server Error"})
		return
	}

	if enemyRequest.Nickname == "" {
		json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname is required"})
		return
	}

	indexEnemy := -1
	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			indexEnemy = i
		}
		if enemy.Nickname == enemyRequest.Nickname {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname already exists"})
			return
		}
	}

	if indexEnemy != -1 {
		enemies[indexEnemy].Nickname = enemyRequest.Nickname
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(enemies[indexEnemy])
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func DeleteEnemy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := r.URL.Query().Get("nickname")

	for i, enemy := range enemies {
		if enemy.Nickname == nickname {
			enemies = append(enemies[:i], enemies[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(PlayerResponse{Message: "Enemy nickname not found"})
}

func AddBattle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var battleRequest BattleRequest
	if err := json.NewDecoder(r.Body).Decode(&battleRequest); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal Server Error"})
		return
	}

	if battleRequest.Enemy == "" {
		json.NewEncoder(w).Encode(Response{Message: "Enemy nickname is required"})
		return
	}

	if battleRequest.Player == "" {
		json.NewEncoder(w).Encode(Response{Message: "Player nickname is required"})
		return
	}

	indexEnemy := -1
	for i, enemy := range enemies {
		if enemy.Nickname == battleRequest.Enemy {
			indexEnemy = i
			break
		}
	}

	if indexEnemy < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Enemy dont exits"})
		return
	}

	indexPlayer := -1
	for i, player := range players {
		if player.Nickname == battleRequest.Player {
			indexPlayer = i
			break
		}
	}

	if indexPlayer < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Player dont exits"})
		return
	}

	diceThrow := rand.Intn(6-1) + 1
	battleID := uuid.New().String()

	if diceThrow >= 1 && diceThrow <= 3 {
		players[indexPlayer].Life -= enemies[indexEnemy].Attack
	}

	if diceThrow >= 4 && diceThrow <= 6 {
		enemies[indexEnemy].Life -= players[indexPlayer].Attack
	}
	battle := BattleRequest{
		ID:         battleID,
		Enemy:      battleRequest.Enemy,
		Player:     battleRequest.Player,
		DiceThrow:  diceThrow,
		
	}
	battles = append(battles, battle)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BattleResponse{
		Id:        battleID,
		DiceThrow: diceThrow,
		Player:    players[indexPlayer],
		Enemy:     enemies[indexEnemy],
	})

}

func LoadBattle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(battles)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /player", AddPlayer)
	mux.HandleFunc("GET /player", LoadPlayers)
	mux.HandleFunc("DELETE /player/{nickname}", DeletePlayer)
	mux.HandleFunc("GET /player/{nickname}", LoadPlayerByNickname)
	mux.HandleFunc("PUT /player/{nickname}", SavePlayer)

	mux.HandleFunc("POST /enemy", AddEnemy)
	mux.HandleFunc("GET /enemy", LoadEnemy)
	mux.HandleFunc("DELETE /enemy/{nickname}", DeleteEnemy)
	mux.HandleFunc("GET /enemy/{nickname}", LoadEnemyByNickname)
	mux.HandleFunc("PUT /enemy/{nickname}", SaveEnemy)

	mux.HandleFunc("POST /battle", AddBattle)
	mux.HandleFunc("GET /battle", LoadBattle)

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println(err)
	}
}
