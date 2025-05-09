-- Drop dulu yang bergantung ke questions
DROP TABLE IF EXISTS question_links;
DROP TABLE IF EXISTS question_saved;
DROP TABLE IF EXISTS question_mistakes;
DROP TABLE IF EXISTS user_questions;

-- Baru drop parent-nya
DROP TABLE IF EXISTS questions;
