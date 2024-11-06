# SQHeavy

## API examples
- curl http://localhost:3517/command  -u admin --json '{"command":"select * from t2", "dbName": "tet4"}' -> {"status":"OK","msg":"enomine,8;enom9ine,98;eine,998","rowsAffected":-1}% 
- curl http://localhost:3517/command  -u admin --json '{"command":"CREATE TABLE t2 (s text, n integer)", "dbName": "tet4"}' -> {"status":"OK","msg":"","rowsAffected":0}%
- curl http://localhost:3517/command  -u admin --json '{"command":"select * from namestb")", "dbName": "tet4"}' -> {"status":"Failed","msg":"invalid character ')' after object key:value pair","rowsAffected":-1}% 

