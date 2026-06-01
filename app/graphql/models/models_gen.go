package gqlmodels

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
)

type PublicUser struct {
	PublicID    string           `json:"publicId"`
	Name        string           `json:"name"`
	DisplayName string           `json:"displayName"`
	Role        enums.UserRole   `json:"role"`
	Plan        enums.UserPlan   `json:"plan"`
	Status      enums.UserStatus `json:"status"`
	CreatedAt   time.Time        `json:"createdAt"`
	UserInfo    *PublicUserInfo  `json:"userInfo"`
}

type PublicUserInfo struct {
	AvatarURL          *string          `json:"avatarURL,omitempty"`
	CoverBackgroundURL *string          `json:"coverBackgroundURL,omitempty"`
	Header             *string          `json:"header,omitempty"`
	Introduction       *string          `json:"introduction,omitempty"`
	Gender             enums.UserGender `json:"gender"`
	Country            *enums.Country   `json:"country,omitempty"`
	BirthDate          time.Time        `json:"birthDate"`
}

type SearchUserCursorFields struct {
	PublicID string `json:"publicId"`
}

type SearchSortOrder string

const (
	SearchSortOrderAsc  SearchSortOrder = "ASC"
	SearchSortOrderDesc SearchSortOrder = "DESC"
)

func (e SearchSortOrder) String() string {
	return string(e)
}

type SearchUserSortBy string

const (
	SearchUserSortByRelevance  SearchUserSortBy = "RELEVANCE"
	SearchUserSortByName       SearchUserSortBy = "NAME"
	SearchUserSortByLastActive SearchUserSortBy = "LAST_ACTIVE"
	SearchUserSortByCreatedAt  SearchUserSortBy = "CREATED_AT"
)

type SearchUserFilters struct {
	Role      *enums.UserRole   `json:"role,omitempty"`
	Plan      *enums.UserPlan   `json:"plan,omitempty"`
	Status    *enums.UserStatus `json:"status,omitempty"`
	HasAvatar *bool             `json:"hasAvatar,omitempty"`
	Country   *enums.Country    `json:"country,omitempty"`
	IsOnline  *bool             `json:"isOnline,omitempty"`
}

type SearchUserInput struct {
	Query     string             `json:"query"`
	After     *string            `json:"after,omitempty"`
	First     *int32             `json:"first,omitempty"`
	Filters   *SearchUserFilters `json:"filters,omitempty"`
	SortBy    *SearchUserSortBy  `json:"sortBy,omitempty"`
	SortOrder *SearchSortOrder   `json:"sortOrder,omitempty"`
}

type SearchPageInfo struct {
	HasNextPage              bool    `json:"hasNextPage"`
	HasPreviousPage          bool    `json:"hasPreviousPage"`
	StartEncodedSearchCursor *string `json:"startEncodedSearchCursor,omitempty"`
	EndEncodedSearchCursor   *string `json:"endEncodedSearchCursor,omitempty"`
}

type SearchUserEdge struct {
	EncodedSearchCursor string      `json:"encodedSearchCursor"`
	Node                *PublicUser `json:"node"`
}

type SearchUserConnection struct {
	SearchEdges    []*SearchUserEdge `json:"searchEdges"`
	SearchPageInfo *SearchPageInfo   `json:"searchPageInfo"`
	TotalCount     int32             `json:"totalCount"`
	SearchTime     float64           `json:"searchTime"`
}
