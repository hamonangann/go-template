package handler

import (
	"context"
	"net/http"
	"strconv"
	"template/internal/phonebook"

	"github.com/gin-gonic/gin"
)

type UserJSON struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AddressJSON struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type UserService interface {
	Register(context.Context, *phonebook.User) (string, error)
	Login(context.Context, *phonebook.User) (string, error)
}

type AddressService interface {
	NewAddress(ctx context.Context, userID int, address *phonebook.Address) error
	Addresses(ctx context.Context) ([]*phonebook.Address, error)
	GetAddressesByUserID(ctx context.Context, userID int) ([]*phonebook.Address, error)
	GetAddressByID(ctx context.Context, ID int) (*phonebook.Address, error)
	UpdateAddress(ctx context.Context, userID int, addressID int, newAddress *phonebook.Address) error
	DeleteAddress(ctx context.Context, userID int, addressID int) error
}

type RESTHandler struct {
	userSvc    UserService
	addressSvc AddressService
}

func NewRESTHandler(userSvc UserService, addressSvc AddressService) RESTHandler {
	return RESTHandler{userSvc, addressSvc}
}

func (h *RESTHandler) Register(ctx *gin.Context) {
	var input UserJSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	user := &phonebook.User{
		Email:    input.Email,
		Password: input.Password,
	}

	token, err := h.userSvc.Register(ctx, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("Authorization", "Bearer "+token)
	ctx.JSON(
		http.StatusCreated,
		gin.H{"message": "success", "access_token": token},
	)
}

func (h *RESTHandler) Login(ctx *gin.Context) {
	var input UserJSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	user := &phonebook.User{
		Email:    input.Email,
		Password: input.Password,
	}

	token, err := h.userSvc.Login(ctx, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("Authorization", "Bearer "+token)
	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "success", "access_token": token},
	)
}

func (h *RESTHandler) NewAddress(ctx *gin.Context) {
	var input AddressJSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	userID := ctx.GetInt("user_id")

	address := &phonebook.Address{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
	}

	err := h.addressSvc.NewAddress(ctx, userID, address)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{"message": "success"},
	)
}

func (h *RESTHandler) Addresses(ctx *gin.Context) {
	addresses, err := h.addressSvc.Addresses(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	addressesResponse := make([]AddressJSON, 0)
	for _, address := range addresses {
		addressesResponse = append(addressesResponse, AddressJSON{
			ID:          address.ID,
			UserID:      address.User.ID,
			Name:        address.Name,
			PhoneNumber: address.PhoneNumber,
		})
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "success", "data": addressesResponse},
	)
}

func (h *RESTHandler) GetAddressesByUserID(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	addresses, err := h.addressSvc.GetAddressesByUserID(ctx, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	addressesResponse := make([]AddressJSON, 0)
	for _, address := range addresses {
		addressesResponse = append(addressesResponse, AddressJSON{
			ID:          address.ID,
			UserID:      address.User.ID,
			Name:        address.Name,
			PhoneNumber: address.PhoneNumber,
		})
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "success", "data": addressesResponse},
	)
}

func (h *RESTHandler) GetAddressByID(ctx *gin.Context) {
	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	address, err := h.addressSvc.GetAddressByID(ctx, addressID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "success", "data": AddressJSON{
			ID:          addressID,
			UserID:      address.User.ID,
			Name:        address.Name,
			PhoneNumber: address.PhoneNumber,
		}},
	)
}

func (h *RESTHandler) UpdateAddress(ctx *gin.Context) {
	var input AddressJSON
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	address := &phonebook.Address{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
	}

	userID := ctx.GetInt("user_id")

	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.addressSvc.UpdateAddress(ctx, userID, addressID, address)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "success"},
	)
}

func (h *RESTHandler) DeleteAddress(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")

	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.addressSvc.DeleteAddress(ctx, userID, addressID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"message": "success"},
	)
}
