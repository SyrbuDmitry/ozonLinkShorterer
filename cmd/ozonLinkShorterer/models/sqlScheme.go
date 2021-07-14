package models

//SQL-запросы

const (
	Scheme = `CREATE TABLE IF NOT EXISTS urls(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				originalUrl VARCHAR(220)
              );
				INSERT INTO SQLITE_SEQUENCE SELECT 'urls',100
					WHERE NOT EXISTS
						(SELECT*FROM SQLITE_SEQUENCE);
			 `
)
