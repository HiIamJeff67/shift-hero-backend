package dtos

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/google/uuid"
)

/* ============================== Company Request DTO ============================== */

type CreateCompanyReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			Name        string `json:"name" validate:"required,min=2,max=128"`
			Description string `json:"description" validate:"omitempty,max=1024"`
			Email       string `json:"email" validate:"required,email"`
		},
		any,
	]
}

type GetMyCompaniesReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		any,
		any,
	]
}

type GetCompanyReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type UpdateCompanyReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId uuid.UUID `json:"companyId" validate:"required,uuid4"`
			PartialUpdateDto[struct {
				Name        *string `json:"name" validate:"omitnil,min=2,max=128"`
				Description *string `json:"description" validate:"omitnil,max=1024"`
				Email       *string `json:"email" validate:"omitnil,email"`
			}]
		},
		any,
	]
}

type GetCompanyMembersReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type AddCompanyMemberReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId    uuid.UUID          `json:"companyId" validate:"required,uuid4"`
			UserId       uuid.UUID          `json:"userId" validate:"required,uuid4"`
			EmployeeRole enums.EmployeeRole `json:"employeeRole" validate:"required,isemployeerole"`
		},
		any,
	]
}

type UpdateCompanyMemberReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId    uuid.UUID          `json:"companyId" validate:"required,uuid4"`
			UserId       uuid.UUID          `json:"userId" validate:"required,uuid4"`
			EmployeeRole enums.EmployeeRole `json:"employeeRole" validate:"required,isemployeerole"`
		},
		any,
	]
}

type DeleteCompanyMemberReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId uuid.UUID `json:"companyId" validate:"required,uuid4"`
			UserId    uuid.UUID `json:"userId" validate:"required,uuid4"`
		},
		any,
	]
}

type CreateCompanyJoinRequestReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId     uuid.UUID           `json:"companyId" validate:"required,uuid4"`
			RequestedRole *enums.EmployeeRole `json:"requestedRole" validate:"omitnil,isemployeerole"`
			Note          string              `json:"note" validate:"omitempty,max=1024"`
		},
		any,
	]
}

type GetCompanyJoinRequestsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			Status *enums.CompanyJoinRequestStatus `form:"status" validate:"omitnil,iscompanyjoinrequeststatus"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type ReviewCompanyJoinRequestReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId     uuid.UUID `json:"companyId" validate:"required,uuid4"`
			JoinRequestId uuid.UUID `json:"joinRequestId" validate:"required,uuid4"`
		},
		any,
	]
}

type GetMyCompanyJoinRequestsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		any,
		any,
	]
}

/* ============================== Company Response DTO ============================== */

type CompanyMemberResDto struct {
	UserId       uuid.UUID          `json:"userId"`
	Name         string             `json:"name"`
	DisplayName  string             `json:"displayName"`
	Email        string             `json:"email"`
	EmployeeRole enums.EmployeeRole `json:"employeeRole"`
}

type CompanyResDto struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Email       string    `json:"email"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CompanyJoinRequestResDto struct {
	Id               uuid.UUID                      `json:"id"`
	CompanyId        uuid.UUID                      `json:"companyId"`
	CompanyName      string                         `json:"companyName"`
	RequesterUserId  uuid.UUID                      `json:"requesterUserId"`
	RequesterName    string                         `json:"requesterName"`
	RequesterEmail   string                         `json:"requesterEmail"`
	RequestedRole    enums.EmployeeRole             `json:"requestedRole"`
	Note             string                         `json:"note"`
	Status           enums.CompanyJoinRequestStatus `json:"status"`
	ReviewedByUserId *uuid.UUID                     `json:"reviewedByUserId"`
	ReviewedAt       *time.Time                     `json:"reviewedAt"`
	CreatedAt        time.Time                      `json:"createdAt"`
	UpdatedAt        time.Time                      `json:"updatedAt"`
}
