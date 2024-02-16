CREATE TABLE IF NOT EXISTS clientes (
  id SERIAL PRIMARY KEY,
  limite INT NOT NULL,
  saldo INT NOT NULL DEFAULT 0
);

CREATE TYPE tipo_transacao AS ENUM ('c', 'd');

-- TIMESTAMP não tem fuso
-- TIMESTAMPTZ tem fuso
-- Vendo o repo da rinha, acho que é pra usar sem fuso?
CREATE TABLE IF NOT EXISTS transacoes (
  id SERIAL PRIMARY KEY,
  valor INT NOT NULL,
  tipo tipo_transacao NOT NULL,
  descricao VARCHAR(10) CHECK (LENGTH(descricao) >= 1),
  realizada_em TIMESTAMP NOT NULL,
  cliente_id INT REFERENCES clientes (id) ON DELETE CASCADE
);
