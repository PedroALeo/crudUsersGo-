package models

import (
	"errors"
	"log"
	"strconv"

	"github.com/PedroALeo/CRUDUsers/database"
	"gopkg.in/validator.v2"
)

var ErrUserNotFound = errors.New("userNotFound")
var ErrBadRequestInvalidID = errors.New("invalidID")
var ErrBadRequest = errors.New("badRequest")

type User struct //struct that defines the User model
{
	ID       int
	Name     string `json:"name" validate:"nonzero"`
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
	Phone    string `json:"phone" validate:"nonzero"`
}

func UserValidation(user *User) error { //UserValidation recives a pointer to User and returns a error if the pointer content doesn't follow the model standarts
	if err := validator.Validate(user); err != nil {
		return err
	}
	return nil
}

func GetAllUsers() []User { //GetAllUsers gets all the user rows in the users table and append the users to a slice and returns this slice
	rows, _ := database.DB.Query("SELECT * from users")

	users := make([]User, 0)

	for rows.Next() {
		var id int64
		var name, email, password, phone string
		if err := rows.Scan(&id, &name, &email, &password, &phone); err != nil {
			log.Fatal(err)
		}
		users = append(users, User{int(id), name, email, password, phone})
	}

	return users
}

func GetUserById(id string) (User, error) { //GetUsersById recives the string that contains the id of the user who will be searched in the database, returns a empty user and an erro if occured in case of succes returns the user found
	idInt, _ := strconv.Atoi(id)
	rows := database.DB.QueryRow("SELECT id, nome, email, password, phone from users where id=$1", idInt)

	var user User

	if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone); err != nil {
		return User{}, ErrBadRequestInvalidID
	}
	if user.ID == 0 {

		return User{}, ErrUserNotFound
	}

	return user, nil
}

func NewUser(user *User) (User, error) { //NewUser recives a pointer to a user the will be used to creat a new entity in the database, in case of error returns a ErrBadRequest and an empty user in case of success returnes the complete user that was insertes in the database
	row := database.DB.QueryRow("INSERT INTO users(nome, email, password, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.Password, user.Password)
	if err := row.Scan(&user.ID); err != nil {
		return User{}, ErrBadRequest
	}
	userR := User{ID: user.ID, Name: user.Name, Email: user.Email, Password: user.Password, Phone: user.Phone}
	return userR, nil
}

func DeleteUserById(id string) (User, error) { //DeleteUserById recives a string that contains the id of the user who will be delete from the database, utilizes the GetById function to check if the user is realy in the database if it isn't returns a empty user and the ErrBadRequestInvalidId or ErrUserNotFound, in case of success returns the user that was deleted as a fail safe
	user, err := GetUserById(id)
	switch err {
	case ErrBadRequestInvalidID:
		return User{}, err
	case ErrUserNotFound:
		return User{}, err
	default:
		database.DB.QueryRow("DELETE FROM users WHERE id=$1", id)
		return user, nil
	}
}

func PutUserById(user *User) (User, error) { //PutUserById recives a pointer to the user that will be completly updated it returns ErrBadResquestInvalidId or ErrUserNotFound in case of errors and the updated user in case of success
	userR, err := GetUserById(strconv.Itoa(user.ID))
	switch err {
	case ErrBadRequestInvalidID:
		return User{}, err
	case ErrUserNotFound:
		return User{}, err
	default:
		row := database.DB.QueryRow("UPDATE users SET(nome, email, password, phone) = ($1, $2, $3, $4) WHERE id=$5 RETURNING *",
			user.Name, user.Email, user.Password, user.Phone, user.ID)
		if err := row.Scan(&userR.ID, &userR.Name, &userR.Email, &userR.Password, &userR.Phone); err != nil {
			return User{}, ErrUserNotFound
		}
		return userR, nil
	}
}

func PatchUserById(user *User) (User, error) { //PatchUserById recives a pointer to the user that will be updated in case of errors it can return ErrBadRequestInvalidId or ErrUserNotFound in case of success returns the updated user after the chances made in the database
	userResp, err := GetUserById(strconv.Itoa(user.ID))
	switch err {
	case ErrBadRequestInvalidID:
		return User{}, err
	case ErrUserNotFound:
		return User{}, err
	default:
		if user.Name != userResp.Name && user.Name != "" {
			userResp.Name = user.Name
		}
		if user.Email != userResp.Email && user.Email != "" {
			userResp.Email = user.Email
		}
		if user.Password != userResp.Password && user.Password != "" {
			userResp.Password = user.Password
		}
		if user.Phone != userResp.Phone && user.Phone != "" {
			userResp.Phone = user.Phone
		}

		row := database.DB.QueryRow("UPDATE users SET(nome, email, password, phone) = ($1, $2, $3, $4) WHERE id=$5 RETURNING *",
			userResp.Name, userResp.Email, userResp.Password, userResp.Phone, userResp.ID)
		if err := row.Scan(&userResp.ID, &userResp.Name, &userResp.Email, &userResp.Password, &userResp.Phone); err != nil {
			return User{}, ErrUserNotFound
		}
	}
	return userResp, nil
}
