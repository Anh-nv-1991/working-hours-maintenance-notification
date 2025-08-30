## **PostgreSQL Cheat Sheet for Docker**

### 1. **Truy cập vào PostgreSQL trong Docker**

`docker exec -it wh-ma-db psql -U admin -d main-1`

- **container_name**: Tên hoặc ID của container PostgreSQL.

- **username**: Tên người dùng PostgreSQL (thường là `postgres`).

- **database_name**: Tên cơ sở dữ liệu cần truy cập.

### 2. **Các Lệnh PostgreSQL Cơ Bản**

#### a. **Liệt kê các bảng trong cơ sở dữ liệu**

`\dt`

#### b. **Xem cấu trúc của một bảng**

`\d <table_name>`

**Ví dụ**:

`\d devices`

#### c. **Xem dữ liệu trong một bảng (dùng SELECT)**

`SELECT * FROM <table_name> LIMIT 10;`

**Ví dụ**:

`SELECT * FROM devices LIMIT 10;`

#### d. **Lọc dữ liệu theo điều kiện**

`SELECT * FROM <table_name> WHERE <column_name> = '<value>';`

**Ví dụ**:

`SELECT * FROM devices WHERE serial_number = 'SN-005';`

#### e. **Tạo bảng mới**

`CREATE TABLE <table_name> (     id SERIAL PRIMARY KEY,     name VARCHAR(255),     status VARCHAR(50) );`

#### f. **Cập nhật dữ liệu trong bảng**

`UPDATE <table_name> SET <column_name> = '<new_value>' WHERE <condition>;`

**Ví dụ**:

`UPDATE devices SET status = 'inactive' WHERE serial_number = 'SN-005';`

#### g. **Xóa dữ liệu trong bảng**

`DELETE FROM <table_name> WHERE <condition>;`

**Ví dụ**:

`DELETE FROM devices WHERE serial_number = 'SN-005';`

#### h. **Thêm dữ liệu vào bảng**

`INSERT INTO <table_name> (<column1>, <column2>, ...) VALUES ('<value1>', '<value2>', ...);`

**Ví dụ**:

`INSERT INTO devices (serial_number, name, model, manufacturer, year, status) VALUES ('SN-006', 'Loader F', 'CAT 980', 'Caterpillar', 2021, 'active');`

### 3. **Các Lệnh Khác**

#### a. **Xem danh sách cơ sở dữ liệu**

`\l`

#### b. **Chuyển sang cơ sở dữ liệu khác**

`\c <database_name>`

**Ví dụ**:

`\c mydb`
