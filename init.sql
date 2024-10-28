-- Создаем таблицу transactions
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    amount NUMERIC(10, 2) NOT NULL,
    transaction_type VARCHAR(10) NOT NULL CHECK (transaction_type IN ('доход', 'расход')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Можно добавить любые начальные данные, если нужно