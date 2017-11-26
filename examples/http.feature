Feature: HTTP
  Scenario: Import HTTP module
    Given a file named "main.tisp" with:
    """
    (import "http")
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly ""

  Scenario: Send GET request
    Given a file named "main.tisp" with:
    """
    (import "http")

    (write ((http.Get "https://google.com") "status"))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly "200"

  Scenario: Send POST request
    Given a file named "main.tisp" with:
    """
    (import "http")

    (write ((http.Post "https://google.com" "Hello, world!") "status"))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly "405"

  Scenario: Run a server
    Given a file named "main.tisp" with:
    """
    (import "http")

    ..(map (\ (r) ((r "respond") "Hello, world!")) (http.GetRequests ":8080"))
    """
    And a file named "main.sh" with:
    """
    tisp main.tisp &
    pid=$!
    sleep 1 &&
    curl http://127.0.0.1:8080 &&
    kill $pid
    """
    When I successfully run `sh ./main.sh`
    Then the stdout should contain exactly "Hello, world!"
