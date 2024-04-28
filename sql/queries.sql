-- all tables from schema
SELECT owner, table_name
FROM all_tables
WHERE owner = 'DEMO';

-- all columns from table
select owner, table_name, column_name, data_type
from all_tab_cols
WHERE owner = 'DEMO'
AND table_name = 'AUTHORS';

-- primary key from table
SELECT c2.owner, c2.table_name, c2.column_name
FROM user_constraints c1, all_cons_columns c2
WHERE c1.constraint_type = 'P'
AND c1.constraint_name = c2.constraint_name
AND c2.table_name = 'AUTHORS';


-- foreign key from table
SELECT a.owner, a.table_name, a.column_name, c_pk.table_name r_table_name
FROM all_cons_columns a
JOIN all_constraints c ON a.owner = c.owner
AND a.constraint_name = c.constraint_name
JOIN all_constraints c_pk ON c.r_owner = c_pk.owner
AND c.r_constraint_name = c_pk.constraint_name
WHERE c.constraint_type = 'R'
AND a.table_name = 'AUTHORISBN';