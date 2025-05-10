-- DROP semua tabel yang bergantung ke `questions`
DROP TABLE IF EXISTS donation_questions;
DROP TABLE IF EXISTS question_links;
DROP TABLE IF EXISTS question_saved;
DROP TABLE IF EXISTS question_mistakes;
DROP TABLE IF EXISTS user_questions;

-- Baru DROP tabel `questions` utama
DROP TABLE IF EXISTS questions;
