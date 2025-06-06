package services

import (
	"math/rand"
	"strconv"
	"time"
)

type Nickname struct {
	Nickname string `json:"nickname"`
}

var (
	prefixA = []string{"Ai", "Aq", "Al", "As", "Ain", "An", "Aya", "Ar", "Azu", "Aby", "Aga", "Ash"}
	suffixA = []string{"gul", "zhan", "khan", "bike", "dana", "zada", "zere", "na", "sulu", "sha", "aru", "tyn", "nur", "ma", "li", "da"}

	prefixZh = []string{"Zh", "Zhe", "Zhan"}
	suffixZh = []string{"an", "nur", "bek", "ar", "at", "lan"}

	prefixN = []string{"Nur", "Nurz", "Nuri"}
	suffixN = []string{"zhan", "bek", "ai", "lan", "gul", "zat"}

	prefixB = []string{"Ba", "Bi", "Be", "Bo", "Bai", "Bat"}
	suffixB = []string{"nur", "tch", "gdat", "lek", "bek"}

	prefixM = []string{"Mar", "Man", "Malik", "Mun", "Mira", "Madi"}
	suffixM = []string{"bek", "gul", "lan", "sha", "nur"}

	prefixT = []string{"Tal", "Tim", "Tay", "Tur", "Tolk", "Tari", "Temu"}
	suffixT = []string{"gul", "nar", "bek", "sha", "at", "tjin"}

	prefixD = []string{"Dar", "Dau", "Dina", "Dani", "Dos"}
	suffixD = []string{"bek", "gul", "at", "zhan", "sha"}

	prefixL = []string{"Lan", "Laz", "Lai", "Liya", "Lazat"}
	suffixL = []string{"gul", "sha", "bek", "nur"}

	prefixC = []string{"Cam", "Cal", "Car", "Caz"}
	suffixC = []string{"gul", "sha", "zan"}

	prefixS = []string{"Sam", "Sak", "San", "Sanz", "Ser", "Sau", "Sultan", "Sa"}
	suffixS = []string{"gul", "zhan", "nur", "sha", "ar", "na"}
)

func RandomNickName() string {
	rand.Seed(time.Now().UnixNano())

	category := rand.Intn(10)

	var prefix, suffix string

	switch category {
	case 0:
		prefix = prefixA[rand.Intn(len(prefixA))]
		suffix = suffixA[rand.Intn(len(suffixA))]
	case 1:
		prefix = prefixZh[rand.Intn(len(prefixZh))]
		suffix = suffixZh[rand.Intn(len(suffixZh))]
	case 2:
		prefix = prefixN[rand.Intn(len(prefixN))]
		suffix = suffixN[rand.Intn(len(suffixN))]
	case 3:
		prefix = prefixB[rand.Intn(len(prefixB))]
		suffix = suffixB[rand.Intn(len(suffixB))]
	case 4:
		prefix = prefixM[rand.Intn(len(prefixM))]
		suffix = suffixM[rand.Intn(len(suffixM))]
	case 5:
		prefix = prefixT[rand.Intn(len(prefixT))]
		suffix = suffixT[rand.Intn(len(suffixT))]
	case 6:
		prefix = prefixD[rand.Intn(len(prefixD))]
		suffix = suffixD[rand.Intn(len(suffixD))]
	case 7:
		prefix = prefixL[rand.Intn(len(prefixL))]
		suffix = suffixL[rand.Intn(len(suffixL))]
	case 8:
		prefix = prefixC[rand.Intn(len(prefixC))]
		suffix = suffixC[rand.Intn(len(suffixC))]
	case 9:
		prefix = prefixS[rand.Intn(len(prefixS))]
		suffix = suffixS[rand.Intn(len(suffixS))]
	}

	number := rand.Intn(90000) + 10000
	return prefix + suffix + strconv.Itoa(number)
}
