package walletHandler

import (
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lutungp/julo-test/database"
	"github.com/lutungp/julo-test/internal/handler/utils"
	"github.com/lutungp/julo-test/internal/model"
)

type ResponseBalance struct {
	Data DataBalance `json:"wallet"`
}

type DataBalance struct {
	Id          uuid.UUID `json:"id"`
	CustomerXid string    `json:"deposited_by"`
	Status      string    `json:"status"`
	EnabledAt   time.Time `json:"enabled_at"`
	Amount      float64   `json:"amount"`
	ReferenceId string    `json:"reference_id"`
}

type ResponseDeposit struct {
	Data DataDeposit `json:"deposit"`
}

type DataDeposit struct {
	Id          uuid.UUID `json:"id"`
	CustomerXid string    `json:"deposited_by"`
	Status      string    `json:"status"`
	EnabledAt   time.Time `json:"enabled_at"`
	Amount      float64   `json:"amount"`
	ReferenceId string    `json:"reference_id"`
}

type ResponseWithdraw struct {
	Data DataWithdraw `json:"withdrawal"`
}

type DataWithdraw struct {
	Id          uuid.UUID `json:"id"`
	CustomerXid string    `json:"withdrawn_by"`
	Status      string    `json:"status"`
	EnabledAt   time.Time `json:"enabled_at"`
	Amount      float64   `json:"amount"`
	ReferenceId string    `json:"reference_id"`
}

type LastDeposit struct {
	CustomerId  string
	WalletId    string
	Value       float64
	ReferenceId string
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type ResponseWallet struct {
	Data DataWallet `json:"wallet"`
}

type DataWallet struct {
	Id          uuid.UUID `json:"id"`
	CustomerXid string    `json:"owned_by"`
	Status      string    `json:"status"`
	DisabledAt  time.Time `json:"disabled_at"`
	Amount      float64   `json:"amount"`
}

func EnableWalletCustomer(c *fiber.Ctx) error {
	db := database.DB
	wallet := new(model.Wallet)

	claims, err := utils.ExtractTokenMetadata(c)

	/**
	 * CEK AVAILABLE WALLET
	 */
	db.Find(&wallet, "status = 'enable' AND customer_id = ?", claims.Id)
	/**
	 * JIKA SUDAH ADA
	 */
	if wallet.ID != uuid.Nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Customer already has wallet", "data": err})
	}

	db.Find(&wallet, "customer_id = ?", claims.Id)

	message := "Created Wallet"
	if wallet.ID != uuid.Nil {
		wallet.CustomerId = claims.Id
		wallet.Status = "enable"

		err = db.Save(wallet).Error
		message = "Wallet Activated"
	} else {
		wallet.ID = uuid.New()
		wallet.CustomerId = claims.Id
		wallet.Status = "enable"
		err = db.Create(&wallet).Error
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create wallet", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": message, "data": wallet})
}

type DisableWalletValidator struct {
	Status bool `json:"is_disabled" xml:"is_disabled" form:"is_disabled" validate:"required"`
}

var validate = validator.New()

func DisableWalletStruct(disablecustomer DisableWalletValidator) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(disablecustomer)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func DisableWalletCustomer(c *fiber.Ctx) error {
	db := database.DB

	p := new(DisableWalletValidator)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	errors := DisableWalletStruct(*p)
	if errors != nil {

		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var account model.Account

	claims, err := utils.ExtractTokenMetadata(c)

	/**
	 * CEK AVAILABLE WALLET
	 */
	db.Find(&account, "id = ?", claims.Id)
	if account.ID == uuid.Nil {
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Customer Account Not Found", "data": err})
		}
	}

	wallet := new(model.Wallet)

	db.Find(&wallet, "customer_id = ?", claims.Id)
	/**
	 * JIKA SUDAH ADA
	 */
	if wallet.ID == uuid.Nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Customer doesn't have a wallet yet", "data": err})
	}

	status := ""
	if p.Status == true {
		status = "disabled"
	} else {
		status = "enable"
	}

	wallet.Status = status
	wallet.DisabledAt = time.Now()
	db.Save(wallet)

	disabled := &ResponseWallet{
		Data: DataWallet{
			Id:          wallet.ID,
			CustomerXid: wallet.CustomerId,
			Status:      "success",
			Amount:      wallet.Balance,
			DisabledAt:  wallet.DisabledAt,
		},
	}

	return c.JSON(fiber.Map{"status": "success", "data": disabled})
}

func GetBalance(c *fiber.Ctx) error {
	db := database.DB
	wallet := new(model.Wallet)

	claims, err := utils.ExtractTokenMetadata(c)
	db.Find(&wallet, "status='enable' AND customer_id = ?", claims.Id)
	/**
	 * JIKA BELUM MEMILIKI WALLET
	 */
	if wallet.ID == uuid.Nil || err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Customer doesn't have a wallet yet"})
	}

	var account model.Account
	db.Find(&account, "id = ?", claims.Id)

	balance := &ResponseBalance{
		Data: DataBalance{
			Id:          wallet.ID,
			CustomerXid: account.CustomerXid,
			Status:      "success",
			EnabledAt:   wallet.EnableAt,
			Amount:      wallet.Balance,
		},
	}

	return c.JSON(fiber.Map{"status": "success", "data": balance})
}

func getLastDeposit(id string) (*LastDeposit, error) {

	db := database.DB
	fillwallet := new(model.FillWallet)

	result := db.Last(&fillwallet, "customer_id = ?", id)
	if result.RowsAffected == 0 {
		return &LastDeposit{
			CustomerId:  "",
			WalletId:    "",
			Value:       0,
			ReferenceId: "",
		}, nil
	}

	return &LastDeposit{
		CustomerId:  fillwallet.CustomerId,
		WalletId:    fillwallet.WalletId,
		Value:       fillwallet.Value,
		ReferenceId: fillwallet.ReferenceId,
	}, nil
}

func DepositBalance(c *fiber.Ctx) error {
	db := database.DB

	wallet := new(model.Wallet)

	fillwallet := new(model.FillWallet)

	claims, err := utils.ExtractTokenMetadata(c)
	db.Find(&wallet, "customer_id = ?", claims.Id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn`t Deposit", "data": err})
	}

	amount, err := strconv.ParseFloat(c.FormValue("amount"), 64)
	reference_id := c.FormValue("reference_id")

	db.Find(&wallet, "customer_id = ?", claims.Id)

	if wallet.ID == uuid.Nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Customer doesn't have a wallet yet"})
	}

	fillwallet.ID = uuid.New()
	fillwallet.CustomerId = claims.Id
	fillwallet.ReferenceId = reference_id
	fillwallet.Value = amount
	db.Save(fillwallet)

	/**
	 * SIMPAN TOTAL PADA WALLET
	 */
	wallet.Balance = wallet.Balance + amount
	db.Save(wallet)

	deposit := &ResponseDeposit{
		Data: DataDeposit{
			Id:          fillwallet.ID,
			CustomerXid: fillwallet.CustomerId,
			Status:      "success",
			Amount:      fillwallet.Value,
			ReferenceId: fillwallet.ReferenceId,
		},
	}

	return c.JSON(fiber.Map{"status": "success", "data": deposit})
}

type DepositValidator struct {
	Amount      float64 `json:"amount" xml:"amount" form:"amount" validate:"required,number"`
	ReferenceId string  `json:"reference_id" xml:"reference_id" form:"reference_id" validate:"required"`
}

func DepositValidateStruct(deposit DepositValidator) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(deposit)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func WithDrawls(c *fiber.Ctx) error {
	p := new(DepositValidator)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	errors := DepositValidateStruct(*p)
	if errors != nil {

		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	db := database.DB
	wallet := new(model.Wallet)
	fillwallet := new(model.FillWallet)

	claims, err := utils.ExtractTokenMetadata(c)
	db.Find(&wallet, "customer_id = ?", claims.Id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn`t Deposit", "data": err})
	}

	amount := p.Amount * (-1)
	reference_id := p.ReferenceId

	db.Find(&wallet, "status='enable' AND customer_id = ?", claims.Id)

	if wallet.ID == uuid.Nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Customer doesn't have a wallet yet"})
	}

	fillwallet.ID = uuid.New()
	fillwallet.CustomerId = claims.Id
	fillwallet.ReferenceId = reference_id
	fillwallet.Value = amount
	db.Save(fillwallet)

	/**
	 * SIMPAN TOTAL PADA WALLET
	 */
	wallet.Balance = wallet.Balance + amount
	db.Save(wallet)

	deposit := &ResponseWithdraw{
		Data: DataWithdraw{
			Id:          fillwallet.ID,
			CustomerXid: fillwallet.CustomerId,
			Status:      "success",
			Amount:      p.Amount,
			ReferenceId: fillwallet.ReferenceId,
		},
	}

	return c.JSON(fiber.Map{"status": "success", "data": deposit})
}
