package storage

import (
	"context"
	"crypto/rsa"
	"encoding/hex"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kimv1 "github.com/crochee/kim/api/kim/v1"
)

type Service struct {
	keys map[string]*rsa.PublicKey
}

type UserStore interface {
	GetUserByID(context.Context, string) (*kimv1.User, error)
	GetUserByUsername(context.Context, string) (*kimv1.User, error)
}

type userStore struct {
	client.Client
}

func (us *userStore) GetUserByID(ctx context.Context, userID string) (*kimv1.User, error) {
	decoded, err := hex.DecodeString(userID)
	if err != nil {
		return nil, err
	}
	return us.GetUserByUsername(ctx, string(decoded))
}

func (us *userStore) GetUserByUsername(ctx context.Context, name string) (*kimv1.User, error) {
	parts := strings.SplitN(name, "/", 2)
	user := &kimv1.User{}
	ns := types.NamespacedName{Name: parts[0], Namespace: parts[1]}
	if err := us.Get(ctx, ns, user); err != nil {
		return nil, err
	}
	return user, nil
}
