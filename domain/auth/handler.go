package auth

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Handler - HTTP auth handler.
type Handler struct {
	useCase UseCase
}

// NewHandler - a factory function of auth handler.
func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h Handler) getRegister(c *fiber.Ctx) error {
	return c.Render("pages/auth/register", fiber.Map{})
}

func (h Handler) postRegister(c *fiber.Ctx) error {
	var cred Credentials
	if err := c.BodyParser(&cred); err != nil {
		return c.Render("pages/auth/register", fiber.Map{"Error": err.Error()})
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(cred.Password), 14)
	user := &User{
		Email: cred.Email,
		Auth: &Auth{
			Type:   "email",
			Secret: string(hashPassword),
		},
	}
	err := h.useCase.register(user)
	if err != nil {
		return c.Render("pages/auth/register", fiber.Map{"Error": err.Error()})
	}
	return c.Redirect("/auth/login")
}

func (h Handler) getLogin(c *fiber.Ctx) error {
	return c.Render("pages/auth/login", fiber.Map{})
}

func (h Handler) postLogin(c *fiber.Ctx) error {
	var cred Credentials
	err := c.BodyParser(&cred)
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	sess := c.Locals("session").(*session.Session)
	sessionToken, ok := sess.Get("session_token").(string)
	if !ok {
		return c.Render("pages/auth/login", fiber.Map{"Error": "Forbidden"})
	}
	err = h.useCase.login(cred, sessionToken)
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	token, err := h.jwtGenerate(&cred)
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
}

func (h Handler) jwtGenerate(cred *Credentials) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = cred.Email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fiber.ErrInternalServerError
	}
	return t, nil
}

func (h Handler) getLogout(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	err := sess.Destroy()
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	return c.Redirect("/auth/login")
}

func (h Handler) getProfile(c *fiber.Ctx) error {
	//sess, err := session.Get("session", c)
	//if err != nil {
	//	return c.SendStatus(http.StatusForbidden)
	//}
	//log.Printf("%v", sess.Values)
	//token, ok := sess.Values["session_token"].(string)
	//if !ok {
	//	b := new(bytes.Buffer)
	//	t.ViewErrorForbiddenPage(b)
	//	return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
	//}
	//_, err = h.useCase.Profile(token)
	//if err != nil {
	//	return c.JSON(http.StatusForbidden, err)
	//}
	//return c.SendStatus(http.StatusForbidden)

	return c.Render("pages/post/index", fiber.Map{
		"Title": "Show post list page",
	}, "layouts/main")
}
