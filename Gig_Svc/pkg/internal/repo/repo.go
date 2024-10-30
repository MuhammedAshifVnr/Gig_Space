package repo

import (
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"gorm.io/gorm"
)

type GigRepo struct {
	DB *gorm.DB
}

func NewGigRepository(db *gorm.DB) *GigRepo {
	return &GigRepo{
		DB: db,
	}
}

func (r *GigRepo) CreateGgi(gig model.Gig) (model.Gig, error) {
	err := r.DB.Create(&gig).Error
	if err != nil {
		return model.Gig{}, err
	}
	return gig, nil
}

func (r *GigRepo) GetGigsByFreelancerID(freelancerID uint) ([]model.Gig, error) {
	var gigs []model.Gig
	query := r.DB.Where("freelancer_id = ?", freelancerID).Preload("Images")

	err := query.Find(&gigs).Error
	if err != nil {
		return nil, err
	}
	return gigs, nil
}

func (r *GigRepo) GetGigByID(Id uint) (model.Gig, error) {
	var gig model.Gig
	// query := `select * from gigs where id =?`
	// err := r.DB.Raw(query, Id).Scan(&gig)
	err := r.DB.Preload("Images").First(&gig, Id).Error
	return gig, err
}

func (r *GigRepo) UpdateGig(gig model.Gig) error {
	err := r.DB.Save(&gig).Error
	return err
}

func (r *GigRepo) DeleteImages(id uint) error {
	fmt.Println("id = ", id)
	query := `DELETE FROM images WHERE gig_id = ?`
	err := r.DB.Exec(query, id).Error
	return err
}

func (r *GigRepo) DeleteGig(id, user_id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Exec(`DELETE FROM images WHERE gig_id = ?`, id).Error; err != nil {
			return err
		}

		if err := tx.Exec(`DELETE FROM gigs WHERE id = ? AND freelancer_id = ?`, id, user_id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GigRepo) CreateOrder(data model.Order) error {
	err := r.DB.Create(&data).Error
	return err
}

func (r *GigRepo) GetOrders(clientID uint) ([]*proto.Order, error) {
	query := `SELECT * FROM ORDERS WHERE client_id = ?`
	var orders []*proto.Order
	var res []model.Order
	err := r.DB.Raw(query, clientID).Scan(&res)
	if err != nil {
		return nil, err.Error
	}
	for _, val := range res {
		order := &proto.Order{
			OrderId:      val.OrderID,
			GigId:        uint32(val.GigID),
			FreelancerId: uint32(val.FreelancerID),
			// PaymentId:    val.PaymentID,
			Amount: int64(val.Amount),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *GigRepo) CreateQuote(Quote model.Quote) error {
	err := r.DB.Create(&Quote)
	return err.Error
}

func (r *GigRepo) GetAllQuotes(freelancerID uint) ([]*proto.Quote, error) {
	var result []model.Quote
	err := r.DB.Preload("Gig").
		Joins("JOIN gigs ON gigs.id = quotes.gig_id").
		Where("gigs.freelancer_id = ?", freelancerID).
		Find(&result)
		//fmt.Println(result)
	var quotes []*proto.Quote
	for _, val := range result {
		//fmt.Println(val.Gig.FreelancerID)
		quotes = append(quotes, &proto.Quote{
			GigId:        val.GigId,
			ClientId:     val.ClientId,
			Describe:     val.Describe,
			Price:        val.Price,
			DeliveryDays: int64(val.DeliveryDays),
			FreelancerId: uint64(val.Gig.FreelancerID),
		})
	}
	return quotes, err.Error
}

func (r *GigRepo) CreateCustomGig(gig model.CustomGig) error {
	err := r.DB.Create(&gig)
	return err.Error
}

func (r *GigRepo) GetAllOffers(clientID uint) ([]*proto.CreateOfferReq, error) {
	var offers []*proto.CreateOfferReq
	err := r.DB.Raw(`
	SELECT * 
	FROM custom_gigs
	WHERE client_id = ?
`, clientID).Scan(&offers)
	return offers, err.Error
}

func (r *GigRepo) GetCustomGig(ID uint) (model.CustomGig, error) {
	query := `Select * from custom_gigs where id = ? `
	var Gig model.CustomGig
	err := r.DB.Raw(query, ID).Scan(&Gig)
	return Gig, err.Error
}

func (r *GigRepo) CreateCustomOrder(order model.CustomOrder) error {
	err := r.DB.Create(&order)
	return err.Error
}

func (r *GigRepo) UpdateOrderStatus(order_id, status string) error {
	query := `UPDATE orders SET status = ? WHERE order_id = ?`
	err := r.DB.Exec(query, status, order_id).Error
	return err
}

func (r *GigRepo) UpdateOfferOrderStatus(order_id, status string) error {
	query := `UPDATE custom_orders SET status = ? WHERE order_id = ?`
	err := r.DB.Exec(query, status, order_id).Error
	return err
}

func (r *GigRepo) GetRequest(userid uint) ([]*proto.Request, []*proto.Request, error) {

	var result []model.Order
	status := "Freelance Approval Pending"
	query := `SELECT * FROM orders WHERE freelancer_id = ? AND status = ?`

	err := r.DB.Raw(query, userid, status).Scan(&result).Error
	if err != nil {
		fmt.Println("Error fetching orders:", err.Error())
		return nil, nil, err
	}

	var orders []*proto.Request
	for _, val := range result {
		orders = append(orders, &proto.Request{
			OrderId:  val.OrderID,
			GigId:    uint64(val.GigID),
			ClinetId: uint64(val.ClinetID),
			Amount:   int64(val.Amount),
		})
	}

	var result2 []model.CustomOrder
	query = `SELECT * FROM custom_orders WHERE freelancer_id = ? AND status = ?`
	err = r.DB.Raw(query, userid, status).Scan(&result2).Error
	if err != nil {
		fmt.Println("Error fetching custom orders:", err.Error())
		return nil, nil, err
	}

	var customOrders []*proto.Request
	for _, val := range result2 {
		customOrders = append(customOrders, &proto.Request{
			OrderId:  val.OrderID,
			GigId:    uint64(val.CustomGigID),
			ClinetId: uint64(val.ClinetID),
			Amount:   int64(val.Amount),
		})
	}

	return orders, customOrders, nil
}

func (r *GigRepo) AcceptOrder(orderID string) error {
	query := `UPDATE orders SET status = ? WHERE order_id = ? and status =?`
	err := r.DB.Exec(query, "Ongoing", orderID, "Freelance Approval Pending")
	if err.RowsAffected == 0 {
		return fmt.Errorf("no orders found with order_id %s and status 'Freelance Approval Pending'", orderID)
	}
	return err.Error
}

func (r *GigRepo) AcceptCustomOrder(orderID string) error {
	query := `UPDATE custom_orders SET status = ? WHERE order_id = ? and status =?`
	err := r.DB.Exec(query, "Ongoing", orderID, "Freelance Approval Pending")
	if err.RowsAffected == 0 {
		return fmt.Errorf("no orders found with order_id %s and status 'Freelance Approval Pending'", orderID)
	}
	return err.Error
}

func (r *GigRepo) RejectOrder(orderID string) error {
	query := `UPDATE orders SET status = ? WHERE order_id = ? and status =?`
	err := r.DB.Exec(query, "Freelancer Rejected", orderID, "Freelance Approval Pending")
	if err.RowsAffected == 0 {
		return fmt.Errorf("no orders found with order_id %s and status 'Freelance Approval Pending'", orderID)
	}
	return err.Error
}

func (r *GigRepo) RejectCustomOrder(orderID string) error {
	query := `UPDATE custom_orders SET status = ? WHERE order_id = ? and status =?`
	err := r.DB.Exec(query, "Freelancer Rejected", orderID, "Freelance Approval Pending")
	if err.RowsAffected == 0 {
		return fmt.Errorf("no orders found with order_id %s and status 'Freelance Approval Pending'", orderID)
	}
	return err.Error
}

func (r *GigRepo) AdminGetOrders() ([]*proto.Request, []*proto.Request, error) {
	var result []model.Order
	status := "Freelancer Rejected"
	query := `SELECT * FROM orders WHERE status = ?`

	err := r.DB.Raw(query, status).Scan(&result).Error
	if err != nil {
		fmt.Println("Error fetching orders:", err.Error())
		return nil, nil, err
	}

	var orders []*proto.Request
	for _, val := range result {
		orders = append(orders, &proto.Request{
			OrderId:      val.OrderID,
			GigId:        uint64(val.GigID),
			ClinetId:     uint64(val.ClinetID),
			FreelancerId: uint64(val.FreelancerID),
			Amount:       int64(val.Amount),
		})
	}

	var result2 []model.CustomOrder
	query = `SELECT * FROM custom_orders WHERE status = ?`
	err = r.DB.Raw(query, status).Scan(&result2).Error
	if err != nil {
		fmt.Println("Error fetching custom orders:", err.Error())
		return nil, nil, err
	}

	var customOrders []*proto.Request
	for _, val := range result2 {
		customOrders = append(customOrders, &proto.Request{
			OrderId:      val.OrderID,
			GigId:        uint64(val.CustomGigID),
			ClinetId:     uint64(val.ClinetID),
			FreelancerId: uint64(val.FreelancerID),
			Amount:       int64(val.Amount),
		})
	}
	return orders, customOrders, nil
}

func (r *GigRepo)GetOrderByID(order_id string)(model.Order,error){
	var order model.Order
	query:=`select * from orders where order_id = ?`
	err:=r.DB.Raw(query,order_id).Scan(&order).Error
	return order,err
}

func (r *GigRepo)GetCustomOrderByID(order_id string)(model.CustomOrder,error){
	var order model.CustomOrder
	query:=`select * from custom_orders where order_id = ?`
	err:=r.DB.Raw(query,order_id).Scan(&order).Error
	return order,err
}

// func (r *GigRepo)AdRefund(order_id string)error{

// }