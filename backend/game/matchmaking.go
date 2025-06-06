package game

import "tic-tac-toe/backend/utils"

func AddToMatchQueue(userId string) error {
	return utils.Rdb.LPush(utils.Ctx, "match_queue", userId).Err()
}

func PopTwoPlayers() ([]string, error) {
	players := []string{}

	for i := 0; i < 2; i++ {
		player, err := utils.Rdb.RPop(utils.Ctx, "match_queue").Result()
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}
