
--To connect to PostgreSQL Press Ctrl+Shift+P in VScode → PostgreSQL: Add Connection

-- +goose UP
-- 1. Learners
CREATE TABLE learners (
  learner_id BIGSERIAL PRIMARY KEY
);

-- 2. Teachers
CREATE TABLE teachers (
  teacher_id  BIGSERIAL PRIMARY KEY,
  teacher_description VARCHAR(255)
);

-- 3. Admins
CREATE TABLE admins (
  admin_id BIGSERIAL PRIMARY KEY
);

-- 4. Users (supertype for Learner, Teacher, Admin)
CREATE TABLE users (
  user_id            BIGSERIAL     PRIMARY KEY,   --Big INT with auto increment
  session_id         VARCHAR(255),
  profile_picture    BYTEA,                       --Byte Array
  first_name         VARCHAR(30) NOT NULL,
  last_name          VARCHAR(30) NOT NULL,
  gender             VARCHAR(6),
  phone_number       VARCHAR(20),
  balance            NUMERIC(12,2) DEFAULT 0,     --MAX BALANCE IS 10 Billion
  created_at         TIMESTAMP DEFAULT NOW(),
  learner_id         BIGINT  UNIQUE,
  teacher_id         BIGINT  UNIQUE,
  admin_id           BIGINT  UNIQUE,

  --Have ID field set to null if that learner/teacher/admin is removed from table
  CONSTRAINT fk_user_learner
    FOREIGN KEY (learner_id) REFERENCES learners(learner_id)  ON DELETE SET NULL,
  CONSTRAINT fk_user_teacher
    FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id)  ON DELETE SET NULL,
  CONSTRAINT fk_user_admin
    FOREIGN KEY (admin_id)   REFERENCES admins(admin_id)      ON DELETE SET NULL
);

--//  How to insert Profile Picture  \\--
--INSERT INTO users (profile_picture)
--VALUES (
--  pg_read_binary_file('/path/to/avatar.png'),
--);



-- 5. Reports
CREATE TABLE reports (
  report_id            BIGSERIAL PRIMARY KEY,
  report_user_id       BIGINT      NOT NULL,  -- who files the report
  reported_user_id     BIGINT      NOT NULL,  -- who is being reported
  report_type          VARCHAR(20) NOT NULL,
  report_description   VARCHAR(255),
  report_picture       BYTEA,
  report_date          TIMESTAMP DEFAULT NOW(),

  CONSTRAINT fk_reporter
    FOREIGN KEY (report_user_id)   REFERENCES users(user_id)   ON DELETE CASCADE,
  CONSTRAINT fk_reported
    FOREIGN KEY (reported_user_id) REFERENCES users(user_id)   ON DELETE CASCADE
);

-- 6. Notifications
CREATE TABLE notifications (
  notification_id           BIGSERIAL   PRIMARY KEY,
  user_id                   BIGINT      NOT NULL,
  notification_type         VARCHAR(30) NOT NULL,
  notification_description  VARCHAR(255),
  notification_date         TIMESTAMP DEFAULT NOW(),
  read_flag                 BOOLEAN   DEFAULT FALSE,

  CONSTRAINT fk_notification_user
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- 7. Classes
CREATE TABLE classes (
  class_id             BIGSERIAL PRIMARY KEY,
  teacher_id           BIGINT    NOT NULL,
  learner_limit        INT       NOT NULL DEFAULT 50, --กันไว้ไม่เกิน 50 ละกัน
  class_description    VARCHAR(1000),
  enrollment_deadline  TIMESTAMP NOT NULL,

  CONSTRAINT fk_class_teacher
    FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id) ON DELETE CASCADE
);

-- 8. ClassSessions
CREATE TABLE class_sessions (
  class_session_id          BIGSERIAL PRIMARY KEY,
  class_id                  BIGINT,
  class_session_description VARCHAR(1000),
  enrollment_deadline       TIMESTAMP   NOT NULL,
  class_start               TIMESTAMP   NOT NULL,
  class_finish              TIMESTAMP   NOT NULL,
  class_status              VARCHAR(20)
);

-- 9. ClassCategories (many-to-many)
CREATE TABLE class_categories (
  class_id       BIGINT      NOT NULL,
  class_category VARCHAR(30) NOT NULL,
  PRIMARY KEY (class_id, class_category),

  CONSTRAINT fk_cc_class
    FOREIGN KEY (class_id) REFERENCES classes(class_id) ON DELETE CASCADE
);

-- 10. Reviews
CREATE TABLE reviews (
  learner_id BIGINT        NOT NULL,
  class_id   BIGINT        NOT NULL,
  rating     INT           NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment    VARCHAR(255),
  PRIMARY KEY (learner_id, class_id),

  CONSTRAINT fk_review_learner
    FOREIGN KEY (learner_id) REFERENCES learners(learner_id) ON DELETE CASCADE,
  CONSTRAINT fk_review_class
    FOREIGN KEY (class_id)   REFERENCES classes(class_id)   ON DELETE CASCADE
);

-- 11. Enrollments
CREATE TABLE enrollments (
  class_id              BIGINT NOT NULL,
  learner_id            BIGINT NOT NULL,
  enrollment_status     VARCHAR(20),
  PRIMARY KEY (class_id, learner_id),

  CONSTRAINT fk_enroll_class
    FOREIGN KEY (class_id)   REFERENCES classes(class_id)   ON DELETE CASCADE,
  CONSTRAINT fk_enroll_learner
    FOREIGN KEY (learner_id) REFERENCES learners(learner_id) ON DELETE CASCADE
);

-- 12. BanDetailsLearner
CREATE TABLE ban_details_learner (
  ban_learner_id   BIGSERIAL PRIMARY KEY,
  learner_id       BIGINT    NOT NULL,
  ban_start        TIMESTAMP DEFAULT NOW(),
  ban_end          TIMESTAMP NOT NULL,
  ban_description  VARCHAR(255),

  CONSTRAINT fk_ban_learner
    FOREIGN KEY (learner_id) REFERENCES learners(learner_id) ON DELETE CASCADE
);

-- 13. BanDetailsTeacher
CREATE TABLE ban_details_teacher (
  ban_teacher_id   BIGSERIAL PRIMARY KEY,
  teacher_id       BIGINT    NOT NULL,
  ban_start        TIMESTAMP DEFAULT NOW(),
  ban_end          TIMESTAMP NOT NULL,
  ban_description  VARCHAR(255),

  CONSTRAINT fk_ban_teacher
    FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id) ON DELETE CASCADE
);


-- +goose DOWN