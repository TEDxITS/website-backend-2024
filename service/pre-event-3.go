package service

import (
	"bytes"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/utils"
	"gorm.io/gorm"
)

type (
	PreEvent3Service interface {
		RegisterPE3(dto.PE3RSVPRegister, string) error
		GetStatus() (dto.PE3RSVPStatus, error)
	}

	preEvent3Service struct {
		eventRepo  repository.EventRepository
		userRepo   repository.UserRepository
		ticketRepo repository.TicketRepository
		bucketRepo repository.BucketRepository
	}
)

func NewPreEvent3Service(
	uRepo repository.UserRepository,
	tRepo repository.TicketRepository,
	eRepo repository.EventRepository,
	bRepo repository.BucketRepository,
) PreEvent3Service {
	return &preEvent3Service{
		eventRepo:  eRepo,
		userRepo:   uRepo,
		ticketRepo: tRepo,
		bucketRepo: bRepo,
	}
}

func (s *preEvent3Service) RegisterPE3(req dto.PE3RSVPRegister, userID string) error {
	event, err := s.eventRepo.GetByID(constants.PreEvent3ID)
	if err != nil {
		return err
	}

	if time.Now().Before(event.StartDate.Add(-7 * time.Hour)) {
		return dto.ErrPreEvent3NotYetOpen
	}

	if time.Now().After(event.EndDate.Add(-7 * time.Hour)) {
		return dto.ErrPreEvent3Closed
	}

	// validating uploaded file
	if req.PaymentFile.Size > dto.MB*5 {
		return dto.ErrMaxFileSize5MB
	}

	file, err := req.PaymentFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	fileBuffer := make([]byte, 512)
	if _, err := file.Read(fileBuffer); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	// only allow for jpeg/jpg/png
	fileType := http.DetectContentType(fileBuffer)
	if fileType != dto.ENUM_FILE_TYPE_JPEG && fileType != dto.ENUM_FILE_TYPE_PNG {
		return dto.ErrFileMustBeImage
	}
	ext := "." + utils.GetExtensions(req.PaymentFile.Filename)

	// generating unique code
	var code string
	for {
		code = utils.GenUniqueCode()
		if _, err := s.ticketRepo.GetTicketById(code); err != nil && err != gorm.ErrRecordNotFound {
			return err
		} else {
			break
		}
	}

	req.PaymentFile.Filename = code + ext
	err = s.bucketRepo.UploadFile(dto.ENUM_STORAGE_FOLDER_MAIN_EVENT, req.PaymentFile)
	if err != nil {
		return dto.ErrFailedToStorePaymentFile
	}

	False := false
	getFileEndpoint := dto.STORAGE_ENDPOINT_MAIN_EVENT + code + ext
	ticket := entity.Ticket{
		TicketID:         code,
		UserID:           userID,
		EventID:          constants.PreEvent3ID,
		Handphone:        req.Handphone,
		Birthdate:        req.Birthdate,
		Payment:          getFileEndpoint,
		PaymentConfirmed: &False,
		CheckedIn:        &False,
	}

	if _, err := s.ticketRepo.CreateTicket(ticket); err != nil {
		return err
	}

	// send email
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return err
	}

	readHtml, err := os.ReadFile("./utils/template/mail_payment_received.html")
	if err != nil {
		return err
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return err
	}

	var price string
	if event.Price >= 1000 {
		price = strconv.Itoa(event.Price)
		price = price[:len(price)-3] + "." + price[len(price)-3:]
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, struct {
		Name       string
		TicketType string
		TotalPrice string
	}{
		Name:       user.Name,
		TicketType: event.Name,
		TotalPrice: price,
	}); err != nil {
		return err
	}

	emailData := utils.Email{
		Email:   user.Email,
		Subject: "Payment Received",
		Body:    strMail.String(),
	}

	err = utils.SendMail(emailData)
	if err != nil {
		return dto.ErrSendEmail
	}

	return nil
}

func (s *preEvent3Service) GetStatus() (dto.PE3RSVPStatus, error) {
	status := true

	event, err := s.eventRepo.GetByID(constants.PreEvent3ID)
	if err != nil {
		return dto.PE3RSVPStatus{}, err
	}

	if time.Now().Before(event.StartDate.Add(-7 * time.Hour)) {
		status = false
	}

	if time.Now().After(event.EndDate.Add(-7 * time.Hour)) {
		status = false
	}

	return dto.PE3RSVPStatus{
		Status: &status,
	}, nil
}
