package service

import (
	"bytes"
	"context"
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
	"github.com/TEDxITS/website-backend-2024/websocket"
	"gorm.io/gorm"
)

type (
	MainEventService interface {
		RegisterMainEvent(context.Context, dto.MainEventRegister, string) error
		ConfirmPayment(context.Context, dto.MainEventConfirmPaymentRequest) error
		CheckIn(context.Context, dto.MainEventCheckInRequest) error
		GetStatus(context.Context) (dto.MainEventStatusResponse, error)
		GetMainEventPaginated(context.Context, dto.PaginationQuery) (dto.MainEventPaginationResponse, error)
		GetMainEventDetail(context.Context, string) (dto.MainEventResponse, error)
		GetMainEventCounter(context.Context) (dto.MainEventCounter, error)
	}

	mainEventService struct {
		eventRepo  repository.EventRepository
		userRepo   repository.UserRepository
		ticketRepo repository.TicketRepository
		bucketRepo repository.BucketRepository
		queueHub   []websocket.QueueHub
	}
)

func NewMainEventService(
	uRepo repository.UserRepository,
	tRepo repository.TicketRepository,
	eRepo repository.EventRepository,
	bRepo repository.BucketRepository,
	qHub []websocket.QueueHub,
) MainEventService {
	return &mainEventService{
		eventRepo:  eRepo,
		userRepo:   uRepo,
		ticketRepo: tRepo,
		bucketRepo: bRepo,
		queueHub:   qHub,
	}
}

func (s *mainEventService) RegisterMainEvent(ctx context.Context, req dto.MainEventRegister, userID string) error {
	hub, err := func() (websocket.QueueHub, error) {
		for _, hub := range s.queueHub {
			if hub.IsEventHandler(req.EventID) {
				return hub, nil
			}
		}
		return nil, dto.ErrEventNotFound
	}()

	if err != nil {
		return err
	}

	event, err := s.eventRepo.GetByID(req.EventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.ErrEventNotFound
		}
		return err
	}

	if time.Now().Before(event.StartDate.Add(-7 * time.Hour)) {
		return dto.ErrMainEventNotYetOpen
	}

	if time.Now().After(event.EndDate.Add(-7 * time.Hour)) {
		return dto.ErrMainEventClosed
	}

	if event.Registers >= event.Capacity {
		return dto.ErrMainEventFull
	}

	client := hub.GetClientInTransactionByUserID(userID)
	if client == nil {
		return dto.ErrUserNotInTransaction
	}

	if client.IsWithMerch() != *event.WithKit {
		return dto.ErrMismatchData
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
		EventID:          req.EventID,
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

	// signal the client to exit the handler thread
	// and sequentially unregister from the hub
	client.Done(nil)

	return nil
}

func (s *mainEventService) ConfirmPayment(ctx context.Context, req dto.MainEventConfirmPaymentRequest) error {
	ticket, err := s.ticketRepo.FindByTicketID(req.Code)
	if err != nil {
		return dto.ErrTicketNotFound
	}

	user, err := s.userRepo.GetUserById(ticket.UserID)
	if err != nil {
		return dto.ErrUserNotFound
	}

	confirmed := true
	ticket.PaymentConfirmed = &confirmed
	_, err = s.ticketRepo.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	readHtml, err := os.ReadFile("./utils/template/mail_confirmation_payment.html")

	if err != nil {
		return err
	}

	data := struct {
		Name     string
		TicketID string
	}{
		Name:     user.Name,
		TicketID: ticket.TicketID,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return err
	}

	emailData := utils.Email{
		Email:   user.Email,
		Subject: "Confirmation Payment",
		Body:    strMail.String(),
	}

	err = utils.SendMail(emailData)
	if err != nil {
		return dto.ErrSendEmail
	}

	return nil
}

func (s *mainEventService) CheckIn(ctx context.Context, req dto.MainEventCheckInRequest) error {
	ticket, err := s.ticketRepo.FindByTicketID(req.Code)
	if err != nil {
		return dto.ErrTicketNotFound
	}

	checked := true
	ticket.CheckedIn = &checked
	_, err = s.ticketRepo.UpdateTicket(ticket)
	if err != nil {
		return err
	}

	return nil
}

func (s *mainEventService) GetStatus(ctx context.Context) (dto.MainEventStatusResponse, error) {
	early_bird, err := s.eventRepo.GetByID(constants.MainEventEarlyBirdNoMerchID)
	if err != nil {
		return dto.MainEventStatusResponse{}, err
	}
	early_bird_with_merch, err := s.eventRepo.GetByID(constants.MainEventEarlyBirdWithMerchID)
	if err != nil {
		return dto.MainEventStatusResponse{}, err
	}

	pre_sale, err := s.eventRepo.GetByID(constants.MainEventPreSaleNoMerchID)
	if err != nil {
		return dto.MainEventStatusResponse{}, err
	}
	pre_sale_with_merch, err := s.eventRepo.GetByID(constants.MainEventPreSaleWithMerchID)
	if err != nil {
		return dto.MainEventStatusResponse{}, err
	}

	normal, err := s.eventRepo.GetByID(constants.MainEventNormalNoMerchID)
	if err != nil {
		return dto.MainEventStatusResponse{}, err
	}
	normal_with_merch, err := s.eventRepo.GetByID(constants.MainEventNormalWithMerchID)
	if err != nil {
		return dto.MainEventStatusResponse{}, err
	}

	preprocess := func(e entity.Event) dto.MainEventStatusDetail {
		status := dto.MAIN_EVENT_OPEN

		if time.Now().Before(e.StartDate) {
			status = dto.MAIN_EVENT_CLOSED
		}

		if time.Now().After(e.EndDate) {
			status = dto.MAIN_EVENT_FULL
		}

		startDateDifference := int(time.Until(e.StartDate.Add(-7 * time.Hour)).Seconds())
		endDateDifference := int(time.Until(e.EndDate.Add(-7 * time.Hour)).Seconds())

		return dto.MainEventStatusDetail{
			Status: status,
			UntilOpen: dto.RemainingTime{
				Days:    int(startDateDifference / (60 * 60 * 24)),
				Hours:   int((startDateDifference / (60 * 60)) % 24),
				Minutes: int((startDateDifference / 60) % 60),
				Seconds: int(startDateDifference % 60),
			},
			UntilClosed: dto.RemainingTime{
				Days:    int(endDateDifference / (60 * 60 * 24)),
				Hours:   int((endDateDifference / (60 * 60)) % 24),
				Minutes: int((endDateDifference / 60) % 60),
				Seconds: int(endDateDifference % 60),
			},
		}

	}

	res := dto.MainEventStatusResponse{
		EarlyBird: preprocess(early_bird),
		PreSale:   preprocess(pre_sale),
		Normal:    preprocess(normal),
	}

	if (early_bird.Capacity-early_bird.Registers)+(early_bird_with_merch.Registers-early_bird_with_merch.Capacity) <= 0 {
		res.EarlyBird.Status = dto.MAIN_EVENT_FULL
	}

	if (pre_sale.Capacity-pre_sale.Registers)+(pre_sale_with_merch.Registers-pre_sale_with_merch.Capacity) <= 0 {
		res.PreSale.Status = dto.MAIN_EVENT_FULL
	}

	if (normal.Capacity-normal.Registers)+(normal_with_merch.Registers-normal_with_merch.Capacity) <= 0 {
		res.Normal.Status = dto.MAIN_EVENT_FULL
	}

	res.EarlyBird.NoMerchID = constants.MainEventEarlyBirdNoMerchID
	res.PreSale.NoMerchID = constants.MainEventPreSaleNoMerchID
	res.Normal.NoMerchID = constants.MainEventNormalNoMerchID
	res.EarlyBird.WithMerchID = constants.MainEventEarlyBirdWithMerchID
	res.PreSale.WithMerchID = constants.MainEventPreSaleWithMerchID
	res.Normal.WithMerchID = constants.MainEventNormalWithMerchID

	return res, nil
}

func (s *mainEventService) GetMainEventPaginated(ctx context.Context, req dto.PaginationQuery) (dto.MainEventPaginationResponse, error) {
	var limit int
	var page int

	limit = req.PerPage
	if limit <= 0 {
		limit = constants.ENUM_PAGINATION_LIMIT
	}

	page = req.Page
	if page <= 0 {
		page = constants.ENUM_PAGINATION_PAGE
	}

	tickets, maxPage, count, err := s.ticketRepo.JoinGetAllPagination(req.Search, limit, page)
	if err != nil {
		return dto.MainEventPaginationResponse{}, err
	}

	var result []dto.MainEventPaginationData
	for _, t := range tickets {
		result = append(result, dto.MainEventPaginationData{
			ID:        t.TicketID,
			Name:      t.User.Name,
			Email:     t.User.Email,
			Confirmed: *t.PaymentConfirmed,
			CheckedIn: *t.CheckedIn,
			EventName: t.Event.Name,
			Price:     t.Event.Price,
		})
	}

	return dto.MainEventPaginationResponse{
		Data: result,
		PaginationMetadata: dto.PaginationMetadata{
			Page:    page,
			PerPage: limit,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}

func (s *mainEventService) GetMainEventDetail(ctx context.Context, id string) (dto.MainEventResponse, error) {
	ticket, err := s.ticketRepo.GetTicketById(id)
	if err != nil {
		return dto.MainEventResponse{}, dto.ErrTicketNotFound
	}

	event, err := s.eventRepo.GetByID(ticket.EventID)
	if err != nil {
		return dto.MainEventResponse{}, dto.ErrEventNotFound
	}

	user, err := s.userRepo.GetUserById(ticket.UserID)
	if err != nil {
		return dto.MainEventResponse{}, dto.ErrUserNotFound
	}

	return dto.MainEventResponse{
		ID:        ticket.TicketID,
		Name:      user.Name,
		Email:     user.Email,
		Confirmed: *ticket.PaymentConfirmed,
		CheckedIn: *ticket.CheckedIn,
		EventName: event.Name,
		Price:     event.Price,

		Handphone:    ticket.Handphone,
		Birthdate:    ticket.Birthdate,
		Seat:         ticket.Seat,
		Payment:      ticket.Payment,
		WithKit:      *event.WithKit,
		RegisterDate: ticket.CreatedAt,
	}, nil
}

func (s *mainEventService) GetMainEventCounter(ctx context.Context) (dto.MainEventCounter, error) {
	total, err := s.ticketRepo.CountTotal()
	if err != nil {
		return dto.MainEventCounter{}, err
	}

	confirmed_payments, err := s.ticketRepo.CountConfirmedPayments()
	if err != nil {
		return dto.MainEventCounter{}, err
	}

	checked_ins, err := s.ticketRepo.CountCheckedIns()
	if err != nil {
		return dto.MainEventCounter{}, err
	}

	return dto.MainEventCounter{
		Total:             total,
		ConfirmedPayments: confirmed_payments,
		CheckedIns:        checked_ins,
	}, nil
}
