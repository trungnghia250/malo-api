package dto

import (
	"github.com/trungnghia250/malo-api/service/model"
	"mime/multipart"
	"time"
)

type GetCustomerByIDRequest struct {
	CustomerID string `json:"customer_id" query:"customer_id"`
}

type GetCustomerGroupByIDRequest struct {
	ID string `json:"id" query:"id"`
}

type GetCustomerByIDResponse struct {
	CustomerID   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	Gender       string `json:"gender"`
	PhoneNumber  string `json:"phone_number"`
}

type Customer struct {
	CustomerID     string    `bson:"customer_id" json:"customer_id"`
	CustomerName   string    `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	Gender         string    `bson:"gender,omitempty" json:"gender,omitempty"`
	PhoneNumber    string    `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Email          string    `bson:"email,omitempty" json:"email,omitempty"`
	Address        string    `json:"address,omitempty" bson:"address,omitempty"`
	Province       string    `json:"province,omitempty" bson:"province,omitempty"`
	DateOfBirth    time.Time `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	CustomerType   string    `json:"customer_type,omitempty" bson:"customer_type,omitempty"`
	CustomerSource string    `json:"customer_source,omitempty" bson:"customer_source,omitempty"`
	Tags           []string  `json:"tags" bson:"tags"`
	Note           string    `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt     time.Time `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy     string    `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	RewardPoint    int32     `json:"reward_point,omitempty" bson:"reward_point,omitempty"`
	RankPoint      int32     `json:"rank_point,omitempty" bson:"rank_point,omitempty"`
}

type ListCustomerRequest struct {
	Limit        int32    `json:"limit,omitempty" query:"limit,omitempty"`
	Offset       int32    `json:"offset,omitempty" query:"offset,omitempty"`
	CustomerName []string `json:"customer_name,omitempty" query:"customer_name,omitempty"`
	Phone        string   `json:"phone,omitempty" query:"phone,omitempty"`
	Address      string   `json:"address,omitempty" query:"address,omitempty"`
	CustomerType []string `json:"customer_type,omitempty" query:"customer_type,omitempty"`
	Tags         []string `json:"tags,omitempty" query:"tags,omitempty"`
	Gender       string   `json:"gender,omitempty" query:"gender,omitempty"`
	Email        string   `json:"email,omitempty" query:"email,omitempty"`
	Source       []string `json:"source,omitempty" query:"source,omitempty"`
	CustomerIDs  []string `json:"customer_ids,omitempty" query:"customer_ids,omitempty"`
	ExceptIDs    []string `json:"except_ids,omitempty" query:"except_ids,omitempty"`
	CreatedAt    []int32  `json:"created_at,omitempty" query:"created_at,omitempty"`
}

type ListCustomerResponse struct {
	Count int32            `json:"count"`
	Data  []model.Customer `json:"data"`
}

type DeleteCustomersRequest struct {
	CustomerIDs []string `json:"customer_ids" query:"customer_ids"`
}

type DeleteCustomerGroupsRequest struct {
	IDs []string `json:"ids" query:"ids"`
}

type UpdateListCustomerRequest struct {
	CustomerIDs []string `json:"customer_ids" query:"customer_ids"`
	Tags        []string `json:"tags" query:"tags"`
}

type ExportCustomerRequest struct {
	CustomerIDs []string `json:"customer_ids,omitempty" query:"customer_ids,omitempty"`
	Filter      string   `json:"filter,omitempty" query:"filter,omitempty"`
}

type ListCustomerGroupRequest struct {
	Limit     int32    `json:"limit,omitempty"`
	Offset    int32    `json:"offset,omitempty"`
	IDs       []string `json:"ids,omitempty"`
	Name      []string `json:"name,omitempty"`
	CreatedAt []int32  `json:"created_at,omitempty"`
}

type ListCustomerGroupResponse struct {
	Count int32                 `json:"count"`
	Data  []model.CustomerGroup `json:"data"`
}

type GetCustomerGroupResponse struct {
	ID         string           `json:"_id,omitempty" bson:"_id,omitempty"`
	GroupName  string           `json:"group_name,omitempty" bson:"group_name,omitempty"`
	Note       string           `json:"note,omitempty" bson:"note,omitempty"`
	CreatedAt  time.Time        `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ModifiedAt time.Time        `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	ModifiedBy string           `json:"modified_by,omitempty" bson:"modified_by,omitempty"`
	TotalCount int32            `json:"totalCount,omitempty" bson:"totalCount,omitempty"`
	Customers  []model.Customer `json:"customers,omitempty"`
}

type CreateCustomerGroup struct {
	Data   *model.CustomerGroup `json:"data,omitempty"`
	Filter ListCustomerRequest  `json:"filter,omitempty"`
}

type ImportCustomerRequest struct {
	File *multipart.FileHeader `form:"file"`
}

type ImportCustomerResponse struct {
	Scan    int32            `json:"scan"`
	Success int32            `json:"success"`
	Insert  int32            `json:"insert"`
	Update  int32            `json:"update"`
	Ignore  int32            `json:"ignore"`
	Data    []model.Customer `json:"data"`
}
