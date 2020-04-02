## Formbot

This is the Go action server implementation for the 
[Rasa Open Source formbot example](https://github.com/RasaHQ/rasa/tree/master/examples/formbot).

### Running it

1 Add the action server to your `endpoints.yml`:

    ```yaml
    action_endpoint:
      url: "http://localhost:5055/webhook"
    ```
2. Run Rasa Open Source using `rasa run`

3.  Run `go run main.go`. This should give you output similar to this:

    ```bash
    ❯ go run main.go
    INFO[0000] Action server running on on port 5055
    ```
 