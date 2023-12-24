package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm_tutorial/models"
	"gorm_tutorial/services"
	"gorm_tutorial/utils"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Create a DSN (data source name) to connect to the MySQL database.
	dsn := "root:hp123412@tcp(localhost:3306)/learn_go?charset=utf8&parseTime=True&loc=Asia%2FHo_Chi_Minh&allowNativePasswords=false"

	// Create a new GORM database connection.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}
	// to avoid unused db when comments all function
	fmt.Println("Connected to MySQL", db)

	// fmt.Println("------------------------Get All Users---------------------------")

	// allUsers(db)

	// fmt.Println("------------------------Get All Messages---------------------------")

	// allMsssages(db)

	// fmt.Println("------------------------Get User By Id---------------------------")

	// oneUser(db, 1)

	// fmt.Println("------------------------Get Users With Time Before CreateAt ---------------------------")

	// t1 := time.Date(2023, time.September, 20, 11, 23, 58, 0, time.UTC)
	// allUsersCreateBefore(db, t1)

	// fmt.Println("------------------------Get Users By Username And Password---------------------------")

	// login(db, "johndoe", "password")

	// fmt.Println("------------------------Insert User---------------------------")

	// register(db, "testuser888", "password", "testuser888@gmail.com")

	// fmt.Println("------------------------Update User By Id (Change Email And Password)---------------------------")

	// updateUser(db, 1, "newpassword", "john.doe99@gmail.com")

	// fmt.Println("------------------------Get all messages by sender user_id---------------------------")

	// allMsssagesBySender(db, 1)

	// fmt.Println("------------------------Get all messages by recipient user_id and sender_id---------------------------")

	// allMsssagesByUserIds(db, 1, 2)

	// fmt.Println("------------------------Get all messages by recipient user_id and sender_id but from view---------------------------")

	// allMsssagesWithUsernameWithView(db)

	// fmt.Println("------------------------Get all users with pagination---------------------------")

	// allUsersByPage(db, 2)

	// fmt.Println("------------------------Get all users with wildcard---------------------------")

	// searchUsername(db, "testuser")

	// fmt.Println("------------------------Count users by status---------------------------")

	// getUserStatusTotal(db)

	// fmt.Println("------------------------Call store procedure users wildcard---------------------------")

	// searchUsernameProcedure(db, "testuser")

	// fmt.Println("------------------------get users with 304 checking---------------------------")

	// t2 := time.Date(2023, time.October, 11, 12, 0, 0, 0, time.UTC)
	// t2 :=nil
	// streamAllUsers(db, t2)

	// fmt.Println("------------------------create and valid token---------------------------")

	token := createToken(db, 1)
	fmt.Println(token)
	verify := verifyToken(db, token)
	fmt.Println(verify)

}

func verifyToken(db *gorm.DB, token string) (valid bool) {
	valid, err := utils.VerifyToken(db, token)
	if err != nil {
		log.Fatalln("Cannot create token:", err)
	}
	return
}

func createToken(db *gorm.DB, id uint) (token string) {
	token, err := utils.CreateToken(db, 1)
	if err != nil {
		log.Fatalln("Cannot create token:", err)
	}
	return
}

func allUsers(db *gorm.DB) {
	users := services.GetAllUsers(db)
	// Serialize the users and messages to JSON strings.
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Cannot serialize users to JSON:", err)
	}
	// Print the JSON strings to the console.
	fmt.Println(string(usersJSON))
}

func allMsssages(db *gorm.DB) {
	messages := services.GetAllMessages(db)
	messagesJSON, err := json.Marshal(messages)
	if err != nil {
		log.Fatalln("Cannot serialize messages to JSON:", err)
	}

	fmt.Println(string(messagesJSON))
}

func oneUser(db *gorm.DB, id uint) {
	user, err := services.GetUserByID(db, id)
	if err != nil {
		log.Fatalln("Cannot get user by ID:", err)
	}
	fmt.Println(user.Status == models.UserStatusInactive)
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Cannot marshal user struct to JSON:", err)
	}
	fmt.Println(string(userJSON))
}

func allUsersCreateBefore(db *gorm.DB, time time.Time) {
	users, err := services.GetUsersByCreatedAtBefore(db, time)
	if err != nil {
		log.Fatalln("Cannot get users by createAt before:", err)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Cannot marshal user struct to JSON:", err)
	}
	fmt.Println(string(usersJSON))
}

func login(db *gorm.DB, username, password string) {
	users, err := services.GetUserByUsernameAndPassword(db, username, password)
	if err != nil {
		log.Fatalln("Cannot not login:", err)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Cannot marshal user struct to JSON:", err)
	}
	fmt.Println(string(usersJSON))
}

func register(db *gorm.DB, username, password, email string) {
	user, err := services.InsertUser(db, username, password, email)
	if err != nil {
		log.Fatalln("Cannot not insert user:", err)
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Cannot marshal user struct to JSON:", err)
	}
	fmt.Println(string(userJSON))
}

func updateUser(db *gorm.DB, id uint, password, email string) {
	user, err := services.UpdateUserWithChangedPasswordAndEmail(db, id, password, email)
	if err != nil {
		log.Fatalln("Cannot update user by ID:", err)
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatalln("Cannot marshal user struct to JSON:", err)
	}
	fmt.Println(string(userJSON))
}

func allUsersByPage(db *gorm.DB, page uint) {
	users, err := services.GetAllUsersWithPagination(db, page)
	if err != nil {
		log.Fatalln("Cannot get all users with pagination:", err)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Cannot serialize users to JSON:", err)
	}
	fmt.Println(string(usersJSON))
}

func searchUsername(db *gorm.DB, username string) {
	users, err := services.GetAllUsersByUsernameLike(db, username)
	if err != nil {
		log.Fatalln("Cannot get all users with wildcard:", err)
	}
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Cannot serialize users to JSON:", err)
	}
	fmt.Println(string(usersJSON))
}

func getUserStatusTotal(db *gorm.DB) {
	usersStatus, err := services.GetCountUsersGroupedByStatus(db)
	if err != nil {
		log.Fatalln("Cannot count users by status:", err)
	}
	usersJSON, err := json.Marshal(usersStatus)
	if err != nil {
		log.Fatalln("Cannot serialize users to JSON:", err)
	}
	fmt.Println(string(usersJSON))
}

func searchUsernameProcedure(db *gorm.DB, username string) {
	usersStatus, err := services.GetAllUsersWithUsernameLikeUsingProcedure(db, username)
	if err != nil {
		log.Fatalln("Cannot get users from store procedure:", err)
	}
	usersJSON, err := json.Marshal(usersStatus)
	if err != nil {
		log.Fatalln("Cannot serialize users to JSON:", err)
	}
	fmt.Println(string(usersJSON))
}

func streamAllUsers(db *gorm.DB, time time.Time) {
	users, lastest, err := services.GetAllUsersWithTime(db, &time)
	// users, lastest, err := services.GetAllUsersWithTime(db, nil)
	if err != nil {
		if errors.Is(err, utils.ErrNoModified) {
			// return empty interface (dynamic)
			res := models.Result{Time: lastest, Result: []interface{}{}, Error: "Not Modified"}
			resJSON, err := json.Marshal(res)
			if err != nil {
				log.Fatalln("Cannot serialize result to JSON:", err)
			}
			fmt.Println(string(resJSON))
		} else {
			res := models.Result{Time: lastest, Result: []interface{}{}, Error: err.Error()}
			resJSON, err := json.Marshal(res)
			if err != nil {
				log.Fatalln("Cannot serialize result to JSON:", err)
			}
			fmt.Println(string(resJSON))
		}
	} else {
		usersInterface := make([]interface{}, len(users))
		for i, v := range users {
			usersInterface[i] = v
		}
		// error in json will empty ""
		res := models.Result{Time: lastest, Result: usersInterface}
		resJSON, err := json.Marshal(res)
		if err != nil {
			log.Fatalln("Cannot serialize result to JSON:", err)
		}
		fmt.Println(string(resJSON))
	}
}

func allMsssagesBySender(db *gorm.DB, id uint) {
	messagesSender, err := services.GetMessagesBySender(db, id)
	if err != nil {
		log.Fatalln("Cannot get all messages by user_id:", err)
	}
	messagesSenderJSON, err := json.Marshal(messagesSender)
	if err != nil {
		log.Fatalln("Cannot serialize messages to JSON:", err)
	}
	fmt.Println(string(messagesSenderJSON))
}

func allMsssagesByUserIds(db *gorm.DB, id1, id2 uint) {
	messages, err := services.GetMessagesWithSenderAndRecipientNames(db, id1, id2)
	if err != nil {
		log.Fatalln("Cannot get messages with 2 user id:", err)
	}
	messagesJson, err := json.Marshal(messages)
	if err != nil {
		log.Fatalln("Cannot serialize messages to JSON:", err)
	}
	fmt.Println(string(messagesJson))
}

func allMsssagesWithUsernameWithView(db *gorm.DB) {
	messages, err := services.GetMessagesWithSenderAndRecipientNamesFromView(db)
	if err != nil {
		log.Fatalln("Cannot get messages with sender and recipient names:", err)
	}
	messagesJson, err := json.Marshal(messages)
	if err != nil {
		log.Fatalln("Cannot serialize messages to JSON:", err)
	}
	fmt.Println(string(messagesJson))
}
