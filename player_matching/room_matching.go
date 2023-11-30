package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Room represents a chat room with a limit on the number of users it can accommodate.
type Room struct {
	ID           int
	Capacity     int
	CurrentUsers int
	mutex        sync.Mutex
}

// User represents a user who can join or leave a room.
type User struct {
	ID     int
	RoomID int
}

// Game represents the room matching system.
type Game struct {
	Rooms map[int]*Room
	Users map[int]*User
	mutex sync.Mutex

	availableRooms []*Room
}

var (
	r  = rand.New(rand.NewSource(time.Now().UnixNano()))
	wg = sync.WaitGroup{}
)

// NewRoom creates a new room with the given capacity.
func NewRoom(id, capacity int) *Room {
	return &Room{
		ID:           id,
		Capacity:     capacity,
		CurrentUsers: 0,
	}
}

// AddUser attempts to add a user to the specified room.
func (room *Room) AddUser(user *User) bool {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if room.CurrentUsers < room.Capacity {
		room.CurrentUsers++
		user.RoomID = room.ID
		fmt.Printf("User %d joined Room %d\n", user.ID, room.ID)
		if room.CurrentUsers == room.Capacity {
			fmt.Printf("Room %d full\n", room.ID)
		}
		return true
	}
	return false
}

// LeaveRoom removes a user from the specified room.
func (room *Room) LeaveRoom(user *User) bool {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	if room.CurrentUsers == 0 {
		return false
	}
	room.CurrentUsers--
	user.RoomID = -1
	fmt.Printf("User %d left Room %d\n", user.ID, room.ID)
	return true
}

// NewGame creates a new room matching system.
func NewGame() *Game {
	return &Game{
		Rooms: make(map[int]*Room),
		Users: make(map[int]*User),
	}
}

// AddRoom adds a new room to the system.
func (game *Game) AddRoom(roomID, capacity int) {
	game.Rooms[roomID] = NewRoom(roomID, capacity)
	game.availableRooms = append(game.availableRooms, game.Rooms[roomID])
}

// JoinGame adds a user to the system and attempts to find an available room.
func (game *Game) JoinGame(userID int) bool {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	user := &User{ID: userID, RoomID: -1}
	game.Users[userID] = user
	/*Assume that we can assign the user to any open room
	with open spots, a possible optimization is keep a
	queue of rooms with empty spots and always assign user to the head.

	Pop from the queue if the head becomes full. Enque if a full room becomes empty.

	Init queue with all rooms.*/

	/* Linear scanning:

	for _, room := range game.Rooms {
		if room.AddUser(user) {
			return true
		}
	}*/
	if len(game.availableRooms) == 0 {
		return false
	}
	headRoom := game.availableRooms[0]
	game.availableRooms[0].AddUser(user)
	if headRoom.Capacity == headRoom.CurrentUsers {
		game.availableRooms = game.availableRooms[1:]
	}
	return true
}

// LeaveGame removes a user from the system and the associated room.
func (game *Game) LeaveGame(userID int) bool {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	user, exists := game.Users[userID]
	if exists {
		roomID := user.RoomID
		delete(game.Users, userID)

		if roomID != -1 {
			room := game.Rooms[roomID]
			room.LeaveRoom(user)
			if room.CurrentUsers == 0 {
				game.availableRooms = append(game.availableRooms, room)
			}
		}
		return true
	}
	return false
}

// Play Game
func playGame(uid int, game *Game) {
	joinTime := time.Now()
	for !game.JoinGame(uid) {
		time.Sleep(time.Millisecond * 500)
	}
	waitTime := time.Now().Sub(joinTime)
	fmt.Printf("User %d waited %s\n", uid, waitTime)
	// Play for an arbitrary amount of time 5-7 seconds
	playTime := r.Intn(2) + 5
	time.Sleep(time.Duration(playTime) * time.Second)
	fmt.Printf("User %d played %s\n", uid, time.Duration(playTime)*time.Second)
	game.LeaveGame(uid)
	wg.Done()
}

func main() {
	// Create a room matching system
	game := NewGame()

	// Add rooms to the system
	game.AddRoom(1, 3)
	game.AddRoom(2, 2)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go playGame(i, game)
	}
	wg.Wait()
}
