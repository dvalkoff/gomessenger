package user

import "github.com/google/uuid"

type UserRegistrationInfo struct {
	Nickname          string   `json:"nickname"`
	Password          string   `json:"password"`
	IdentityPubKey    string   `json:"identityPubKey"`
	SignedPubPreKey   string   `json:"signedPubPrekey"`
	OneTimePubPreKeys []string `json:"oneTimePubPrekeys"`
}

type UserRegistratedInfo struct {
	Id uuid.UUID `json:"id"`
}

type UserInfo struct {
	Id       uuid.UUID `json:"id"`
	Nickname string    `json:"nickname"`
}

type KeyBundle struct {
	UserId           uuid.UUID `json:"userId"`
	IdentityPubKey   string    `json:"identityPubKey"`
	SignedPubPreKey  string    `json:"signedPubPrekey"`
	OneTimePubPreKey string    `json:"oneTimePubPrekey"`
}
