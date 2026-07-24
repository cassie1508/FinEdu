CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS flashcards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    why_it_matters TEXT NOT NULL,
    definition TEXT NOT NULL,
    example TEXT NOT NULL,
    common_misconceptions TEXT DEFAULT ''::TEXT,
    review_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============ SQL QUERIES FOR FLASHCARD OPERATIONS ============

-- 1. GET ALL FLASHCARDS (GetAllFlashcards)
-- SELECT * FROM flashcards WHERE user_id = $1 ORDER BY created_at DESC;

-- 2. GET FLASHCARD BY ID (GetFlashcardByID)
-- SELECT * FROM flashcards WHERE id = $1 AND user_id = $2;

-- 3. CREATE FLASHCARD (CreateFlashcard)
-- INSERT INTO flashcards (user_id, title, category, why_it_matters, definition, example, common_misconceptions, review_count, created_at, updated_at)
-- VALUES ($1, $2, $3, $4, $5, $6, $7, 0, NOW(), NOW())
-- RETURNING id, user_id, title, category, why_it_matters, definition, example, common_misconceptions, review_count, created_at, updated_at;

-- 4. UPDATE FLASHCARD (UpdateFlashcard)
-- UPDATE flashcards 
-- SET title = $1, category = $2, why_it_matters = $3, definition = $4, example = $5, common_misconceptions = $6, updated_at = NOW()
-- WHERE id = $7 AND user_id = $8
-- RETURNING id, user_id, title, category, why_it_matters, definition, example, common_misconceptions, review_count, created_at, updated_at;

-- 5. DELETE FLASHCARD (DeleteFlashcard)
-- DELETE FROM flashcards WHERE id = $1 AND user_id = $2;

-- 6. REVIEW FLASHCARD - INCREMENT REVIEW COUNT (ReviewFlashcard)
-- UPDATE flashcards 
-- SET review_count = review_count + 1, updated_at = NOW()
-- WHERE id = $1 AND user_id = $2
-- RETURNING id, user_id, title, category, why_it_matters, definition, example, common_misconceptions, review_count, created_at, updated_at;

;