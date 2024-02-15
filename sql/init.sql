CREATE TABLE IF NOT EXISTS clientes (
  id SERIAL PRIMARY KEY,
  limite INT NOT NULL,
  saldo INT NOT NULL DEFAULT 0
);

CREATE TYPE tipo_transacao AS ENUM ('c', 'd');

CREATE TABLE IF NOT EXISTS transacoes (
  id SERIAL PRIMARY KEY,
  tipo tipo_transacao NOT NULL,
  descricao VARCHAR(10) CHECK (LENGTH(descricao) >= 1),
  cliente_id INT UNIQUE REFERENCES clientes (id) ON DELETE CASCADE
);
