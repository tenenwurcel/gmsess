package config

import (
	"gmsess/utils"
)

var Verifier *utils.Verifier

func SetupVerifier() {
	v := utils.NewVerifier()
	setVerifier(v)
	go v.Start()
}

func setVerifier(v *utils.Verifier) {
	Verifier = v
}
