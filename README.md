# Air Invoicing
## Simple and free invoicing for freelancers

> ⚠️ Check your local laws and regulations to ensure that you are compliant with invoicing requirements

Air Invoicing is a simple and free invoicing tool for freelancers, designed to create invoices quickly and easily with an extensible template system.

### Features
- Save user data
- Save clients
- Create invoices
- Save created invoices

### Screenshots

#### CLI Interface

```
       d8888 8888888 8888888b.  
      d88888   888   888   Y88b 
     d88P888   888   888    888 
    d88P 888   888   888   d88P 
   d88P  888   888   8888888P"  
  d88P   888   888   888 T88b   
 d8888888888   888   888  T88b  
d88P     888 8888888 888   T88b 


Welcome to Air!
Version: 0.1.0 


------ Main Menu ------
Please select an option:
  001 - Set issuer data
  002 - List issuer data
  003 - Delete issuer data
  
  101 - Create new customer
  102 - List existing customer
  103 - Delete existing customer
  
  201 - Create new invoice
  202 - List existing invoice
  
  999 - Quit

> 
```

#### Example invoice

![ExampleInvoice](https://github.com/user-attachments/assets/6e1b9d40-077e-4e0b-ad4a-24274e9c5fbe)

### Installation

Build or [download](https://github.com/RLado/Air/releases) a package for your platform:

Build:
```bash
make package # Optionally specify 'DESTDIR' to change the output directory
```

Then run the executable initiallizing the database:
```bash
./air -i
```

After that, you can run the application normally:
```bash
./air
```
