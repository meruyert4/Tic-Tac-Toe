package game

import "tic-tac-toe/backend/utils"

func SaveGameState(gameId string, state string) error {
	return utils.Rdb.Set(utils.Ctx, "game:"+gameId, state, 0).Err()
}

func GetGameState(gameId string) (string, error) {
	return utils.Rdb.Get(utils.Ctx, "game:"+gameId).Result()
}
