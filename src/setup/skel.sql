PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS "User";
CREATE TABLE IF NOT EXISTS "User" (
	"Id" INTEGER,
    "Name" TEXT NOT NULL,
    "TinNumber" TEXT NOT NULL,
    "Address" TEXT NOT NULL,
    "City" TEXT NOT NULL,
    "PostalCode" TEXT NOT NULL,
    "Country" TEXT NOT NULL,
    "Phone" TEXT NOT NULL,
    "Email" TEXT NOT NULL,
	PRIMARY KEY("Id" AUTOINCREMENT)
);

DROP TABLE IF EXISTS "Customers";
CREATE TABLE IF NOT EXISTS "Customers" (
	"Id" INTEGER,
    "Name" TEXT NOT NULL,
    "TinNumber" TEXT NOT NULL,
    "Address" TEXT NOT NULL,
    "City" TEXT NOT NULL,
    "PostalCode" TEXT NOT NULL,
    "Country" TEXT NOT NULL,
    "Phone" TEXT NOT NULL,
    "Email" TEXT NOT NULL,
	PRIMARY KEY("Id" AUTOINCREMENT)
);

DROP TABLE IF EXISTS "IssuedInvoices";
CREATE TABLE IF NOT EXISTS "IssuedInvoices" (
    "Id" INTEGER,
    "Number" TEXT NOT NULL UNIQUE,
    "Date" INTEGER NOT NULL,
    "Data" TEXT NOT NULL,
    "Result" TEXT NOT NULL,
    PRIMARY KEY("Id" AUTOINCREMENT)
);
CREATE INDEX IF NOT EXISTS "IssuedInvoices_Number" ON "IssuedInvoices" ("Number");
