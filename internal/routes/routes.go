package routes

import (
	"LevelUp_Hub_Backend/internal/modules/auth"
	"LevelUp_Hub_Backend/internal/modules/booking"
	"LevelUp_Hub_Backend/internal/modules/connections"
	"LevelUp_Hub_Backend/internal/modules/courses"
	"LevelUp_Hub_Backend/internal/modules/favorites"
	"LevelUp_Hub_Backend/internal/modules/mentor_discovery"
	"LevelUp_Hub_Backend/internal/modules/message"
	"LevelUp_Hub_Backend/internal/modules/payment"
	"LevelUp_Hub_Backend/internal/modules/profile"
	"LevelUp_Hub_Backend/internal/modules/ratings"
	"LevelUp_Hub_Backend/internal/modules/slot"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUp(
	app *fiber.App,
	db *gorm.DB,
	rdb *redis.Client,
	jwtSecret string,
	rzpKey string,
	rzpSecret string,
) {

	// upgrade middleware globali
	app.Use("/api/messages/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	api := app.Group("/api")

	//module routes
	profile.RegisterRoutes(api, db, jwtSecret)
	auth.RegisterRoutes(api, db, rdb, jwtSecret)
	mentordiscovery.RegisterRoutes(api, db)
	courses.RegisterRoutes(api, db, jwtSecret)
	slot.RegisterRoutes(api, db, jwtSecret)
	message.RegisterRoutes(api, db, jwtSecret)
	ratings.RegisterRoutes(api, db, jwtSecret)
	favorites.RegisterRoutes(api,db,jwtSecret)
	connections.RegisterRoutes(api,db,jwtSecret)

	//////////// for booking and payment dependency wiring///////////////
	// ---------- REPOSITORIES ----------
	bookingRepo := booking.NewRepository(db)
	slotRepo := slot.NewRepository(db)
	mentorRepo := profile.NewMentorRepository(db)
	connectionRepo:=connections.NewRepository(db)
	paymentRepo := payment.NewRepository(db)

	// ---------- Connection Service ----------
	connection:=connections.NewService(connectionRepo,mentorRepo)

	// ---------- RAZORPAY ----------
	rzp := payment.NewRazorpayClient(rzpKey, rzpSecret)

	// ---------- SERVICES ----------
	bookingService := booking.NewService(
		bookingRepo,
		slotRepo,
		mentorRepo,
		connection,
		nil, // payment port set later
	)

	paymentService := payment.NewService(
		paymentRepo,
		bookingService, // BookingPort
		rzp,
		rzpKey,
	)

	// inject escrow dependency
	bookingService.SetPaymentPort(paymentService)

	// ---------- HANDLERS ----------
	bookingHandler := booking.NewHandler(bookingService)
	paymentHandler := payment.NewHandler(paymentService)

	// ---------- REGISTER ROUTES ----------
	booking.RegisterRoutes(api, jwtSecret, bookingHandler)
	payment.RegisterRoutes(api, jwtSecret, paymentHandler)
}
