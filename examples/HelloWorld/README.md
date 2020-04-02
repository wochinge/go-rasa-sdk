# Hello World

This is the Go action server implementation for a custom action `action_hello_world` which sends the message
`Hello world` to the user when it's triggered.

### Running it

1. Add `action_hello_world` to your `domain.yml`
2. Write a story to trigger the action, e.g:

    ```
    ## My Go test story
    * greet
      - action_hello_world
    ```
3. Add the action server to your `endpoints.yml`:

    ```yaml
    action_endpoint:
      url: "http://localhost:5055/webhook"
    ```
4. Run Rasa Open Source using `rasa run`

5.  Run `go run main.go`. This should give you output similar to this:

    ```bash
    ❯ go run main.go
    INFO[0000] Action server running on on port 5055
    ```