# Meilo
Meilo is an implementation of a test/development SMTP server written in Go. Meilo opens sent emails instead of sending them, 
this allows development workflows where ONLY the SMTP credentials are changed.

## Prerequisites
- Meilo requies Go 1.23 or later, you can download it from [here](https://golang.org/dl/)

- Meilo package:
```bash
go get github.com/wawandco/meilo
```

## Usage
1.  Create a new instance of Meilo

You have the ability to set the options for the SMTP server, these are totally optional, if you don't set them, Meilo will use the default values.

```go
// Start the SMTP server
creds, err := meilo.Start(
	// Sender options
	meilo.WithSenderOptions(
		meilo.Only([]string{"text/html", "text/plain"}),
		meilo.WithDir(os.Getenv("TMP_DIR")),
	),

	// SMTP server options
	meilo.WithHost("smtp.example.com"),
	meilo.WithUser("example@emai.com"),
	meilo.WithPort("1025"),
	meilo.WithPassword("password"),
)
 
if err != nil {
	// Handle the error
	...
}
```
This will start the SMTP server and return the credentials to be used in the email sending process.

2.  Send an email using the SMTP server

```go
//Then you can use the credentials to send an email
err = smtp.SendMail(
	creds.Addr(), 				    // Addr
	smtp.PlainAuth("", creds.User, creds.Password, creds.Host), // Authentication

    "username@example.com", 		// From
	[]string{"example@example.com"}, 	// To
	
	[]byte("Hello from meilo!"), // Body
)

if err != nil {
        // Handle the sending error
        ...

}
```
### Options explained

- `meilo.Only([]string{"text/html", "text/plain"})`: allows you to specify which content types you want to open in the browser, by default it will open all content types.

- `meilo.WithDir(os.Getenv("TMP_DIR"))`: allows you to specify the directory where the emails will be stored, by default it will use the system's temporary directory.

- `meilo.WithHost("smtp.example.com")`: allows you to specify the host of the SMTP server.

- `meilo.WithUser("user@user.com")`: allows you to specify the user of the SMTP server.

- `meilo.WithPort("1025")`: allows you to specify the port of the SMTP server.

- `meilo.WithPassword("password")`: allows you to specify the password of the SMTP server.


## Roadmap / Ideas
- Web interface
- UI Improvements.
- Deployable service for Staging/Testing.
- Listing historically sent emails.

