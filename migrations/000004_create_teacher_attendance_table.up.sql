CREATE TABLE IF NOT EXISTS teacher_attendance_logs (
    id SERIAL PRIMARY KEY,
    teacher_name VARCHAR(255) NOT NULL,
    group_name VARCHAR(50) NOT NULL,
    subject_name VARCHAR(255) NOT NULL,
    lesson_type VARCHAR(50) NOT NULL,
    academic_year VARCHAR(20) NOT NULL,
    semester VARCHAR(20) NOT NULL,
    scheduled_start_time TIME NOT NULL,
    actual_arrival_time TIMESTAMP NOT NULL,
    arrival_date DATE NOT NULL,
    status VARCHAR(20),
    delay_minutes INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO teacher_attendance_logs (teacher_name, group_name, subject_name, lesson_type, academic_year, semester, scheduled_start_time, actual_arrival_time, arrival_date, status, delay_minutes)
VALUES
-- =========================================================
-- 1-SEMESTR (KUZGI SEMESTR: 2024 Sentyabr - Dekabr) - 50 ta
-- =========================================================
-- Sunnatov T.R (Falsafa)
('Sunnatov T.R', 'MAT 24/1', 'Falsafa', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-02 08:25:00', '2024-09-02', 'present', 0),
('Sunnatov T.R', 'MAT 24/2', 'Falsafa', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-16 08:42:00', '2024-09-16', 'late', 12),
('Sunnatov T.R', 'MAT 24/3', 'Falsafa', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-10-07 08:30:00', '2024-10-07', 'present', 0),
('Sunnatov T.R', 'MAT 24/4', 'Falsafa', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-10-21 10:15:00', '2024-10-21', 'late', 15),
('Sunnatov T.R', 'MAT 24/1', 'Falsafa', 'Amaliy', '2024-2025', '1-semestr', '11:30:00', '2024-11-11 11:28:00', '2024-11-11', 'present', 0),
('Sunnatov T.R', 'MAT 24/2', 'Falsafa', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-12-02 10:25:00', '2024-12-02', 'late', 25),

-- Ismanova M.A. (Bolalarning ijtimoiy moslashuvi)
('Ismanova M.A.', 'MAT 24/1', 'Bolalarning ijtimoiy moslashuvi', 'Ma''ruza', '2024-2025', '1-semestr', '10:00:00', '2024-09-05 09:58:00', '2024-09-05', 'present', 0),
('Ismanova M.A.', 'MAT 24/2', 'Bolalarning ijtimoiy moslashuvi', 'Ma''ruza', '2024-2025', '1-semestr', '11:30:00', '2024-09-19 11:40:00', '2024-09-19', 'late', 10),
('Ismanova M.A.', 'MAT 24/6', 'Bolalarning ijtimoiy moslashuvi', 'Amaliy', '2024-2025', '1-semestr', '08:30:00', '2024-10-03 08:29:00', '2024-10-03', 'present', 0),
('Ismanova M.A.', 'MAT 24/7', 'Bolalarning ijtimoiy moslashuvi', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-10-24 10:20:00', '2024-10-24', 'late', 20),
('Ismanova M.A.', 'MAT 24/1', 'Bolalarning ijtimoiy moslashuvi', 'Amaliy', '2024-2025', '1-semestr', '11:30:00', '2024-11-07 11:30:00', '2024-11-07', 'present', 0),
('Ismanova M.A.', 'MAT 24/2', 'Bolalarning ijtimoiy moslashuvi', 'Amaliy', '2024-2025', '1-semestr', '08:30:00', '2024-12-12 08:45:00', '2024-12-12', 'late', 15),

-- Xurvaliyeva T.L. (Bolalar psixologiyasi)
('Xurvaliyeva T.L.', 'MAT 24/1', 'Bolalar psixologiyasi', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-10 08:25:00', '2024-09-10', 'present', 0),
('Xurvaliyeva T.L.', 'MAT 24/2', 'Bolalar psixologiyasi', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-24 08:40:00', '2024-09-24', 'late', 10),
('Xurvaliyeva T.L.', 'MAT 24/3', 'Bolalar psixologiyasi', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-10-08 09:55:00', '2024-10-08', 'present', 0),
('Xurvaliyeva T.L.', 'MAT 24/4', 'Bolalar psixologiyasi', 'Amaliy', '2024-2025', '1-semestr', '11:30:00', '2024-10-22 11:45:00', '2024-10-22', 'late', 15),
('Xurvaliyeva T.L.', 'MAT 24/1', 'Bolalar psixologiyasi', 'Amaliy', '2024-2025', '1-semestr', '08:30:00', '2024-11-12 08:30:00', '2024-11-12', 'present', 0),
('Xurvaliyeva T.L.', 'MAT 24/2', 'Bolalar psixologiyasi', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-12-03 10:05:00', '2024-12-03', 'late', 5),

-- Maripova N.X (Bolalarni sahnalashtirish)
('Maripova N.X', 'MAT 24/1', 'Bolalarni sahnalashtirish', 'Ma''ruza', '2024-2025', '1-semestr', '11:30:00', '2024-09-13 11:28:00', '2024-09-13', 'present', 0),
('Maripova N.X', 'MAT 24/2', 'Bolalarni sahnalashtirish', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-27 08:42:00', '2024-09-27', 'late', 12),
('Maripova N.X', 'MAT 24/3', 'Bolalarni sahnalashtirish', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-10-11 10:00:00', '2024-10-11', 'present', 0),
('Maripova N.X', 'MAT 24/4', 'Bolalarni sahnalashtirish', 'Amaliy', '2024-2025', '1-semestr', '11:30:00', '2024-10-25 11:50:00', '2024-10-25', 'late', 20),
('Maripova N.X', 'MAT 24/5', 'Bolalarni sahnalashtirish', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-11-15 08:29:00', '2024-11-15', 'present', 0),

-- Mutalova D.A. (Ilk va maktabgacha pedagogika)
('Mutalova D.A.', 'MAT 24/3', 'Ilk va maktabgacha pedagogika', 'Ma''ruza', '2024-2025', '1-semestr', '10:00:00', '2024-09-18 09:55:00', '2024-09-18', 'present', 0),
('Mutalova D.A.', 'MAT 24/4', 'Ilk va maktabgacha pedagogika', 'Ma''ruza', '2024-2025', '1-semestr', '11:30:00', '2024-10-02 11:38:00', '2024-10-02', 'late', 8),
('Mutalova D.A.', 'MAT 24/5', 'Ilk va maktabgacha pedagogika', 'Amaliy', '2024-2025', '1-semestr', '08:30:00', '2024-10-16 08:30:00', '2024-10-16', 'present', 0),
('Mutalova D.A.', 'MAT 24/3', 'Ilk va maktabgacha pedagogika', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-11-20 10:25:00', '2024-11-20', 'late', 25),
('Mutalova D.A.', 'MAT 24/4', 'Ilk va maktabgacha pedagogika', 'Ma''ruza', '2024-2025', '1-semestr', '11:30:00', '2024-12-11 11:25:00', '2024-12-11', 'present', 0),

-- Urmanova T.I. (Rivojlantiruvchi markazlar)
('Urmanova T.I.', 'MAT 24/1', 'Rivojlantiruvchi markazlar', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-04 08:30:00', '2024-09-04', 'present', 0),
('Urmanova T.I.', 'MAT 24/2', 'Rivojlantiruvchi markazlar', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-10-09 10:15:00', '2024-10-09', 'late', 15),
('Urmanova T.I.', 'MAT 24/7', 'Rivojlantiruvchi markazlar', 'Ma''ruza', '2024-2025', '1-semestr', '11:30:00', '2024-11-06 11:28:00', '2024-11-06', 'present', 0),

-- Tursunbayeva M.D (Bolalarning ijtimoiy moslashuvi)
('Tursunbayeva M.D', 'MAT 24/3', 'Bolalarning ijtimoiy moslashuvi', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-09-11 09:58:00', '2024-09-11', 'present', 0),
('Tursunbayeva M.D', 'MAT 24/4', 'Bolalarning ijtimoiy moslashuvi', 'Amaliy', '2024-2025', '1-semestr', '11:30:00', '2024-10-23 11:45:00', '2024-10-23', 'late', 15),

-- Madalimov T.A. (Falsafa)
('Madalimov T.A.', 'MAT 24/5', 'Falsafa', 'Ma''ruza', '2024-2025', '1-semestr', '08:30:00', '2024-09-26 08:25:00', '2024-09-26', 'present', 0),
('Madalimov T.A.', 'MAT 24/5', 'Falsafa', 'Amaliy', '2024-2025', '1-semestr', '10:00:00', '2024-11-14 10:10:00', '2024-11-14', 'late', 10),


-- =========================================================
-- 2-SEMESTR (BAHORGI SEMESTR: 2025 Fevral - May) - 50 ta
-- =========================================================
-- Nafasov A.K. (O'zbekistonning eng yangi tarixi)
('Nafasov A.K.', 'MAT 24/1', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '08:30:00', '2025-02-03 08:28:00', '2025-02-03', 'present', 0),
('Nafasov A.K.', 'MAT 24/2', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '10:00:00', '2025-02-17 10:12:00', '2025-02-17', 'late', 12),
('Nafasov A.K.', 'MAT 24/1', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-03-03 11:30:00', '2025-03-03', 'present', 0),
('Nafasov A.K.', 'MAT 24/2', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-03-17 08:50:00', '2025-03-17', 'late', 20),
('Nafasov A.K.', 'MAT 24/1', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '10:00:00', '2025-04-07 09:55:00', '2025-04-07', 'present', 0),
('Nafasov A.K.', 'MAT 24/2', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-05-05 11:45:00', '2025-05-05', 'late', 15),

-- Xudaynazarov A (Mediasavodxonlik)
('Xudaynazarov A', 'MAT 24/1', 'Mediasavodxonlik', 'Ma''ruza', '2024-2025', '2-semestr', '11:30:00', '2025-02-04 11:25:00', '2025-02-04', 'present', 0),
('Xudaynazarov A', 'MAT 24/2', 'Mediasavodxonlik', 'Ma''ruza', '2024-2025', '2-semestr', '08:30:00', '2025-02-18 08:45:00', '2025-02-18', 'late', 15),
('Xudaynazarov A', 'MAT 24/3', 'Mediasavodxonlik', 'Ma''ruza', '2024-2025', '2-semestr', '10:00:00', '2025-03-04 10:00:00', '2025-03-04', 'present', 0),
('Xudaynazarov A', 'MAT 24/4', 'Mediasavodxonlik', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-03-18 11:55:00', '2025-03-18', 'late', 25),
('Xudaynazarov A', 'MAT 24/1', 'Mediasavodxonlik', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-04-15 08:30:00', '2025-04-15', 'present', 0),
('Xudaynazarov A', 'MAT 24/2', 'Mediasavodxonlik', 'Amaliy', '2024-2025', '2-semestr', '10:00:00', '2025-05-13 10:10:00', '2025-05-13', 'late', 10),

-- Sa'dullayeva M.X. (O'zbekistonning eng yangi tarixi)
('Sa''dullayeva M.X.', 'MAT 24/3', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '10:00:00', '2025-02-05 09:58:00', '2025-02-05', 'present', 0),
('Sa''dullayeva M.X.', 'MAT 24/4', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '11:30:00', '2025-02-19 11:42:00', '2025-02-19', 'late', 12),
('Sa''dullayeva M.X.', 'MAT 24/6', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '08:30:00', '2025-03-05 08:30:00', '2025-03-05', 'present', 0),
('Sa''dullayeva M.X.', 'MAT 24/3', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '10:00:00', '2025-03-19 10:20:00', '2025-03-19', 'late', 20),
('Sa''dullayeva M.X.', 'MAT 24/4', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-04-23 11:28:00', '2025-04-23', 'present', 0),
('Sa''dullayeva M.X.', 'MAT 24/6', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-05-14 08:38:00', '2025-05-14', 'late', 8),

-- Arzimova S.N (Amaliy o'zbek tili)
('Arzimova S.N', 'MAT 24/4', 'Amaliy o''zbek tili', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-02-06 08:25:00', '2025-02-06', 'present', 0),
('Arzimova S.N', 'MAT 24/6', 'Amaliy o''zbek tili', 'Ma''ruza', '2024-2025', '2-semestr', '10:00:00', '2025-02-20 10:15:00', '2025-02-20', 'late', 15),
('Arzimova S.N', 'MAT 24/4', 'Amaliy o''zbek tili', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-03-06 11:30:00', '2025-03-06', 'present', 0),
('Arzimova S.N', 'MAT 24/6', 'Amaliy o''zbek tili', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-04-17 08:50:00', '2025-04-17', 'late', 20),
('Arzimova S.N', 'MAT 24/4', 'Amaliy o''zbek tili', 'Ma''ruza', '2024-2025', '2-semestr', '10:00:00', '2025-05-08 09:55:00', '2025-05-08', 'present', 0),

-- Jurayev H (O'zbekistonning eng yangi tarixi)
('Jurayev H', 'MAT 24/7', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '08:30:00', '2025-02-07 08:30:00', '2025-02-07', 'present', 0),
('Jurayev H', 'MAT 24/7', 'O''zbekistonning eng yangi tarixi', 'Amaliy', '2024-2025', '2-semestr', '10:00:00', '2025-03-14 10:10:00', '2025-03-14', 'late', 10),
('Jurayev H', 'MAT 24/7', 'O''zbekistonning eng yangi tarixi', 'Ma''ruza', '2024-2025', '2-semestr', '11:30:00', '2025-04-11 11:25:00', '2025-04-11', 'present', 0),

-- Imamova U. (Mediasavodxonlik)
('Imamova U.', 'MAT 24/1', 'Mediasavodxonlik', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-02-25 08:29:00', '2025-02-25', 'present', 0),
('Imamova U.', 'MAT 24/2', 'Mediasavodxonlik', 'Amaliy', '2024-2025', '2-semestr', '10:00:00', '2025-03-25 10:15:00', '2025-03-25', 'late', 15),
('Imamova U.', 'MAT 24/3', 'Mediasavodxonlik', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-04-22 11:30:00', '2025-04-22', 'present', 0),

-- Kadirova X.N (Mediasavodxonlik)
('Kadirova X.N', 'MAT 24/2', 'Mediasavodxonlik', 'Ma''ruza', '2024-2025', '2-semestr', '08:30:00', '2025-02-26 08:25:00', '2025-02-26', 'present', 0),
('Kadirova X.N', 'MAT 24/3', 'Mediasavodxonlik', 'Ma''ruza', '2024-2025', '2-semestr', '10:00:00', '2025-03-26 10:25:00', '2025-03-26', 'late', 25),

-- Yusupova F.M (Amaliy o'zbek tili)
('Yusupova F.M', 'MAT 24/1', 'Amaliy o''zbek tili', 'Amaliy', '2024-2025', '2-semestr', '11:30:00', '2025-02-27 11:28:00', '2025-02-27', 'present', 0),
('Yusupova F.M', 'MAT 24/1', 'Amaliy o''zbek tili', 'Amaliy', '2024-2025', '2-semestr', '08:30:00', '2025-03-27 08:35:00', '2025-03-27', 'late', 5);

INSERT INTO teacher_attendance_logs (teacher_name, group_name, subject_name, lesson_type, academic_year, semester, scheduled_start_time, actual_arrival_time, arrival_date, status, delay_minutes)
VALUES
-- =========================================================
-- 3-SEMESTR (KUZGI SEMESTR: 2025 Sentyabr - Dekabr)
-- =========================================================

-- Janbayeva M (Umumiy pedagogika - Ma'ruza)
('Janbayeva M', 'MAT 24/1', 'Umumiy pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-09-03 08:25:00', '2025-09-03', 'present', 0),
('Janbayeva M', 'MAT 24/2', 'Umumiy pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-09-17 08:40:00', '2025-09-17', 'late', 10),
('Janbayeva M', 'MAT 24/3', 'Umumiy pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '10:00:00', '2025-10-01 10:00:00', '2025-10-01', 'present', 0),
('Janbayeva M', 'MAT 24/4', 'Umumiy pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '10:00:00', '2025-11-12 10:15:00', '2025-11-12', 'late', 15),
('Janbayeva M', 'MAT 24/6', 'Umumiy pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '11:30:00', '2025-12-03 11:28:00', '2025-12-03', 'present', 0),

-- Vafoyeva G (Umumiy pedagogika - Amaliy)
('Vafoyeva G', 'MAT 24/1', 'Umumiy pedagogika', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-09-04 10:12:00', '2025-09-04', 'late', 12),
('Vafoyeva G', 'MAT 24/2', 'Umumiy pedagogika', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-09-18 11:30:00', '2025-09-18', 'present', 0),
('Vafoyeva G', 'MAT 24/3', 'Umumiy pedagogika', 'Amaliy', '2025-2026', '3-semestr', '08:30:00', '2025-10-16 08:50:00', '2025-10-16', 'late', 20),
('Vafoyeva G', 'MAT 24/4', 'Umumiy pedagogika', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-11-06 09:55:00', '2025-11-06', 'present', 0),

-- Nosirova Z.X. (Bolalar adabiyoti - Ma'ruza)
('Nosirova Z.X.', 'MAT 24/1', 'Bolalar adabiyoti', 'Ma''ruza', '2025-2026', '3-semestr', '11:30:00', '2025-09-05 11:25:00', '2025-09-05', 'present', 0),
('Nosirova Z.X.', 'MAT 24/2', 'Bolalar adabiyoti', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-09-19 08:35:00', '2025-09-19', 'late', 5),
('Nosirova Z.X.', 'MAT 24/4', 'Bolalar adabiyoti', 'Ma''ruza', '2025-2026', '3-semestr', '10:00:00', '2025-10-24 10:00:00', '2025-10-24', 'present', 0),
('Nosirova Z.X.', 'MAT 24/5', 'Bolalar adabiyoti', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-11-21 08:45:00', '2025-11-21', 'late', 15),

-- Kodirova M.E. (Bolalar adabiyoti - Amaliy & Tarbiyachining ish hujjatlari)
('Kodirova M.E.', 'MAT 24/1', 'Bolalar adabiyoti', 'Amaliy', '2025-2026', '3-semestr', '08:30:00', '2025-09-09 08:30:00', '2025-09-09', 'present', 0),
('Kodirova M.E.', 'MAT 24/3', 'Bolalar adabiyoti', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-10-07 10:10:00', '2025-10-07', 'late', 10),
('Kodirova M.E.', 'MAT 24/6', 'Tarbiyachining ish hujjatlari', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-11-18 11:28:00', '2025-11-18', 'present', 0),
('Kodirova M.E.', 'MAT 24/6', 'Tarbiyachining ish hujjatlari', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-12-09 11:50:00', '2025-12-09', 'late', 20),

-- Achilova M.S. (Tasviriy faoliyatga o'rgatish metodikasi)
('Achilova M.S.', 'MAT 24/1', 'Tasviriy faoliyatga o''rgatish metodikasi', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-09-10 09:58:00', '2025-09-10', 'present', 0),
('Achilova M.S.', 'MAT 24/2', 'Tasviriy faoliyatga o''rgatish metodikasi', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-10-08 08:42:00', '2025-10-08', 'late', 12),
('Achilova M.S.', 'MAT 24/4', 'Tasviriy faoliyatga o''rgatish metodikasi', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-11-05 11:30:00', '2025-11-05', 'present', 0),
('Achilova M.S.', 'MAT 24/6', 'Tasviriy faoliyatga o''rgatish metodikasi', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-12-10 10:25:00', '2025-12-10', 'late', 25),

-- Shanasirova Z.Y (Tabiat bilan tanishtirish metodikasi)
('Shanasirova Z.Y', 'MAT 24/1', 'Tabiat bilan tanishtirish metodikasi', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-09-11 08:25:00', '2025-09-11', 'present', 0),
('Shanasirova Z.Y', 'MAT 24/2', 'Tabiat bilan tanishtirish metodikasi', 'Ma''ruza', '2025-2026', '3-semestr', '10:00:00', '2025-10-09 10:15:00', '2025-10-09', 'late', 15),
('Shanasirova Z.Y', 'MAT 24/5', 'Tabiat bilan tanishtirish metodikasi', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-11-20 11:28:00', '2025-11-20', 'present', 0),

-- Ibadullayeva SH. (Inklyuziv ta'lim. Gospital pedagogika)
('Ibadullayeva SH.', 'MAT 24/1', 'Inklyuziv ta''lim. Gospital pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '10:00:00', '2025-09-12 09:55:00', '2025-09-12', 'present', 0),
('Ibadullayeva SH.', 'MAT 24/2', 'Inklyuziv ta''lim. Gospital pedagogika', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-10-10 11:45:00', '2025-10-10', 'late', 15),
('Ibadullayeva SH.', 'MAT 24/3', 'Inklyuziv ta''lim. Gospital pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-11-07 08:30:00', '2025-11-07', 'present', 0),
('Ibadullayeva SH.', 'MAT 24/4', 'Inklyuziv ta''lim. Gospital pedagogika', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-12-05 10:20:00', '2025-12-05', 'late', 20),

-- Jumaniyazov J (Inklyuziv ta'lim. Gospital pedagogika) - MAT 24/5
('Jumaniyazov J', 'MAT 24/5', 'Inklyuziv ta''lim. Gospital pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-09-15 08:29:00', '2025-09-15', 'present', 0),
('Jumaniyazov J', 'MAT 24/5', 'Inklyuziv ta''lim. Gospital pedagogika', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-10-13 10:08:00', '2025-10-13', 'late', 8),
('Jumaniyazov J', 'MAT 24/5', 'Inklyuziv ta''lim. Gospital pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '11:30:00', '2025-11-17 11:30:00', '2025-11-17', 'present', 0),
('Jumaniyazov J', 'MAT 24/5', 'Inklyuziv ta''lim. Gospital pedagogika', 'Amaliy', '2025-2026', '3-semestr', '08:30:00', '2025-12-15 08:45:00', '2025-12-15', 'late', 15),

-- Urmanova T.I. (Tarbiyachining ish hujjatlari)
('Urmanova T.I.', 'MAT 24/1', 'Tarbiyachining ish hujjatlari', 'Ma''ruza', '2025-2026', '3-semestr', '11:30:00', '2025-09-16 11:28:00', '2025-09-16', 'present', 0),
('Urmanova T.I.', 'MAT 24/2', 'Tarbiyachining ish hujjatlari', 'Amaliy', '2025-2026', '3-semestr', '08:30:00', '2025-10-14 08:40:00', '2025-10-14', 'late', 10),
('Urmanova T.I.', 'MAT 24/3', 'Tarbiyachining ish hujjatlari', 'Ma''ruza', '2025-2026', '3-semestr', '10:00:00', '2025-11-11 10:00:00', '2025-11-11', 'present', 0),
('Urmanova T.I.', 'MAT 24/4', 'Tarbiyachining ish hujjatlari', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-12-02 11:55:00', '2025-12-02', 'late', 25),

-- Karimov R.R (Umumiy pedagogika & Inklyuziv ta'lim) - MAT 24/7
('Karimov R.R', 'MAT 24/7', 'Umumiy pedagogika', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-09-17 08:25:00', '2025-09-17', 'present', 0),
('Karimov R.R', 'MAT 24/7', 'Inklyuziv ta''lim. Gospital pedagogika', 'Amaliy', '2025-2026', '3-semestr', '10:00:00', '2025-10-15 10:15:00', '2025-10-15', 'late', 15),

-- Boliyeva L (Tarbiyachining ish hujjatlari)
('Boliyeva L', 'MAT 24/5', 'Tarbiyachining ish hujjatlari', 'Amaliy', '2025-2026', '3-semestr', '11:30:00', '2025-09-25 11:30:00', '2025-09-25', 'present', 0),
('Boliyeva L', 'MAT 24/7', 'Tarbiyachining ish hujjatlari', 'Ma''ruza', '2025-2026', '3-semestr', '08:30:00', '2025-11-20 08:42:00', '2025-11-20', 'late', 12);