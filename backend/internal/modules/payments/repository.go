package payments

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionPayments = "payments"

// Repository defines persistence operations for the payments collection.
type Repository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, p *Payment) (*Payment, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*Payment, error)
	List(ctx context.Context, filter PaymentFilter) ([]Payment, error)
	GetByStudentID(ctx context.Context, studentID primitive.ObjectID) ([]Payment, error)
	GetPendingByStudentID(ctx context.Context, studentID primitive.ObjectID) (*Payment, error)
	Update(ctx context.Context, p *Payment) (*Payment, error)
	SoftDelete(ctx context.Context, id primitive.ObjectID) error
	GetFinancialSummary(ctx context.Context) (*FinancialSummary, error)
}

type repository struct {
	coll *mongo.Collection
}

// NewRepository creates a new payments repository with the provided MongoDB database.
func NewRepository(db *mongo.Database) Repository {
	if db == nil {
		return &repository{coll: nil}
	}
	return &repository{
		coll: db.Collection(CollectionPayments),
	}
}

// EnsureIndexes creates required indexes for fast queries and soft delete filtration.
func (r *repository) EnsureIndexes(ctx context.Context) error {
	if r.coll == nil {
		return errors.New("database collection is nil")
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetName("idx_payments_student_id"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_payments_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_payments_created_at"),
		},
		{
			Keys:    bson.D{{Key: "archived_at", Value: 1}},
			Options: options.Index().SetName("idx_payments_archived_at"),
		},
	}

	_, err := r.coll.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create inserts a new payment record.
func (r *repository) Create(ctx context.Context, p *Payment) (*Payment, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}

	_, err := r.coll.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetByID retrieves an active payment by its ObjectID.
func (r *repository) GetByID(ctx context.Context, id primitive.ObjectID) (*Payment, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	filter := bson.M{
		"_id":         id,
		"archived_at": bson.M{"$exists": false},
	}

	var p Payment
	err := r.coll.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

// List retrieves active payments matching the given query filters.
func (r *repository) List(ctx context.Context, filter PaymentFilter) ([]Payment, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	query := bson.M{
		"archived_at": bson.M{"$exists": false},
	}

	if filter.StudentID != "" {
		stID, err := primitive.ObjectIDFromHex(filter.StudentID)
		if err == nil {
			query["student_id"] = stID
		}
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.Search != "" {
		query["student_name"] = bson.M{
			"$regex":   primitive.Regex{Pattern: filter.Search, Options: "i"},
		}
	}

	if filter.From != "" || filter.To != "" {
		dateFilter := bson.M{}
		if filter.From != "" {
			dateFilter["$gte"] = filter.From
		}
		if filter.To != "" {
			dateFilter["$lte"] = filter.To
		}
		query["payment_date"] = dateFilter
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	if payments == nil {
		payments = []Payment{}
	}
	return payments, nil
}

// GetByStudentID retrieves all active payments for a student.
func (r *repository) GetByStudentID(ctx context.Context, studentID primitive.ObjectID) ([]Payment, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	filter := bson.M{
		"student_id":  studentID,
		"archived_at": bson.M{"$exists": false},
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []Payment
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	if list == nil {
		list = []Payment{}
	}
	return list, nil
}

// GetPendingByStudentID finds the first active pending payment for a student.
func (r *repository) GetPendingByStudentID(ctx context.Context, studentID primitive.ObjectID) (*Payment, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	filter := bson.M{
		"student_id":  studentID,
		"status":      StatusPending,
		"archived_at": bson.M{"$exists": false},
	}

	var p Payment
	err := r.coll.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // No pending payment found
		}
		return nil, err
	}
	return &p, nil
}

// Update updates an existing payment in MongoDB.
func (r *repository) Update(ctx context.Context, p *Payment) (*Payment, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	p.UpdatedAt = time.Now().UTC()
	filter := bson.M{
		"_id":         p.ID,
		"archived_at": bson.M{"$exists": false},
	}

	res, err := r.coll.ReplaceOne(ctx, filter, p)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, ErrPaymentNotFound
	}

	return p, nil
}

// SoftDelete marks a payment as archived without destroying financial history.
func (r *repository) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	if r.coll == nil {
		return errors.New("database collection is nil")
	}

	now := time.Now().UTC()
	filter := bson.M{
		"_id":         id,
		"archived_at": bson.M{"$exists": false},
	}
	update := bson.M{
		"$set": bson.M{
			"archived_at": now,
			"updated_at":  now,
		},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrPaymentNotFound
	}
	return nil
}

// GetFinancialSummary aggregates total amounts and transaction counts from active payments.
func (r *repository) GetFinancialSummary(ctx context.Context) (*FinancialSummary, error) {
	if r.coll == nil {
		return nil, errors.New("database collection is nil")
	}

	matchStage := bson.D{{Key: "$match", Value: bson.M{
		"archived_at": bson.M{"$exists": false},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id": nil,
		"totalCount": bson.M{"$sum": 1},
		"paidCount": bson.M{
			"$sum": bson.M{
				"$cond": bson.A{bson.M{"$eq": bson.A{"$status", StatusPaid}}, 1, 0},
			},
		},
		"pendingCount": bson.M{
			"$sum": bson.M{
				"$cond": bson.A{bson.M{"$eq": bson.A{"$status", StatusPending}}, 1, 0},
			},
		},
		"totalPaidAmount": bson.M{
			"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$status", StatusPaid}},
						bson.M{"$ne": bson.A{"$amount", nil}},
					}},
					"$amount",
					0,
				},
			},
		},
		"totalPendingAmount": bson.M{
			"$sum": bson.M{
				"$cond": bson.A{
					bson.M{"$and": bson.A{
						bson.M{"$eq": bson.A{"$status", StatusPending}},
						bson.M{"$ne": bson.A{"$amount", nil}},
					}},
					"$amount",
					0,
				},
			},
		},
	}}}

	cursor, err := r.coll.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []FinancialSummary
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &FinancialSummary{
			TotalPaidAmount:    0,
			TotalPendingAmount: 0,
			PaidCount:          0,
			PendingCount:       0,
			TotalCount:         0,
		}, nil
	}

	return &results[0], nil
}
