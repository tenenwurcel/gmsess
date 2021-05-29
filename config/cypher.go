package config

import (
	"crypto/rand"
	"gmsess/utils"
)

var cypher *utils.Cypher

func SetupCypher() {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)

	Cypher := utils.Cypher(randBytes)

	setCypher(&Cypher)
}

func setCypher(Cypher *utils.Cypher) {
	cypher = Cypher
}

func GetCypher() utils.Cypher {
	return *cypher
}
