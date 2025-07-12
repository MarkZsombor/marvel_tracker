-- Create heroes table
CREATE TABLE IF NOT EXISTS heroes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create scenarios table
CREATE TABLE IF NOT EXISTS scenarios (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create plays table
CREATE TABLE IF NOT EXISTS plays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL,
    outcome TEXT NOT NULL CHECK(outcome IN ('win', 'loss')),
    difficulty TEXT NOT NULL,
    notes TEXT,
    scenario_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (scenario_id) REFERENCES scenarios(id)
);

-- Create decks table (many-to-many relationship between plays and heroes)
CREATE TABLE IF NOT EXISTS decks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    play_id INTEGER NOT NULL,
    hero_id INTEGER NOT NULL,
    aspect TEXT NOT NULL CHECK(aspect IN ('leadership', 'justice', 'aggression', 'protection')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (play_id) REFERENCES plays(id) ON DELETE CASCADE,
    FOREIGN KEY (hero_id) REFERENCES heroes(id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_plays_date ON plays(date);
CREATE INDEX IF NOT EXISTS idx_plays_scenario_id ON plays(scenario_id);
CREATE INDEX IF NOT EXISTS idx_decks_play_id ON decks(play_id);
CREATE INDEX IF NOT EXISTS idx_decks_hero_id ON decks(hero_id);