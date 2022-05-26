package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/alextanhongpin/promise"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	n := 10
	promises := make([]*promise.Promise[*User], n)
	for i := 0; i < n; i++ {
		agg := NewUserAggregate()
		promises[i] = agg.LoadAll(i)
	}

	fmt.Println("promise all")
	res, err := promise.AllSettled[*User](promises...).Await()
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

type UserAggregate struct {
	User *User
}

func NewUserAggregate() *UserAggregate {
	return &UserAggregate{
		User: new(User),
	}
}

func (u *UserAggregate) LoadUser(userID int) *promise.Promise[*User] {
	return promise.New(func() (*User, error) {
		fmt.Println("loading user", userID)
		time.Sleep(1 * time.Second)

		return &User{
			ID:         userID,
			Name:       fmt.Sprintf("name-%d", userID),
			IdentityID: fmt.Sprintf("identity-%d", userID),
			CountryID:  fmt.Sprintf("country-%d", userID),
		}, nil
	})
}

func (u *UserAggregate) LoadIdentity(identityID string) *promise.Promise[*Identity] {
	return promise.New(func() (*Identity, error) {
		fmt.Println("loading identity", identityID)
		time.Sleep(1 * time.Second)

		return &Identity{
			ID:         identityID,
			CardNumber: fmt.Sprintf("card-number-%s", identityID),
		}, nil
	})
}

func (u *UserAggregate) LoadCountry(countryID string) *promise.Promise[*Country] {
	return promise.New(func() (*Country, error) {
		fmt.Println("loading country", countryID)
		time.Sleep(1 * time.Second)

		return &Country{
			ISO: countryID,
		}, nil
	})
}

func (u *UserAggregate) LoadIdentityVerification(identityID string) *promise.Promise[*IdentityVerification] {
	return promise.New(func() (*IdentityVerification, error) {
		fmt.Println("loading identity verification", identityID)
		time.Sleep(1 * time.Second)

		return &IdentityVerification{
			IdentityID: identityID,
			Verified:   rand.Intn(2) == 0,
		}, nil
	})
}

func (u *UserAggregate) LoadAll(userID int) *promise.Promise[*User] {
	return promise.Then(u.LoadUser(userID), func(user *User) *promise.Promise[*User] {
		u.User = user

		loadIdentity := u.LoadIdentity(user.IdentityID)
		loadVerification := func(identity *Identity) *promise.Promise[*IdentityVerification] {
			u.User.Identity = identity

			return u.LoadIdentityVerification(identity.ID)
		}
		loadIdentityVerification := func(identityVerification *IdentityVerification) *promise.Promise[bool] {
			u.User.Identity.IdentityVerification = identityVerification
			return promise.Resolve(true)
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			_ = promise.Then(promise.Then(loadIdentity, loadVerification), loadIdentityVerification).AwaitResult()
		}()

		go func() {
			defer wg.Done()

			u.LoadCountry(user.CountryID).Await()
		}()

		return promise.Resolve(user)
	})
}

type User struct {
	ID   int
	Name string

	IdentityID string
	CountryID  string

	Identity *Identity
	Country  *Country
}

type Country struct {
	ISO string
}

type Identity struct {
	ID         string
	CardNumber string

	IdentityVerification *IdentityVerification
}

type IdentityVerification struct {
	IdentityID string
	Verified   bool
}
