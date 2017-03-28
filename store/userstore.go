package store

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/gobuild/log"

	"google.golang.org/appengine/datastore"
)

type Kind string

const UserKind = "userkind"

type User struct {
	ID     int64
	Email  string
	Notify bool
}

func (u *User) Kind() string {
	return UserKind
}

func NewUser(ctx context.Context, user *User) error {

	// check user doesn't already exist, otherwise update
	hits, err := GetUsers(ctx, user)
	if err != nil {
		return err
	}
	switch {
	case len(hits) > 1:
		log.Errorf("%#v\n", hits)
		return errors.New("too many users with this email address")
	case len(hits) == 1:
		hit := hits[0]
		// Only field that can be updated
		hit.Notify = user.Notify
		return UpdateUser(ctx, user)
	}

	key := datastore.NewIncompleteKey(ctx, user.Kind(), nil)
	if _, err := datastore.Put(ctx, key, user); err != nil {
		return err
	}

	user.ID = key.IntID()
	return nil
}

func GetUsers(ctx context.Context, user *User) ([]User, error) {
	if len(user.Email) == 0 {
		return nil, errors.New("email not provided")
	}
	q := datastore.NewQuery(user.Kind()).Filter("Email = ", user.Email)
	q = q.Limit(5)
	var users []User
	keys, err := q.GetAll(ctx, &users)
	for i := range users {
		users[i].ID = keys[i].IntID()
	}
	return users, err
}

func UpdateUser(ctx context.Context, user *User) error {
	key := datastore.NewKey(ctx, user.Kind(), "", user.ID, nil)
	_, err := datastore.Put(ctx, key, user)
	return err
}

func GetEmails(ctx context.Context, filter *User) ([]User, error) {
	var users []User
	q := datastore.NewQuery(filter.Kind()).Filter("Notify =", true)
	t := q.Run(ctx)
	for {
		var u User
		key, err := t.Next(&u)
		if err != nil {
			continue
		}
		u.ID = key.IntID()
		users = append(users, u)
	}
	return users, nil
}
