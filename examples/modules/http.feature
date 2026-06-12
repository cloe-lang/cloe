Feature: HTTP
  Scenario: Import HTTP module
    Given a file named "main.cloe" with:
    """
    (import "http")
    """
    When I successfully run `cloe main.cloe`
    Then the stdout should contain exactly ""

  Scenario: Send GET request
    Given a file named "server.cloe" with:
    """
    (import "http")

    ..(map (\ (r) ((@ r "respond") "Hello, world!")) (http.getRequests ":8081"))
    """
    And a file named "main.cloe" with:
    """
    (import "http")

    (print (@ (http.get "http://127.0.0.1:8081") "status"))
    """
    And a file named "main.sh" with:
    """
    set -e

    cloe server.cloe &
    pid=$!
    sleep 1
    cloe main.cloe
    kill $pid
    """
    When I successfully run `sh ./main.sh`
    Then the stdout should contain exactly "200"

  Scenario: Send POST request
    Given a file named "server.cloe" with:
    """
    (import "http")

    ..(map (\ (r) ((@ r "respond") "Hello, world!")) (http.getRequests ":8082"))
    """
    And a file named "main.cloe" with:
    """
    (import "http")

    (print (@ (http.post "http://127.0.0.1:8082" "Hello, world!") "status"))
    """
    And a file named "main.sh" with:
    """
    set -e

    cloe server.cloe &
    pid=$!
    sleep 1
    cloe main.cloe
    kill $pid
    """
    When I successfully run `sh ./main.sh`
    Then the stdout should contain exactly "200"

  Scenario: Run a server
    Given a file named "main.cloe" with:
    """
    (import "http")

    ..(map (\ (r) ((@ r "respond") "Hello, world!")) (http.getRequests ":8080"))
    """
    And a file named "main.sh" with:
    """
    set -e

    cloe main.cloe &
    pid=$!
    sleep 1
    curl http://127.0.0.1:8080
    kill $pid
    """
    When I successfully run `sh ./main.sh`
    Then the stdout should contain exactly "Hello, world!"
