package usecase

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/trungnghia250/malo-api/service/model"
	"github.com/trungnghia250/malo-api/service/model/dto"
	"github.com/trungnghia250/malo-api/service/repo"
	mail2 "github.com/trungnghia250/malo-api/service/sender/mail"
	"github.com/trungnghia250/malo-api/service/sender/sms"
	"github.com/trungnghia250/malo-api/utils"
	"time"
)

type ICampaignUseCase interface {
	GetCampaignByID(ctx *fiber.Ctx, campaignID string) (*model.Campaign, error)
	DeleteCampaignsByID(ctx *fiber.Ctx, campaignIDs []string) error
	ListCampaign(ctx *fiber.Ctx, req dto.ListCampaignRequest) ([]model.Campaign, error)
	UpdateCampaign(ctx *fiber.Ctx, data *model.Campaign) (*model.Campaign, error)
	CreateCampaign(ctx *fiber.Ctx, data *model.Campaign) (*model.Campaign, error)
}

type campaignUseCase struct {
	repo repo.IRepo
}

func NewCampaignUseCase(repo repo.IRepo) ICampaignUseCase {
	return &campaignUseCase{
		repo: repo,
	}
}

func (c *campaignUseCase) GetCampaignByID(ctx *fiber.Ctx, campaignID string) (*model.Campaign, error) {
	campaign, err := c.repo.NewCampaignRepo().GetCampaignByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return campaign, nil
}

func (c *campaignUseCase) DeleteCampaignsByID(ctx *fiber.Ctx, campaignIDs []string) error {
	err := c.repo.NewCampaignRepo().DeleteCampaignByID(ctx, campaignIDs)

	return err
}

func (c *campaignUseCase) ListCampaign(ctx *fiber.Ctx, req dto.ListCampaignRequest) ([]model.Campaign, error) {
	campaigns, err := c.repo.NewCampaignRepo().ListCampaign(ctx, req)
	if err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (c *campaignUseCase) CreateCampaign(ctx *fiber.Ctx, data *model.Campaign) (*model.Campaign, error) {
	data.CreatedAt = time.Now()
	campaignID, err := c.repo.NewCounterRepo().GetSequenceNextValue(ctx, "campaign_id")
	if err != nil {
		return nil, err
	}
	data.ID = fmt.Sprintf("CP%d", campaignID)
	err = c.repo.NewCampaignRepo().CreateCampaign(ctx, data)
	if err != nil {
		return nil, err
	}

	groups, _ := c.repo.NewCustomerGroupRepo().ListCustomerGroup(ctx, dto.ListCustomerGroupRequest{
		IDs:   data.CustomerGroupIDs,
		Limit: 10,
	})
	var customerIDs []string
	for _, group := range groups {
		for _, customerID := range group.CustomerIDs {
			if !utils.IsStringContains(customerIDs, customerID) {
				customerIDs = append(customerIDs, customerID)
			}
		}
	}

	customers, _ := c.repo.NewCustomerRepo().ListCustomer(ctx, dto.ListCustomerRequest{
		CustomerIDs: customerIDs,
		Limit:       100,
	})

	if data.Channel == "email" {
		var personalizations []*mail.Personalization
		for _, customer := range customers {
			personalization := mail.NewPersonalization()
			personalization.AddTos(mail.NewEmail(customer.CustomerName, customer.Email))
			personalization.Substitutions = map[string]string{
				"%customer_name%": customer.CustomerName,
				"%voucher_code%":  data.VoucherCode,
			}
			personalization.Subject = data.Title
			personalization.SendAt = data.SendAt
			personalizations = append(personalizations, personalization)
		}
		mailOpt := []mail2.OptMessage{
			mail2.WithMessageHTML(data.Message),
		}

		err = mail2.Send(personalizations, mailOpt...)
		if err != nil {
			return nil, err
		}
	}

	if data.Channel == "sms" {
		var receivers []string
		for _, customer := range customers {
			receivers = append(receivers, customer.PhoneNumber)
		}

		smsRequest := sms.Request{
			Receivers: receivers,
			Content:   data.Message,
		}

		if data.SendAt > 0 {
			smsRequest.SendAt = time.Unix(int64(data.SendAt), 0)
		}

		_ = sms.Send(smsRequest)
	}

	return nil, nil
}

func (c *campaignUseCase) UpdateCampaign(ctx *fiber.Ctx, data *model.Campaign) (*model.Campaign, error) {
	data.ModifiedAt = time.Now()

	err := c.repo.NewCampaignRepo().UpdateCampaignByID(ctx, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}