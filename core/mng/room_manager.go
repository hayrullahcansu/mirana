package mng

const (
	DefaultMode1RoomLenght = 100
)

type RoomManager struct {
	// NormalGames    map[*NormalGame]int
	mode1RoomLimit int
}

// func CreateNewRoomManager(limit int) *RoomManager {
// 	return &RoomManager{
// 		NormalGames:    make(map[*NormalGame]int),
// 		mode1RoomLimit: limit,
// 	}
// }

// func (roomManager *RoomManager) Init() {
// 	for i := 0; i < DefaultMode1RoomLenght; i++ {
// 		// room := CreateNewNormalGame()
// 		// roomManager.NormalGames[room] = 0
// 		go room.Run()
// 	}
// }

// func (roomManager *RoomManager) GetMode1Room() (*NormalGame, bool) {
// for room, status := range roomManager.NormalGames {
// 	if status == 0 {
// 		return room, true
// 	}
// }
// 	return nil, false
// }
