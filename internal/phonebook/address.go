package phonebook

import (
	"context"
	"errors"
)

type Address struct {
	ID          int
	User        User
	Name        string
	PhoneNumber string
}

type AddressRepository interface {
	NewAddress(context.Context, *Address) error
	Addresses(context.Context) ([]*Address, error)
	GetAddressesByUserID(context.Context, int) ([]*Address, error)
	GetAddressByID(context.Context, int) (*Address, error)
	GetAddressByName(context.Context, string) (*Address, error)
	GetAddressByPhoneNumber(context.Context, string) (*Address, error)
	UpdateAddress(context.Context, int, *Address) error
	DeleteAddress(context.Context, int) error
}

type AddressService struct {
	repo AddressRepository
}

func NewAddressService(repo AddressRepository) AddressService {
	return AddressService{repo}
}

func (s *AddressService) NewAddress(ctx context.Context, userID int, address *Address) error {
	address.User = User{ID: userID}

	err := s.repo.NewAddress(ctx, address)
	if err != nil {
		return err
	}

	return nil
}

func (s *AddressService) Addresses(ctx context.Context) ([]*Address, error) {
	addresses, err := s.repo.Addresses(ctx)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *AddressService) GetAddressesByUserID(ctx context.Context, userID int) ([]*Address, error) {
	addresses, err := s.repo.GetAddressesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *AddressService) GetAddressByID(ctx context.Context, ID int) (*Address, error) {
	address, err := s.repo.GetAddressByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *AddressService) GetAddressByName(ctx context.Context, name string) (*Address, error) {
	address, err := s.repo.GetAddressByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *AddressService) GetAddressByPhoneNumber(ctx context.Context, phoneNumber string) (*Address, error) {
	address, err := s.repo.GetAddressByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, userID int, addressID int, newAddress *Address) error {
	address, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return err
	}

	if address.User.ID != userID {
		return errors.New("unauthorized update")
	}

	err = s.repo.UpdateAddress(ctx, addressID, newAddress)
	if err != nil {
		return err
	}

	return nil
}

func (s *AddressService) DeleteAddress(ctx context.Context, userID int, addressID int) error {
	address, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return err
	}

	if address.User.ID != userID {
		return errors.New("unauthorized delete")
	}

	err = s.repo.DeleteAddress(ctx, address.ID)
	if err != nil {
		return err
	}

	return nil
}
