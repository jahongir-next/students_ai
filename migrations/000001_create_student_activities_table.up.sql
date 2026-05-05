CREATE TABLE IF NOT EXISTS student_activities (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
        student_group VARCHAR(50) NOT NULL,
    academic_year VARCHAR(20) NOT NULL,
    semester VARCHAR(20) NOT NULL,
    payment_type VARCHAR(50),
    activity_type VARCHAR(20) NOT NULL CHECK (activity_type IN ('grade', 'attendance')),
    subject_name VARCHAR(255) NOT NULL,
    activity_value VARCHAR(50) NOT NULL,
    activity_date DATE DEFAULT CURRENT_DATE,

    lesson_pair INT DEFAULT 1,

    is_grade BOOLEAN DEFAULT FALSE,
    is_attendance BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE INDEX idx_activities_full_name ON student_activities(full_name);
CREATE INDEX idx_activities_group ON student_activities(student_group);
CREATE INDEX idx_activities_type_date ON student_activities(activity_type, activity_date);

CREATE UNIQUE INDEX idx_unique_student_grade
    ON student_activities (full_name, academic_year, semester, subject_name)
    WHERE (activity_type = 'grade');

CREATE UNIQUE INDEX idx_unique_student_attendance
    ON student_activities (full_name, subject_name, activity_date, lesson_pair)
    WHERE (activity_type = 'attendance');