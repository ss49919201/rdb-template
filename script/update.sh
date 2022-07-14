#!bin/bash

BEGIN;

SELECT * FROM users WHERE id = 1;

UPDATE users SET count = count + 1 WHERE id = '1';

COMMIT;
