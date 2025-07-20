package main

import (
	"crypto/sha256"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"

	"github.com/crochee/kim/internal/handle"
	"github.com/crochee/kim/internal/storage"
)

const (
	pathLoggedOut = "/logout"
)

type Storage interface {
	op.Storage
	handle.Authenticate
	handle.DeviceAuthenticate
}

func getUserStore(issuer, usersFile string) (storage.UserStore, error) {
	if usersFile == "" {
		return storage.NewUserStore(issuer), nil
	}
	return storage.StoreFromFile(usersFile)
}

// SetupServer creates an OIDC server with Issuer=http://localhost:<port>
//
// Use one of the pre-made clients in storage/clients.go or register a new one.
func SetupServer(issuer, usersFile string, redirectURI []string) (chi.Router, error) {
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}),
	)

	storage.RegisterClients(
		storage.NativeClient("native", redirectURI...),
		storage.WebClient("web", "secret", redirectURI...),
		storage.WebClient("api", "secret", redirectURI...),
	)

	// the OpenIDProvider interface needs a Storage interface handling various checks and state manipulations
	// this might be the layer for accessing your database
	// in this example it will be handled in-memory
	store, err := getUserStore(issuer, usersFile)
	if err != nil {
		mainLog.Error(err, "cannot create UserStore")
		return nil, err
	}
	storage := storage.NewStorage(store)
	// the OpenID Provider requires a 32-byte key for (token) encryption
	// be sure to create a proper crypto random key and manage it securely!

	router := chi.NewRouter()
	router.Use(logging.Middleware(
		logging.WithLogger(logger),
	))

	// for simplicity, we provide a very small default page for users who have signed out
	router.HandleFunc(pathLoggedOut, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("signed out successfully"))
		// no need to check/log error, this will be handled by the middleware.
	})

	// creation of the OpenIDProvider with the just created in-memory Storage
	provider, err := newOP(storage, issuer, logger)
	if err != nil {
		return nil, err
	}

	// the provider will only take care of the OpenID Protocol, so there must be some sort of UI for the login process
	// for the simplicity of the example this means a simple page with username and password field
	// be sure to provide an IssuerInterceptor with the IssuerFromRequest from the OP so the login can select / and pass it to the storage
	l := handle.NewLogin(storage, op.AuthCallbackURL(provider), op.NewIssuerInterceptor(provider.IssuerFromRequest))

	// regardless of how many pages / steps there are in the process, the UI must be registered in the router,
	// so we will direct all calls to /login to the login UI
	router.Mount("/login/", http.StripPrefix("/login", l))

	router.Route("/device", func(r chi.Router) {
		handle.RegisterDeviceAuth(storage, r)
	})

	handler := http.Handler(provider)
	// we register the http handler of the OP on the root, so that the discovery endpoint (/.well-known/openid-configuration)
	// is served on the correct path
	//
	// if your issuer ends with a path (e.g. http://localhost:9998/custom/path/),
	// then you would have to set the path prefix (/custom/path/)
	router.Mount("/", handler)

	return router, nil
}

// newOP will create an OpenID Provider for localhost on a specified port with a given encryption key
// and a predefined default logout uri
// it will enable all options (see descriptions)
func newOP(storage op.Storage, issuer string, logger *slog.Logger, extraOptions ...op.Option) (op.OpenIDProvider, error) {
	config := &op.Config{
		CryptoKey: sha256.Sum256([]byte("test")),

		// will be used if the end_session endpoint is called without a post_logout_redirect_uri
		DefaultLogoutRedirectURI: pathLoggedOut,

		// enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,

		// enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,

		// enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,

		// enables refresh_token grant use
		GrantTypeRefreshToken: true,

		// enables use of the `request` Object parameter
		RequestObjectSupported: true,

		// this example has only static texts (in English), so we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.English},

		DeviceAuthorization: op.DeviceAuthorizationConfig{
			Lifetime:     5 * time.Minute,
			PollInterval: 5 * time.Second,
			UserFormPath: "/device",
			UserCode:     op.UserCodeBase20,
		},
	}
	return op.NewProvider(config, storage,
		op.IssuerFromForwardedOrHost("https://localhost:9998"),
		// we must explicitly allow the use of the http issuer
		op.WithAllowInsecure(),
		// as an example on how to customize an endpoint this will change the authorization_endpoint from /authorize to /auth
		op.WithCustomAuthEndpoint(op.NewEndpoint("auth")),
		// Pass our logger to the OP
		op.WithLogger(logger.WithGroup("op")),
	)
}
