package models

import (
	"fmt"
	"log"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&Admin{},
		&Learner{},
		&Teacher{},
		&BanDetailsLearner{},
		&BanDetailsTeacher{},
		&Class{},
		&ClassCategory{},
		&ClassSession{},
		&Enrollment{},
		&Notification{},
		&Report{},
		&Review{},
		&Transaction{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("Database migrated successfully")

	if config.STATUS() == "development" {
		err := Seed(db)
		if err != nil {
			log.Fatalf("seed failed: %v", err)
		}
		log.Println("Database seeded successfully")
	} else {
		log.Println("Skip seeded")
	}
}

/* -------------------- Helper for seed the database ,it will do nothing if entry already exist -------------------- */

func seedHelper[T any](tx *gorm.DB, items []T, conflictCols ...string) error {
	if len(items) == 0 {
		return nil
	}
	on := clause.OnConflict{DoNothing: true}
	if len(conflictCols) > 0 {
		on.Columns = make([]clause.Column, 0, len(conflictCols))
		for _, c := range conflictCols {
			on.Columns = append(on.Columns, clause.Column{Name: c})
		}
	}
	return tx.Clauses(on).Create(&items).Error
}

func idMap[K comparable, M any](tx *gorm.DB, model any, col string, keys []K, into map[K]M) error {
	for _, k := range keys {
		var row M
		q := tx.Model(model)
		if err := q.Where(fmt.Sprintf("%s = ?", col), k).First(&row).Error; err != nil {
			return err
		}
		into[k] = row
	}
	return nil
}

func fetchAllRows[T any](tx *gorm.DB) ([]T, error) {
	var rows []T
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

/* -------------------- Seed -------------------- */

func Seed(db *gorm.DB) error {
	now := time.Now()
	deadline := now.Add(7 * 24 * time.Hour)
	start := now.Add(10 * 24 * time.Hour)
	finish := start.Add(2 * time.Hour)

	return db.Transaction(func(tx *gorm.DB) error {
		/* ---------- Users ---------- */
		userKeys := []string{
			"b6600000000",
			"b6600000001",
			"b6600000002",
			"b6600000003",
			"b6600000004",
			"b6600000005",
			"b65000000001",
			"b6700000002",
		}

		users := []User{
			{StudentID: userKeys[0], FirstName: "Alice", LastName: "Admin", Gender: "Female", PhoneNumber: "+66000000000", Balance: 100000, BanCount: 0},
			{StudentID: userKeys[1], FirstName: "Bob", LastName: "Learner", Gender: "Male", PhoneNumber: "+66000000001", Balance: 100, BanCount: 0},
			{StudentID: userKeys[2], FirstName: "Carol", LastName: "Teacher", Gender: "Female", PhoneNumber: "+66000000002", Balance: 5000, BanCount: 0},
			{StudentID: userKeys[3], FirstName: "Dave", LastName: "Learner", Gender: "Male", PhoneNumber: "+66000000003", Balance: 600, BanCount: 0},
			{StudentID: userKeys[4], FirstName: "Eve", LastName: "Teacher", Gender: "Female", PhoneNumber: "+66000000004", Balance: 800, BanCount: 0},
			{StudentID: userKeys[5], FirstName: "Frank", LastName: "Admin", Gender: "Male", PhoneNumber: "+66000000005", Balance: 100000, BanCount: 0},
			{StudentID: userKeys[6], FirstName: "Ban", LastName: "Teacher", Gender: "Male", PhoneNumber: "+65000000001", Balance: 1000, BanCount: 1},
			{StudentID: userKeys[7], FirstName: "Ban", LastName: "Learner", Gender: "Male", PhoneNumber: "+67000000002", Balance: 10, BanCount: 1},
		}
		if err := seedHelper(tx, users, "student_id"); err != nil {
			return err
		}

		userBySID := map[string]User{}
		if err := idMap(tx, &User{}, "student_id", userKeys, userBySID); err != nil {
			return err
		}

		/* ---------- Admins / Learners / Teachers ---------- */
		adminRows := []Admin{
			{UserID: userBySID[userKeys[0]].ID},
			{UserID: userBySID[userKeys[5]].ID},
		}
		if err := seedHelper(tx, adminRows, "user_id"); err != nil {
			return err
		}

		learnerRows := []Learner{
			{UserID: userBySID[userKeys[1]].ID, FlagCount: 0},
			{UserID: userBySID[userKeys[3]].ID, FlagCount: 1},
			{UserID: userBySID[userKeys[7]].ID, FlagCount: 3},
		}
		if err := seedHelper(tx, learnerRows, "user_id"); err != nil {
			return err
		}

		teacherRows := []Teacher{
			{UserID: userBySID[userKeys[2]].ID, Email: "carol.teacher@example.com", Description: "Experienced teacher", FlagCount: 0},
			{UserID: userBySID[userKeys[4]].ID, Email: "eve.teacher@example.com", Description: "Experienced teacher", FlagCount: 1},
			{UserID: userBySID[userKeys[6]].ID, Email: "Ban.teacher@example.com", Description: "Experienced teacher", FlagCount: 3},
		}
		if err := seedHelper(tx, teacherRows, "user_id"); err != nil {
			return err
		}

		allTeachers, err := fetchAllRows[Teacher](tx)
		if err != nil {
			return err
		}
		teacherByEmail := map[string]Teacher{}
		for _, t := range allTeachers {
			teacherByEmail[t.Email] = t
		}

		/* ---------- Categories ---------- */
		categoryRows := []ClassCategory{
			{ClassCategory: "Mathematics"},
			{ClassCategory: "Science"},
			{ClassCategory: "Languages"},
			{ClassCategory: "Programming"},
			{ClassCategory: "Art"},
			{ClassCategory: "History"},
		}
		if err := seedHelper(tx, categoryRows); err != nil {
			return err
		}

		/* ---------- Classes  ---------- */
		classRows := []Class{
			{TeacherID: teacherByEmail["carol.teacher@example.com"].ID, ClassName: "Algebra I", ClassDescription: "Foundations of algebra", Rating: 4},
			{TeacherID: teacherByEmail["carol.teacher@example.com"].ID, ClassName: "Intro to Go", ClassDescription: "Go language for beginners", Rating: 2},
			{TeacherID: teacherByEmail["eve.teacher@example.com"].ID, ClassName: "English A1", ClassDescription: "Basic English course", Rating: 3},
			{TeacherID: teacherByEmail["eve.teacher@example.com"].ID, ClassName: "Watercolor Basics", ClassDescription: "Painting techniques", Rating: 4},
			{TeacherID: teacherByEmail["carol.teacher@example.com"].ID, ClassName: "Physics 101", ClassDescription: "Mechanics and waves", Rating: 3.5},
			{TeacherID: teacherByEmail["carol.teacher@example.com"].ID, ClassName: "History of Tech", ClassDescription: "From abacus to AI", Rating: 4},
		}
		if err := seedHelper(tx, classRows); err != nil {
			return err
		}

		classNames := []string{
			"Algebra I", "Intro to Go", "English A1", "Watercolor Basics", "Physics 101", "History of Tech",
		}
		classByName := map[string]Class{}
		if err := idMap(tx, &Class{}, "class_name", classNames, classByName); err != nil {
			return err
		}

		/* ---------- Sessions (by ClassID) ---------- */
		sessionRows := []ClassSession{
			{ClassID: classByName["Algebra I"].ID, Description: "Weekday evening", Price: 1500, LearnerLimit: 30, EnrollmentDeadline: deadline, ClassStart: start, ClassFinish: finish, ClassStatus: "open"},
			{ClassID: classByName["Intro to Go"].ID, Description: "Weekend bootcamp", Price: 2500, LearnerLimit: 25, EnrollmentDeadline: deadline, ClassStart: start.Add(24 * time.Hour), ClassFinish: finish.Add(24 * time.Hour), ClassStatus: "open"},
			{ClassID: classByName["English A1"].ID, Description: "Morning session", Price: 1200, LearnerLimit: 40, EnrollmentDeadline: deadline, ClassStart: start.Add(48 * time.Hour), ClassFinish: finish.Add(48 * time.Hour), ClassStatus: "open"},
			{ClassID: classByName["Watercolor Basics"].ID, Description: "Afternoon studio", Price: 1800, LearnerLimit: 20, EnrollmentDeadline: deadline, ClassStart: start.Add(72 * time.Hour), ClassFinish: finish.Add(72 * time.Hour), ClassStatus: "open"},
			{ClassID: classByName["Physics 101"].ID, Description: "Evening lab", Price: 1600, LearnerLimit: 35, EnrollmentDeadline: deadline, ClassStart: start.Add(96 * time.Hour), ClassFinish: finish.Add(96 * time.Hour), ClassStatus: "open"},
			{ClassID: classByName["Algebra I"].ID, Description: "Extra section", Price: 1500, LearnerLimit: 25, EnrollmentDeadline: deadline, ClassStart: start.Add(120 * time.Hour), ClassFinish: finish.Add(120 * time.Hour), ClassStatus: "open"},
		}

		if err := seedHelper(tx, sessionRows); err != nil {
			return err
		}

		sessionDescs := []string{
			"Weekday evening", "Weekend bootcamp", "Morning session", "Afternoon studio", "Evening lab", "Extra section",
		}
		sessionByDesc := map[string]ClassSession{}
		if err := idMap(tx, &ClassSession{}, "description", sessionDescs, sessionByDesc); err != nil {
			return err
		}

		/* ---------- Learner mapping by student_id ---------- */
		allLearners, err := fetchAllRows[Learner](tx)
		if err != nil {
			return err
		}
		learnerBySID := map[string]Learner{}
		for _, ln := range allLearners {
			for sid, u := range userBySID {
				if ln.UserID == u.ID {
					learnerBySID[sid] = ln
				}
			}
		}

		/* ---------- Enrollments ---------- */
		enrollmentRows := []Enrollment{
			{LearnerID: learnerBySID[userKeys[1]].ID, ClassSessionID: sessionByDesc["Weekday evening"].ID, EnrollmentStatus: "active"},
			{LearnerID: learnerBySID[userKeys[3]].ID, ClassSessionID: sessionByDesc["Weekday evening"].ID, EnrollmentStatus: "active"},
			{LearnerID: learnerBySID[userKeys[1]].ID, ClassSessionID: sessionByDesc["Weekend bootcamp"].ID, EnrollmentStatus: "active"},
			{LearnerID: learnerBySID[userKeys[3]].ID, ClassSessionID: sessionByDesc["Morning session"].ID, EnrollmentStatus: "active"},
			{LearnerID: learnerBySID[userKeys[1]].ID, ClassSessionID: sessionByDesc["Morning session"].ID, EnrollmentStatus: "active"},
		}
		if err := seedHelper(tx, enrollmentRows); err != nil {
			return err
		}

		/* ---------- Notifications ---------- */
		notificationRows := []Notification{
			{UserID: userBySID[userKeys[1]].ID, NotificationType: "system", NotificationDescription: "Welcome!", NotificationDate: now, ReadFlag: false},
			{UserID: userBySID[userKeys[1]].ID, NotificationType: "enrollment", NotificationDescription: "You enrolled successfully.", NotificationDate: now, ReadFlag: false},
			{UserID: userBySID[userKeys[2]].ID, NotificationType: "class", NotificationDescription: "Your class got a new review.", NotificationDate: now, ReadFlag: false},
			{UserID: userBySID[userKeys[1]].ID, NotificationType: "system", NotificationDescription: "Password changed.", NotificationDate: now, ReadFlag: true},
			{UserID: userBySID[userKeys[0]].ID, NotificationType: "report", NotificationDescription: "A report needs your attention.", NotificationDate: now, ReadFlag: false},
			{UserID: userBySID[userKeys[5]].ID, NotificationType: "system", NotificationDescription: "Balance updated.", NotificationDate: now, ReadFlag: false},
		}
		if err := seedHelper(tx, notificationRows); err != nil {
			return err
		}

		/* ---------- Reports ---------- */
		reportRows := []Report{
			{
				ReportUserID:      userBySID[userKeys[1]].ID,
				ReportedUserID:    userBySID[userKeys[4]].ID,
				ClassSessionID:    sessionByDesc["Weekday evening"].ID,
				ReportType:        "behavior",
				ReportReason:      "spam",
				ReportDescription: "Spam messages",
				ReportDate:        now,
				ReportStatus:      "pending",
			},
			{
				ReportUserID:      userBySID[userKeys[3]].ID,
				ReportedUserID:    userBySID[userKeys[6]].ID,
				ClassSessionID:    sessionByDesc["Weekday evening"].ID,
				ReportType:        "content",
				ReportReason:      "inappropriate",
				ReportDescription: "Inappropriate content",
				ReportDate:        now,
				ReportStatus:      "pending",
			},
		}
		if err := seedHelper(tx, reportRows); err != nil {
			return err
		}

		/* ---------- Reviews ---------- */
		reviewRows := []Review{
			{LearnerID: learnerBySID[userKeys[1]].ID, ClassID: classByName["Algebra I"].ID, Rating: 5, Comment: "Great class!"},
			{LearnerID: learnerBySID[userKeys[3]].ID, ClassID: classByName["Algebra I"].ID, Rating: 4, Comment: "Very helpful."},
			{LearnerID: learnerBySID[userKeys[3]].ID, ClassID: classByName["Intro to Go"].ID, Rating: 5, Comment: "Excellent!"},
			{LearnerID: learnerBySID[userKeys[3]].ID, ClassID: classByName["Physics 101"].ID, Rating: 4, Comment: "Good session."},
			{LearnerID: learnerBySID[userKeys[1]].ID, ClassID: classByName["Watercolor Basics"].ID, Rating: 5, Comment: "Loved it!"},
		}
		if err := seedHelper(tx, reviewRows); err != nil {
			return err
		}

		/* ---------- Transactions ---------- */
		uid1 := userBySID[userKeys[1]].ID
		uid2 := userBySID[userKeys[2]].ID
		uid3 := userBySID[userKeys[3]].ID
		uid4 := userBySID[userKeys[4]].ID
		uid5 := userBySID[userKeys[5]].ID
		failCode := "insufficient_funds"
		failMsg := "Insufficient funds"
		transactionRows := []Transaction{
			{UserID: &uid1, ChargeID: "ch_0000000000000001", AmountSatang: 150000, Currency: "THB", Channel: "card", Status: "paid"},
			{UserID: &uid2, ChargeID: "ch_0000000000000002", AmountSatang: 250000, Currency: "THB", Channel: "card", Status: "paid"},
			{UserID: &uid3, ChargeID: "ch_0000000000000003", AmountSatang: 120000, Currency: "THB", Channel: "bank_transfer", Status: "pending"},
			{UserID: &uid4, ChargeID: "ch_0000000000000004", AmountSatang: 180000, Currency: "THB", Channel: "card", Status: "paid"},
			{UserID: &uid5, ChargeID: "ch_0000000000000005", AmountSatang: 160000, Currency: "THB", Channel: "card", Status: "failed", FailureCode: &failCode, FailureMessage: &failMsg},
		}
		if err := seedHelper(tx, transactionRows, "charge_id"); err != nil {
			return err
		}

		/* ---------- Bans ---------- */
		banLearnerRows := []BanDetailsLearner{
			{LearnerID: learnerBySID[userKeys[7]].ID, BanStart: now, BanEnd: now.Add(24 * time.Hour), BanDescription: "Test ban L01"},
		}
		if err := seedHelper(tx, banLearnerRows); err != nil {
			return err
		}

		teacherByUserID := map[uint]Teacher{}
		for _, t := range allTeachers {
			teacherByUserID[t.UserID] = t
		}
		banTeacherRows := []BanDetailsTeacher{
			{TeacherID: teacherByUserID[userBySID[userKeys[6]].ID].ID, BanStart: now, BanEnd: now.Add(24 * time.Hour), BanDescription: "Test ban T01"},
		}
		if err := seedHelper(tx, banTeacherRows); err != nil {
			return err
		}

		return nil
	})
}
