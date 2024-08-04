PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS "User";
CREATE TABLE IF NOT EXISTS "User" (
	"Id" INTEGER NOT NULL,
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
	"Id" INTEGER NOT NULL,
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
    "Id" INTEGER NOT NULL,
    "Series" TEXT NOT NULL,
    "Number" INT NOT NULL,
    "CustomerId" INTEGER NOT NULL,
    "Date" INTEGER NOT NULL,
    "Data" TEXT NOT NULL,
    PRIMARY KEY("Id" AUTOINCREMENT)
    FOREIGN KEY("CustomerId") REFERENCES "Customers"("Id")
);
CREATE INDEX IF NOT EXISTS "IssuedInvoices_Number" ON "IssuedInvoices" ("Number");
