package handlers

import (
	"errors"

	"polomnik/internal/application"
	"polomnik/internal/domain"
	"polomnik/internal/ports"
)

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func MapError(err error) *AppError {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return &AppError{Status: 404, Code: "NOT_FOUND", Message: "Ресурс не найден"}
	case errors.Is(err, application.ErrTourInactive):
		return &AppError{Status: 404, Code: "TOUR_INACTIVE", Message: "Тур недоступен"}
	case errors.Is(err, application.ErrTourExpired):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Тур уже завершился"}
	case errors.Is(err, application.ErrUnauthorized):
		return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Нужно войти в аккаунт"}
	case errors.Is(err, domain.ErrInsufficientSlots):
		return &AppError{Status: 409, Code: "INSUFFICIENT_SLOTS", Message: "Недостаточно свободных мест"}
	case errors.Is(err, domain.ErrInvalidPeopleCount):
		return &AppError{Status: 422, Code: "INVALID_PEOPLE_COUNT", Message: "Укажите количество участников больше нуля"}
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Нельзя изменить статус заявки"}
	case errors.Is(err, domain.ErrInvalidBookingStatus):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный статус заявки"}
	case errors.Is(err, domain.ErrInvalidCredentials):
		return &AppError{Status: 401, Code: "UNAUTHORIZED", Message: "Неверный логин или пароль"}
	case errors.Is(err, application.ErrPhoneVerificationUnavailable):
		return &AppError{Status: 503, Code: "PHONE_UNAVAILABLE", Message: "Пока что недоступно, используйте другой вариант."}
	case errors.Is(err, application.ErrPhoneVerificationRequired):
		return &AppError{Status: 422, Code: "PHONE_VERIFICATION_REQUIRED", Message: "Подтвердите телефон звонком с вашего номера"}
	case errors.Is(err, application.ErrPhoneVerificationNotConfirmed):
		return &AppError{Status: 422, Code: "PHONE_NOT_CONFIRMED", Message: "Звонок ещё не подтверждён или время истекло"}
	case errors.Is(err, application.ErrPhoneUserNotFound):
		return &AppError{Status: 404, Code: "PHONE_USER_NOT_FOUND", Message: "Аккаунт с этим телефоном не найден. Зарегистрируйтесь."}
	case errors.Is(err, ports.ErrPhoneChallengeFailed):
		return &AppError{Status: 502, Code: "PHONE_PROVIDER_ERROR", Message: "Не удалось начать проверку телефона. Попробуйте позже."}
	case errors.Is(err, domain.ErrDuplicateEmail):
		return &AppError{Status: 409, Code: "DUPLICATE_EMAIL", Message: "Этот email уже зарегистрирован"}
	case errors.Is(err, domain.ErrDuplicatePhone):
		return &AppError{Status: 409, Code: "DUPLICATE_PHONE", Message: "Этот телефон уже зарегистрирован"}
	case errors.Is(err, domain.ErrDuplicateSlug):
		return &AppError{Status: 409, Code: "DUPLICATE_SLUG", Message: "Статья или страница с таким адресом уже есть"}
	case errors.Is(err, domain.ErrDuplicatePath):
		return &AppError{Status: 409, Code: "DUPLICATE_PATH", Message: "Страница с таким путём уже есть"}
	default:
		return mapValidationError(err)
	}
}

func MapBookingError(err error) *AppError {
	if errors.Is(err, domain.ErrNotFound) {
		return &AppError{Status: 404, Code: "TOUR_NOT_FOUND", Message: "Тур не найден"}
	}
	return MapError(err)
}

func mapValidationError(err error) *AppError {
	switch {
	case errors.Is(err, domain.ErrInvalidID):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный идентификатор"}
	case errors.Is(err, domain.ErrInvalidSlug), errors.Is(err, domain.ErrInvalidPath):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный адрес страницы"}
	case errors.Is(err, domain.ErrInvalidTitle):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите название"}
	case errors.Is(err, domain.ErrInvalidPrice), errors.Is(err, domain.ErrInvalidTotalPrice):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректная стоимость"}
	case errors.Is(err, domain.ErrInvalidCurrency):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректная валюта"}
	case errors.Is(err, domain.ErrInvalidDateRange):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный период дат"}
	case errors.Is(err, domain.ErrInvalidSlots):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректное число мест"}
	case errors.Is(err, domain.ErrInvalidContactName), errors.Is(err, domain.ErrInvalidClientName):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите имя"}
	case errors.Is(err, domain.ErrInvalidPhone):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите корректный телефон"}
	case errors.Is(err, domain.ErrInvalidEmail):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите корректный email"}
	case errors.Is(err, domain.ErrInvalidPassword):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Пароль слишком короткий"}
	case errors.Is(err, domain.ErrInvalidRating):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Оценка должна быть от 1 до 5"}
	case errors.Is(err, domain.ErrInvalidReviewText):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Введите текст отзыва"}
	case errors.Is(err, domain.ErrInvalidSupportMessage):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Введите текст сообщения"}
	case errors.Is(err, domain.ErrInvalidTelegramUsername):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите корректный Telegram username (латиница, 5–32 символа)"}
	case errors.Is(err, domain.ErrInvalidNotificationChannel):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Неизвестный канал уведомлений"}
	case errors.Is(err, domain.ErrInvalidNotificationAddress):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный адрес получателя"}
	case errors.Is(err, domain.ErrInvalidNotificationEvent):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный тип события уведомления"}
	case errors.Is(err, domain.ErrInvalidAdminRoleName):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Имя роли: латиница, цифры, _ или -, 2–64 символа"}
	case errors.Is(err, domain.ErrInvalidPermission):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Неизвестное право доступа"}
	case errors.Is(err, domain.ErrDuplicateAdminRoleName):
		return &AppError{Status: 409, Code: "DUPLICATE_ROLE", Message: "Роль с таким именем уже есть"}
	case errors.Is(err, domain.ErrForbidden):
		return &AppError{Status: 403, Code: "FORBIDDEN", Message: "Недостаточно прав"}
	case errors.Is(err, domain.ErrInvalidSupportSender):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Некорректный отправитель сообщения"}
	case errors.Is(err, domain.ErrInvalidBlockType):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Неизвестный тип блока"}
	case errors.Is(err, domain.ErrInvalidExcerpt):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите анонс статьи"}
	case errors.Is(err, domain.ErrInvalidArticleBody):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите текст статьи"}
	case errors.Is(err, domain.ErrInvalidPublishedAt):
		return &AppError{Status: 422, Code: "VALIDATION_ERROR", Message: "Укажите дату публикации"}
	default:
		return nil
	}
}
