package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"fyne.io/fyne/v2/app"
	"github.com/dvalkoff/gomessenger/gui/internal/config"
	"github.com/dvalkoff/gomessenger/gui/internal/controller"
	"github.com/dvalkoff/gomessenger/gui/internal/events"
	"github.com/dvalkoff/gomessenger/gui/internal/integration/api"
	"github.com/dvalkoff/gomessenger/gui/internal/integration/repository"
	"github.com/dvalkoff/gomessenger/gui/internal/model"
	"github.com/dvalkoff/gomessenger/gui/internal/view"
)

const (
	AppId                  = "minmessenger"
	x25519SharedSecretSize = 32
)

func main() {
	minMessengerApp := app.NewWithID(AppId)
	minMessengerWindow := minMessengerApp.NewWindow(AppId)
	appDataPath := minMessengerApp.Storage().RootURI().Path()
	db, err := config.ConfigureDB(appDataPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	workspaceRepository := repository.NewWorkspaceRepository(db)
	userRepository := repository.NewUserRepository(db)

	userApi := api.NewUserApi()

	eventStream := make(chan events.Event, 1)

	workspaceService := model.NewWorkspaceService(workspaceRepository, eventStream)
	userService := model.NewUserService(userApi, userRepository, eventStream)

	workspaceView := &view.WorkspaceView{
		Window:      minMessengerWindow,
		EventStream: eventStream,
	}
	signInView := &view.SignInView{
		Window:      minMessengerWindow,
		EventStream: eventStream,
	}
	signUpView := &view.SignUpView{
		Window:      minMessengerWindow,
		EventStream: eventStream,
	}
	eventProcessor := application.NewEventProcessor(
		eventStream,
		minMessengerWindow,

		workspaceView,
		signInView,
		signUpView,

		workspaceService,
		userService,
	)

	minMessengerWindow.Show()
	go eventProcessor.Run()
	minMessengerApp.Run()
}

func x3dh() {
	rand := rand.Reader

	// hkdf data
	hash := sha256.New
	keyLen := hash().Size()
	salt := make([]byte, 32)
	info := "X3DH"
	// Bob keys

	curve := ecdh.X25519()
	bobIdentityKey, _ := curve.GenerateKey(rand)
	bobSignedPreKey, _ := curve.GenerateKey(rand)
	bobOneTimePreKey, _ := curve.GenerateKey(rand)

	// Alice keys
	aliceIdentityKey, _ := curve.GenerateKey(rand)
	// aliceSignedPreKey, _ := curve.GenerateKey(rand)
	// aliceOneTimePreKey, _ := curve.GenerateKey(rand)

	// Alice wants to establish a shared secret. To do that, she gets Bob's key bundle
	bobPubIK := bobIdentityKey.PublicKey()
	bobPubSPK := bobSignedPreKey.PublicKey()
	bobPubOPK := bobOneTimePreKey.PublicKey()

	// then Alice needs to verify Bob's identity. But we will skip that part for now

	// then Alice should generate a Ephemeral Key
	aliceEphemeralKey, _ := curve.GenerateKey(rand)
	// and now she has to perform DH calculations
	dh := make([]byte, x25519SharedSecretSize, x25519SharedSecretSize*5)

	dhTemp, _ := aliceIdentityKey.ECDH(bobPubSPK)
	dh = append(dh, dhTemp...)
	dhTemp, _ = aliceEphemeralKey.ECDH(bobPubIK)
	dh = append(dh, dhTemp...)
	dhTemp, _ = aliceEphemeralKey.ECDH(bobPubSPK)
	dh = append(dh, dhTemp...)
	dhTemp, _ = aliceEphemeralKey.ECDH(bobPubOPK)
	dh = append(dh, dhTemp...)

	sharedSecretAlice, _ := hkdf.Key(hash, dh, salt, info, keyLen)

	// Now alice has aquired a secret. Alice sends Bob her ephemeral pub key and her identity key
	alicePubIK := aliceIdentityKey.PublicKey()
	alicePupEK := aliceEphemeralKey.PublicKey()
	// Now Bob can do the same thing Alice did with his keys
	dh = make([]byte, x25519SharedSecretSize, x25519SharedSecretSize*5)

	dhTemp, _ = bobSignedPreKey.ECDH(alicePubIK)
	dh = append(dh, dhTemp...)
	dhTemp, _ = bobIdentityKey.ECDH(alicePupEK)
	dh = append(dh, dhTemp...)
	dhTemp, _ = bobSignedPreKey.ECDH(alicePupEK)
	dh = append(dh, dhTemp...)
	dhTemp, _ = bobOneTimePreKey.ECDH(alicePupEK)
	dh = append(dh, dhTemp...)

	sharedSecretBob, _ := hkdf.Key(hash, dh, salt, info, keyLen)

	plaintext := []byte("hello bob")
	block, _ := aes.NewCipher(sharedSecretAlice)
	aead, _ := cipher.NewGCM(block)
	nonce := make([]byte, aead.NonceSize())
	rand.Read(nonce)
	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	block, _ = aes.NewCipher(sharedSecretBob)
	aead, _ = cipher.NewGCM(block)
	decrypted, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(decrypted))
}
