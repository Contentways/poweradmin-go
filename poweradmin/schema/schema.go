// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT

// Package schema contains the raw JSON types that map 1:1 to the Poweradmin API.
// These types are not meant to be used directly by consumers of the SDK.
// Use the domain types in the parent package instead.
//
// Shapes follow Poweradmin API v2 OpenAPI spec:
//   - Envelope: {success, message, data, pagination, meta, error}
//   - Pagination lives at envelope level, not inside data.
//   - List endpoints return data as a flat array.
//   - Single-resource endpoints typically wrap the object under data.<resource>
//     (e.g. data.zone, data.record); a few return the object directly under data
//     (notably RRSet single).
package schema

import "encoding/json"

// APIResponse represents the Poweradmin API envelope. Data is RawMessage so
// callers can unmarshal it into the concrete shape without a marshal roundtrip.
type APIResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data,omitempty"`
	Pagination *Pagination     `json:"pagination,omitempty"`
	Meta       *APIMeta        `json:"meta,omitempty"`
	Error      *ErrorPayload   `json:"error,omitempty"`
}

type APIMeta struct {
	Timestamp string `json:"timestamp,omitempty"`
}

type ErrorPayload struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

// ── Zone ─────────────────────────────────────────────────────────────────────

type Zone struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Masters      string `json:"masters,omitempty"`
	Account      string `json:"account,omitempty"`
	Description  string `json:"description,omitempty"`
	SOASerial    int    `json:"soa_serial,omitempty"`
	DNSSECSigned bool   `json:"dnssec_signed,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type ZoneResponse struct {
	Zone Zone `json:"zone"`
}

// ZoneListResponse is the wrapper used by /v2/zones — alone among the list
// endpoints in this API, this one nests the array under a "zones" key inside
// data instead of returning a flat array.
type ZoneListResponse struct {
	Zones []Zone `json:"zones"`
}

type ZoneCreateRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Masters     string `json:"masters,omitempty"`
	Account     string `json:"account,omitempty"`
	Description string `json:"description,omitempty"`
	Template    string `json:"template,omitempty"`
}

type ZoneCreateResponse struct {
	ZoneID int `json:"zone_id"`
}

type ZoneUpdateRequest struct {
	Type        *string `json:"type,omitempty"`
	Masters     *string `json:"masters,omitempty"`
	Account     *string `json:"account,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ── Record ───────────────────────────────────────────────────────────────────

type Record struct {
	ID       int64  `json:"id,omitempty"`
	ZoneID   int64  `json:"zone_id,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority,omitempty"`
	Disabled bool   `json:"disabled"`
	Auth     bool   `json:"auth,omitempty"`
}

// RecordResponse is used for both single-record GET and the POST create
// response. The API wraps the record under data.record in both cases.
type RecordResponse struct {
	Record Record `json:"record"`
}

type RecordCreateRequest struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Priority  int    `json:"priority,omitempty"`
	Disabled  bool   `json:"disabled,omitempty"`
	CreatePTR bool   `json:"create_ptr,omitempty"`
}

type RecordUpdateRequest struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Content  string `json:"content,omitempty"`
	TTL      *int   `json:"ttl,omitempty"`
	Priority *int   `json:"priority,omitempty"`
	Disabled *bool  `json:"disabled,omitempty"`
}

type BulkRecordOperation struct {
	Action   string `json:"action"`
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Content  string `json:"content,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

type BulkRecordsRequest struct {
	Operations []BulkRecordOperation `json:"operations"`
}

type BulkRecordsResponse struct {
	SuccessCount int      `json:"success_count,omitempty"`
	FailureCount int      `json:"failure_count,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// ── RRSet ────────────────────────────────────────────────────────────────────
// RRSet (Resource Record Set) groups records sharing a name+type. The API
// identifies them by name+type within a zone, not by ID. Records are objects,
// not strings.

type RRSet struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	TTL     int           `json:"ttl"`
	Records []RRSetRecord `json:"records"`
}

type RRSetRecord struct {
	Content  string `json:"content"`
	Priority int    `json:"priority,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// RRSetUpsertRequest is sent via PUT to create-or-replace a full RRSet.
type RRSetUpsertRequest struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	TTL     int           `json:"ttl"`
	Records []RRSetRecord `json:"records"`
}

// ── User ─────────────────────────────────────────────────────────────────────

type User struct {
	UserID      int      `json:"user_id,omitempty"`
	Username    string   `json:"username"`
	Fullname    string   `json:"fullname"`
	Email       string   `json:"email"`
	Description string   `json:"description,omitempty"`
	Active      bool     `json:"active"`
	PermTempl   int      `json:"perm_templ,omitempty"`
	UseLdap     bool     `json:"use_ldap,omitempty"`
	IsAdmin     bool     `json:"is_admin,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	ZoneCount   int      `json:"zone_count,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UserCreateResponse struct {
	UserID int `json:"user_id"`
}

type UserCreateRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active"`
	PermTempl   int    `json:"perm_templ,omitempty"`
	UseLdap     bool   `json:"use_ldap,omitempty"`
}

type UserUpdateRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Fullname    string `json:"fullname,omitempty"`
	Email       string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
	// Active must be sent explicitly even for false to take effect — see API docs.
	Active    *bool `json:"active,omitempty"`
	PermTempl int   `json:"perm_templ,omitempty"`
	UseLdap   *bool `json:"use_ldap,omitempty"`
}

// UserPatchRequest is sent via PATCH /v2/users/{id} for partial updates such as
// assigning a permission template.
type UserPatchRequest struct {
	PermTempl int `json:"perm_templ,omitempty"`
}

// ── Permission ───────────────────────────────────────────────────────────────

type Permission struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Descr string `json:"descr"`
}

type PermissionResponse struct {
	Permission Permission `json:"permission"`
}

// ── Permission Template ──────────────────────────────────────────────────────

type PermissionTemplate struct {
	ID           int          `json:"id,omitempty"`
	Name         string       `json:"name"`
	Descr        string       `json:"descr"`
	TemplateType string       `json:"template_type,omitempty"`
	Permissions  []Permission `json:"permissions,omitempty"`
}

type PermissionTemplateResponse struct {
	Template PermissionTemplate `json:"template"`
}

type PermissionTemplateRequest struct {
	Name         string `json:"name"`
	Descr        string `json:"descr"`
	TemplateType string `json:"template_type,omitempty"`
	Permissions  []int  `json:"permissions,omitempty"`
}

// ── Group ────────────────────────────────────────────────────────────────────

type Group struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PermTemplID int    `json:"perm_templ_id,omitempty"`
	MemberCount int    `json:"member_count,omitempty"`
	ZoneCount   int    `json:"zone_count,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type GroupCreateResponse struct {
	Group Group `json:"group"`
}

type GroupCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	PermTemplID int    `json:"perm_templ_id"`
}

type GroupUpdateRequest struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	PermTemplID int     `json:"perm_templ_id,omitempty"`
}

type GroupMember struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
	JoinedAt string `json:"joined_at,omitempty"`
}

type GroupMemberRequest struct {
	UserID int `json:"user_id"`
}

type GroupZone struct {
	ZoneID    int    `json:"zone_id"`
	ZoneName  string `json:"zone_name"`
	ZoneType  string `json:"zone_type"`
	CreatedAt string `json:"created_at,omitempty"`
}

type GroupZoneRequest struct {
	ZoneID int `json:"zone_id"`
}

// ── Zone Owner ───────────────────────────────────────────────────────────────

type ZoneOwner struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"fullname,omitempty"`
}

// ZoneOwnerAddRequest supports single or batch add.
type ZoneOwnerAddRequest struct {
	UserID  int   `json:"user_id,omitempty"`
	UserIDs []int `json:"user_ids,omitempty"`
}

// ── Zone Template ────────────────────────────────────────────────────────────

type ZoneTemplate struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       int    `json:"owner,omitempty"`
	IsGlobal    bool   `json:"is_global,omitempty"`
	ZonesLinked int    `json:"zones_linked,omitempty"`
}

type ZoneTemplateResponse struct {
	Template ZoneTemplate `json:"template"`
}

type ZoneTemplateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsGlobal    bool   `json:"is_global,omitempty"`
}

type ZoneTemplateRecord struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

type ZoneTemplateRecordResponse struct {
	Record ZoneTemplateRecord `json:"record"`
}

type ZoneTemplateRecordRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
}
