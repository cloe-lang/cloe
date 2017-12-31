Feature: Performance
  Scenario: Benchmark requests per minute of HTTP server
    Given a file named "main.coel" with:
    """
    (import "http")

    ..(map (\ (r) ((r "respond") "Hello, world!")) (http.getRequests ":8080"))
    """
    And a file named "main.go" with:
    """
    package main

    import "net/http"

    type handler struct{}

    func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("Hello, world!\n"))
    }

    func main() {
      err := http.ListenAndServe(":8080", handler{})

      if err != nil {
        panic(err)
      }
    }
    """
    And a file named "main.sh" with:
    """
    wait() {
      sleep 1
    }

    bench() {
      $@ &
      pid=$!

      wait

      rpm -n 50000 -c 1000 http://localhost:8080
      kill -9 $pid
    }

    coel=$(bench coel main.coel)
    wait
    go=$(bench ./main)

    python -c "if $coel / $go < 0.15: exit(1)"
    """
    When I successfully run `go build main.go`
    And I successfully run `sh main.sh`
