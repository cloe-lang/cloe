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
    When I run `cloe server.cloe` in the background
    And I wait 1 second for the command to start up
    And I successfully run `cloe main.cloe`
    Then the stdout from "cloe main.cloe" should contain "200"

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
    When I run `cloe server.cloe` in the background
    And I wait 1 second for the command to start up
    And I successfully run `cloe main.cloe`
    Then the stdout from "cloe main.cloe" should contain "200"

  Scenario: Run a server
    Given a file named "main.cloe" with:
    """
    (import "http")

    ..(map (\ (r) ((@ r "respond") "Hello, world!")) (http.getRequests ":8080"))
    """
    When I run `cloe main.cloe` in the background
    And I wait 1 second for the command to start up
    And I successfully run `curl http://127.0.0.1:8080`
    Then the stdout from "curl http://127.0.0.1:8080" should contain "Hello, world!"
