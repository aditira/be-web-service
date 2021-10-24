CREATE TABLE mahasiswa (
    id INT IDENTITY(1,1) PRIMARY KEY,
    nim character varying(50) NOT NULL,
    name character varying(255) NOT NULL,
    gender character varying(30) NOT NULL,
    class_name character varying(20) NOT NULL,
    class_type character varying(20) NOT NULL,
    subject character varying(50) NOT NULL
);