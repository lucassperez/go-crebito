INSERT INTO clientes (limite)
VALUES (100000), (80000), (1000000), (10000000), (500000);

INSERT INTO transacoes (valor, tipo, descricao, realizada_em, cliente_id)
VALUES (1000, 'c', 'teste-cred', CURRENT_TIMESTAMP, 1);

INSERT INTO transacoes (valor, tipo, descricao, realizada_em, cliente_id)
VALUES (50, 'd', 'teste-cred', CURRENT_TIMESTAMP, 1);
