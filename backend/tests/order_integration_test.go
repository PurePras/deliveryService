//go:build integration

package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PurePras/shri-ram-service/backend/internal/apperror"
	"github.com/PurePras/shri-ram-service/backend/internal/dto"
	"github.com/PurePras/shri-ram-service/backend/internal/model"
	"github.com/PurePras/shri-ram-service/backend/internal/repository"
	"github.com/PurePras/shri-ram-service/backend/internal/service"
)

// TestOrderFlow_StockDecrementsAndRollsBackOnFailure is the highest-risk path in the
// app: money and stock must move together or not at all. It places an order for
// exactly the available stock (must succeed, draining it to zero), then places a
// second order against the same now-empty product (must be rejected) and checks
// stock is still zero — not negative, not partially decremented — proving
// OrderRepository.Create's transaction actually rolled back rather than the failure
// simply being caught after partial writes.
func TestOrderFlow_StockDecrementsAndRollsBackOnFailure(t *testing.T) {
	ctx := context.Background()
	fx := newOrderFixture(t, ctx, "10")
	orderSvc := service.NewOrderService(repository.NewOrderRepository(testPool))
	deliveryDate := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")

	firstOrder, err := orderSvc.Create(ctx, fx.userID, dto.CreateOrderRequest{
		DeliveryAreaID:  fx.areaID,
		DeliverySlotID:  fx.slotID,
		DeliveryDate:    deliveryDate,
		DeliveryAddress: "123 Test Street",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: fx.productID, Quantity: "10"},
		},
	})
	if err != nil {
		t.Fatalf("Create (order for the full stock): %v", err)
	}
	if len(firstOrder.Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(firstOrder.Items))
	}
	if firstOrder.Status != model.OrderStatusPending {
		t.Fatalf("status = %q, want %q", firstOrder.Status, model.OrderStatusPending)
	}
	if firstOrder.TotalAmount != "100.00" {
		t.Fatalf("total_amount = %q, want %q (10 units at 10.00)", firstOrder.TotalAmount, "100.00")
	}

	stockAfterFirst := fx.currentStock(t, ctx)
	if stockAfterFirst != "0.00" {
		t.Fatalf("stock after the first order = %q, want %q", stockAfterFirst, "0.00")
	}

	_, err = orderSvc.Create(ctx, fx.userID, dto.CreateOrderRequest{
		DeliveryAreaID:  fx.areaID,
		DeliverySlotID:  fx.slotID,
		DeliveryDate:    deliveryDate,
		DeliveryAddress: "123 Test Street",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: fx.productID, Quantity: "1"},
		},
	})
	if err == nil {
		t.Fatal("expected an error creating a second order against an out-of-stock product")
	}
	if !errors.Is(err, apperror.ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}

	stockAfterSecond := fx.currentStock(t, ctx)
	if stockAfterSecond != stockAfterFirst {
		t.Fatalf("stock changed after a failed order: was %q, now %q — the transaction did not roll back cleanly",
			stockAfterFirst, stockAfterSecond)
	}
}

// TestOrderFlow_InsufficientStockLeavesStockUntouched is the same rollback guarantee
// from a cold start: a single order that asks for more than what's in stock must be
// rejected in full, not partially fulfilled.
func TestOrderFlow_InsufficientStockLeavesStockUntouched(t *testing.T) {
	ctx := context.Background()
	fx := newOrderFixture(t, ctx, "5")
	orderSvc := service.NewOrderService(repository.NewOrderRepository(testPool))
	deliveryDate := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")

	_, err := orderSvc.Create(ctx, fx.userID, dto.CreateOrderRequest{
		DeliveryAreaID:  fx.areaID,
		DeliverySlotID:  fx.slotID,
		DeliveryDate:    deliveryDate,
		DeliveryAddress: "123 Test Street",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: fx.productID, Quantity: "6"},
		},
	})
	if err == nil {
		t.Fatal("expected an error creating an order for more than the available stock")
	}
	if !errors.Is(err, apperror.ErrValidation) {
		t.Fatalf("expected ErrValidation, got: %v", err)
	}

	if stock := fx.currentStock(t, ctx); stock != "5.00" {
		t.Fatalf("stock after a rejected order = %q, want unchanged %q", stock, "5.00")
	}
}

// orderFixture is a minimal, self-contained set of rows an order can legally
// reference: a user, a category+product with known stock, a delivery area, and a
// delivery slot. Every field is randomized so repeated runs never collide, and
// t.Cleanup tears everything down in FK-safe order (orders/order_items first, via
// ON DELETE CASCADE from orders, down to the independent parent rows last).
type orderFixture struct {
	userID     string
	categoryID string
	productID  string
	areaID     string
	slotID     string
}

func newOrderFixture(t *testing.T, ctx context.Context, stock string) *orderFixture {
	t.Helper()
	suffix := randomSuffix()

	passwordHash := "test-hash"
	user, err := repository.NewUserRepository(testPool).Create(ctx, &model.User{
		Name:         "Order Fixture User " + suffix,
		Phone:        randomTestPhone(),
		PasswordHash: &passwordHash,
	})
	if err != nil {
		t.Fatalf("fixture: create user: %v", err)
	}

	category, err := repository.NewCategoryRepository(testPool).Create(ctx, &model.Category{
		Name:     "Fixture Category " + suffix,
		Slug:     "fixture-category-" + suffix,
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("fixture: create category: %v", err)
	}

	product, err := repository.NewProductRepository(testPool).Create(ctx, &model.Product{
		CategoryID:    category.ID,
		Name:          "Fixture Product " + suffix,
		Slug:          "fixture-product-" + suffix,
		Unit:          "kg",
		Price:         "10.00",
		StockQuantity: stock,
		IsAvailable:   true,
	})
	if err != nil {
		t.Fatalf("fixture: create product: %v", err)
	}

	area, err := repository.NewDeliveryAreaRepository(testPool).Create(ctx, &model.DeliveryArea{
		Name:     "Fixture Area " + suffix,
		City:     "Test City",
		Pincode:  suffix,
		IsActive: true,
	})
	if err != nil {
		t.Fatalf("fixture: create delivery area: %v", err)
	}

	slot, err := repository.NewDeliverySlotRepository(testPool).Create(ctx, &model.DeliverySlot{
		Label:     "Fixture Slot " + suffix,
		StartTime: "09:00:00",
		EndTime:   "12:00:00",
		IsActive:  true,
	})
	if err != nil {
		t.Fatalf("fixture: create delivery slot: %v", err)
	}

	fx := &orderFixture{
		userID:     user.ID,
		categoryID: category.ID,
		productID:  product.ID,
		areaID:     area.ID,
		slotID:     slot.ID,
	}
	t.Cleanup(func() { fx.cleanup(t) })
	return fx
}

func (fx *orderFixture) currentStock(t *testing.T, ctx context.Context) string {
	t.Helper()
	var stock string
	if err := testPool.QueryRow(ctx, `SELECT stock_quantity::text FROM products WHERE id = $1`, fx.productID).Scan(&stock); err != nil {
		t.Fatalf("read current stock: %v", err)
	}
	return stock
}

func (fx *orderFixture) cleanup(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	// orders cascades to order_items; everything after it is a RESTRICT-protected
	// parent that only becomes deletable once those are gone.
	if _, err := testPool.Exec(ctx, `DELETE FROM orders WHERE user_id = $1`, fx.userID); err != nil {
		t.Logf("cleanup: delete orders for user %s: %v", fx.userID, err)
	}
	if _, err := testPool.Exec(ctx, `DELETE FROM products WHERE id = $1`, fx.productID); err != nil {
		t.Logf("cleanup: delete product %s: %v", fx.productID, err)
	}
	if _, err := testPool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, fx.categoryID); err != nil {
		t.Logf("cleanup: delete category %s: %v", fx.categoryID, err)
	}
	if _, err := testPool.Exec(ctx, `DELETE FROM delivery_areas WHERE id = $1`, fx.areaID); err != nil {
		t.Logf("cleanup: delete delivery area %s: %v", fx.areaID, err)
	}
	if _, err := testPool.Exec(ctx, `DELETE FROM delivery_slots WHERE id = $1`, fx.slotID); err != nil {
		t.Logf("cleanup: delete delivery slot %s: %v", fx.slotID, err)
	}
	if _, err := testPool.Exec(ctx, `DELETE FROM users WHERE id = $1`, fx.userID); err != nil {
		t.Logf("cleanup: delete user %s: %v", fx.userID, err)
	}
}
