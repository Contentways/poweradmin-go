// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"fmt"

	"contentways.dev/contentways/poweradmin-go/poweradmin/schema"
)

// Record represents a DNS record in Poweradmin.
type Record struct {
	ID       int64
	ZoneID   int64
	Name     string
	Type     string
	Content  string
	TTL      int
	Priority int
	Disabled bool
	Auth     bool // authoritative flag (managed by PowerDNS)
}

// RecordCreateOpts configures a record creation request.
type RecordCreateOpts struct {
	Name      string
	Type      string
	Content   string
	TTL       int
	Priority  int
	Disabled  bool
	CreatePTR bool
}

// RecordUpdateOpts configures a record update request.
// Pointer fields are only sent when non-nil.
type RecordUpdateOpts struct {
	Name     string
	Type     string
	Content  string
	TTL      *int
	Priority *int
	Disabled *bool
}

// BulkRecordOperation represents a single operation in a bulk records request.
// Action is one of "create", "update", "delete". RecordID is required for
// update/delete; Name/Type/Content are required for create.
type BulkRecordOperation struct {
	Action   string
	RecordID int
	Name     string
	Type     string
	Content  string
	TTL      int
	Priority int
	Disabled bool
}

// BulkRecordsResult is the outcome of a bulk records call.
type BulkRecordsResult struct {
	SuccessCount int
	FailureCount int
	Errors       []string
}

// RecordClient provides access to the record-related Poweradmin API endpoints.
type RecordClient struct {
	client *Client
}

// GetByID returns a single [Record] by zone and record ID.
func (r *RecordClient) GetByID(ctx context.Context, zoneID int, recordID int64) (*Record, *Response, error) {
	var result schema.RecordResponse
	resp, err := r.client.get(ctx, fmt.Sprintf("zones/%d/records/%d", zoneID, recordID), &result)
	if err != nil {
		return nil, resp, err
	}
	rec := RecordFromSchema(result.Record)
	return &rec, resp, nil
}

// RecordListOpts adds a record-type filter on top of the standard pagination options.
type RecordListOpts struct {
	ListOpts
	Type string // optional record-type filter (e.g. "A", "MX")
}

// List returns one page of [Record]s for the given zone.
// Use [RecordListOpts] with Type set to filter by record type.
func (r *RecordClient) List(ctx context.Context, zoneID int, opts RecordListOpts) ([]*Record, *Response, error) {
	v := opts.ListOpts.values()
	if opts.Type != "" {
		v.Set("type", opts.Type)
	}
	path := appendQuery(fmt.Sprintf("zones/%d/records", zoneID), v)
	var result []schema.Record
	resp, err := r.client.get(ctx, path, &result)
	if err != nil {
		return nil, resp, err
	}
	records := make([]*Record, len(result))
	for i, s := range result {
		rec := RecordFromSchema(s)
		records[i] = &rec
	}
	return records, resp, nil
}

// All returns all [Record]s in a zone across all pages.
func (r *RecordClient) All(ctx context.Context, zoneID int) ([]*Record, error) {
	var all []*Record
	opts := RecordListOpts{ListOpts: ListOpts{Page: 1, PerPage: 100}}
	for {
		records, resp, err := r.List(ctx, zoneID, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, records...)
		if resp.Meta.Pagination == nil || resp.Meta.Pagination.Page >= resp.Meta.Pagination.LastPage {
			return all, nil
		}
		opts.Page++
	}
}

// Create creates a new [Record] in the given zone and returns the new ID.
func (r *RecordClient) Create(ctx context.Context, zoneID int, opts RecordCreateOpts) (int64, *Response, error) {
	req := schema.RecordCreateRequest{
		Name:      opts.Name,
		Type:      opts.Type,
		Content:   opts.Content,
		TTL:       opts.TTL,
		Priority:  opts.Priority,
		Disabled:  opts.Disabled,
		CreatePTR: opts.CreatePTR,
	}
	var result schema.RecordResponse
	resp, err := r.client.post(ctx, fmt.Sprintf("zones/%d/records", zoneID), req, &result)
	if err != nil {
		return 0, resp, err
	}
	return result.Record.ID, resp, nil
}

// Update updates an existing [Record] and returns the updated state.
func (r *RecordClient) Update(ctx context.Context, zoneID int, recordID int64, opts RecordUpdateOpts) (*Record, *Response, error) {
	req := schema.RecordUpdateRequest{
		Name:     opts.Name,
		Type:     opts.Type,
		Content:  opts.Content,
		TTL:      opts.TTL,
		Priority: opts.Priority,
		Disabled: opts.Disabled,
	}
	var result schema.RecordResponse
	resp, err := r.client.put(ctx, fmt.Sprintf("zones/%d/records/%d", zoneID, recordID), req, &result)
	if err != nil {
		return nil, resp, err
	}
	rec := RecordFromSchema(result.Record)
	return &rec, resp, nil
}

// Delete deletes the [Record] with the given ID.
func (r *RecordClient) Delete(ctx context.Context, zoneID int, recordID int64) (*Response, error) {
	return r.client.delete(ctx, fmt.Sprintf("zones/%d/records/%d", zoneID, recordID))
}

// Bulk executes multiple record operations atomically against the given zone.
func (r *RecordClient) Bulk(ctx context.Context, zoneID int, ops []BulkRecordOperation) (*BulkRecordsResult, *Response, error) {
	schemaOps := make([]schema.BulkRecordOperation, len(ops))
	for i, op := range ops {
		schemaOps[i] = schema.BulkRecordOperation{
			Action:   op.Action,
			ID:       op.RecordID,
			Name:     op.Name,
			Type:     op.Type,
			Content:  op.Content,
			TTL:      op.TTL,
			Priority: op.Priority,
			Disabled: op.Disabled,
		}
	}
	req := schema.BulkRecordsRequest{Operations: schemaOps}
	var result schema.BulkRecordsResponse
	resp, err := r.client.post(ctx, fmt.Sprintf("zones/%d/records/bulk", zoneID), req, &result)
	if err != nil {
		return nil, resp, err
	}
	return &BulkRecordsResult{
		SuccessCount: result.SuccessCount,
		FailureCount: result.FailureCount,
		Errors:       result.Errors,
	}, resp, nil
}
